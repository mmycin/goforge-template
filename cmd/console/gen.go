package console

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var genSqlcCmd = &cobra.Command{
	Use:   "gen:sqlc",
	Short: "Run code generation",
	Long:  `Execute sqlc generate to create database query code.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running code generation...")
		makeGen()
	},
}

func makeGen() {
	fmt.Println("Executing sqlc generate...")
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: sqlc generate failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Code generation completed successfully")
}
