package main

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/mmycin/goforge/internal/services/todo"
)

func Model() []any {
	return []any{
		&todo.Todo{},
	}
}

func main() {
	models := Model()

	loader := gormschema.New("sqlite")

	stmts, err := loader.Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, stmts)
}
