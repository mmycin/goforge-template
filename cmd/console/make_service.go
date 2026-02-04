package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeServiceCmd = &cobra.Command{
	Use:   "make:service [name]",
	Short: "Create a new service",
	Long:  `Generate a new service with handler, repository, model, routes, and proto files.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Creating service: %s\n", name)
		makeService(name)
	},
}
