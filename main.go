package main

import (
	"fmt"
	"os"

	"github.com/mmycin/goforge/cmd/console"
)

func main() {
	if err := console.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
