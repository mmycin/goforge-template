package console

import (
	"log"
	"os"
	"os/exec"
)

func makeGen() {
	log.Println("Running sqlc generate...")
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("sqlc generate failed: %v", err)
	}
	log.Println("Code generation completed.")
}
