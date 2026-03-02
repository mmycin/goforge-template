package console

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/services"
	"github.com/spf13/cobra"
)

var loaderCmd = &cobra.Command{
	Use:   "loader",
	Short: "Run GORM schema loader",
	Long:  `Load and display GORM schema definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLoader()
	},
}

func runLoader() {
	models := services.Model()
	driver := config.DB.Connection
	if driver == "" {
		driver = "sqlite"
	}
	loader := gormschema.New(driver)
	stmts, err := loader.Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load GORM schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, stmts)
}
