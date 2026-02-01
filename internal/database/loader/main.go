// internal/database/loader/main.go
package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"github.com/mmycin/goforge/internal/database" // your package with Model()
	// Make sure this import path matches your module name + path
)

func main() {
	models := database.Model() // ← calls your function that returns []any{&todo.Todo{}, ...}

	loader := gormschema.New("sqlite") // or WithConfig(...) if needed, e.g. to disable FKs

	stmts, err := loader.Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, stmts)
}
