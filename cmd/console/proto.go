package console

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

var genProtoCmd = &cobra.Command{
	Use:   "gen:proto [service]",
	Short: "Compile proto files and generate gRPC code",
	Long:  `Compiles .proto files using protoc and generates Go server scaffolding in internal/services.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		if err := generateProto(name); err != nil {
			Error("%v", err)
			os.Exit(1)
		}
	},
}

func generateProto(serviceName string) error {
	protoDir := "proto"
	servicesDir := "internal/services"

	// Find proto files
	var protoFiles []string
	if serviceName != "" {
		// Specific service
		p := filepath.Join(protoDir, serviceName, serviceName+".proto")
		if _, err := os.Stat(p); err == nil {
			protoFiles = append(protoFiles, p)
		} else {
			return fmt.Errorf("proto file not found: %s", p)
		}
	} else {
		// Scan all
		entries, err := os.ReadDir(protoDir)
		if err != nil {
			if os.IsNotExist(err) {
				Info("No proto directory found.")
				return nil
			}
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				p := filepath.Join(protoDir, e.Name(), e.Name()+".proto")
				if _, err := os.Stat(p); err == nil {
					protoFiles = append(protoFiles, p)
				}
			}
		}
	}

	if len(protoFiles) == 0 {
		Info("No proto files found.")
		return nil
	}

	moduleName, err := getModuleName()
	if err != nil {
		Warning("Could not determine module name: %v. Generated paths might be incorrect.", err)
	}

	for _, p := range protoFiles {
		Info("Compiling %s...", p)

		args := []string{
			"--go_out=.",
			"--go-grpc_out=.",
			p,
		}

		if moduleName != "" {
			args = append(args, "--go_opt=module="+moduleName)
			args = append(args, "--go-grpc_opt=module="+moduleName)
		}

		cmd := exec.Command("protoc", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("protoc failed for %s: %w", p, err)
		}

		// Generate Scaffold
		// Extract service name from path: proto/<service>/<service>.proto
		parts := strings.Split(filepath.ToSlash(p), "/")
		if len(parts) >= 2 {
			svcName := parts[len(parts)-2]
			if err := generateGrpcScaffold(servicesDir, svcName); err != nil {
				return err
			}
		}
	}

	Success("Proto compilation and scaffolding complete.")
	return nil
}

func getModuleName() (string, error) {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module directive not found in go.mod")
}

func generateGrpcScaffold(servicesDir, serviceName string) error {
	svcDir := filepath.Join(servicesDir, serviceName)
	if err := os.MkdirAll(svcDir, 0755); err != nil {
		return err
	}

	grpcFile := filepath.Join(svcDir, "grpc.go")
	Info("  Generating scaffold: %s", grpcFile)

	camelName := toCamelCase(serviceName)
	methods, err := parseProtoMethods(serviceName)
	if err != nil {
		Warning("  Could not parse proto file for stub generation: %v", err)
	}

	data := struct {
		Package   string
		Title     string
		Methods   []string
		Timestamp string
	}{
		Package:   serviceName,
		Title:     camelName,
		Methods:   methods,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	tmpl := `// Generated at {{.Timestamp}}
package {{.Package}}

import (
	"context"

	"github.com/mmycin/goforge/internal/server"
	pb "github.com/mmycin/goforge/proto/{{.Package}}/gen"
	"google.golang.org/grpc"
)

type {{.Title}}GRPC struct {
	pb.Unimplemented{{.Title}}ServiceServer
}

func init() {
	server.RegisterGRPC(func(s *grpc.Server) {
		pb.Register{{.Title}}ServiceServer(s, &{{.Title}}GRPC{})
	})
}

{{range .Methods}}
func (t *{{$.Title}}GRPC) {{.}} {
	return nil, nil
}
{{end}}
`

	t, err := template.New("grpc").Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(grpcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			r := []rune(part)
			r[0] = unicode.ToUpper(r[0])
			parts[i] = string(r)
		}
	}
	return strings.Join(parts, "")
}

func parseProtoMethods(serviceName string) ([]string, error) {
	protoFile := filepath.Join("proto", serviceName, serviceName+".proto")
	content, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var methods []string
	inService := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "service ") {
			inService = true
			continue
		}
		if inService {
			if strings.HasPrefix(line, "}") {
				break
			}
			if strings.HasPrefix(line, "rpc ") {
				// rpc Create(CreateRequest) returns (CreateResponse);
				parts := strings.Fields(line)
				if len(parts) >= 4 {
					methodName := parts[1]
					// Remove parentheses if name is like "Create("
					methodName = strings.Split(methodName, "(")[0]

					// Extract request type: (CreateRequest)
					reqPart := strings.Join(parts[1:], " ")
					reqStart := strings.Index(reqPart, "(")
					reqEnd := strings.Index(reqPart, ")")
					if reqStart != -1 && reqEnd != -1 {
						reqType := reqPart[reqStart+1 : reqEnd]

						// Extract response type: (CreateResponse)
						respPart := reqPart[reqEnd+1:]
						resStart := strings.Index(respPart, "(")
						resEnd := strings.Index(respPart, ")")
						if resStart != -1 && resEnd != -1 {
							resType := respPart[resStart+1 : resEnd]
							// Construct method signature
							methods = append(methods, fmt.Sprintf("%s(ctx context.Context, in *pb.%s) (*pb.%s, error)", methodName, reqType, resType))
						}
					}
				}
			}
		}
	}

	return methods, nil
}
