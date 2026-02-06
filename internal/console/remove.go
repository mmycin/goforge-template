package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var removeServiceCmd = &cobra.Command{
	Use:   "rem:service [name]",
	Short: "Remove an existing service",
	Long:  `Permanently remove a service, including its directory, proto files, and kernel registration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		Info("Removing service: %s", name)
		removeService(name)
	},
}

func removeService(name string) {
	servicesDir := filepath.Join("internal/services", name)
	protoDir := filepath.Join("proto", name)

	// Check if service exists
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		Error("Service '%s' does not exist", name)
		os.Exit(1)
	}

	// Remove internal/services/<name>
	if err := os.RemoveAll(servicesDir); err != nil {
		Error("Failed to remove service directory: %v", err)
		os.Exit(1)
	}

	// Remove proto/<name>
	if err := os.RemoveAll(protoDir); err != nil {
		Error("Failed to remove proto directory: %v", err)
	}

	// Remove from kernel
	removeFromKernel(name)

	Success("Service '%s' removed successfully", name)
}

func removeFromKernel(name string) {
	kernelPath := "internal/services/kernel.go"
	importLine := fmt.Sprintf("\"github.com/mmycin/goforge/internal/services/%s\"", name)

	content, err := os.ReadFile(kernelPath)
	if err != nil {
		Warning("Could not read %s: %v", kernelPath, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	removed := false

	for _, line := range lines {
		if strings.Contains(line, importLine) {
			removed = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !removed {
		Warning("Could not find registration in %s. Please check manually.", kernelPath)
		return
	}

	err = os.WriteFile(kernelPath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		Warning("Could not write to %s: %v", kernelPath, err)
	}
}
