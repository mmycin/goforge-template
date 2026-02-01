package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Execute() {
	loadEnv()

	if len(os.Args) < 2 {
		help()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "make:migration":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run . make:migration <name>")
		}
		makeMigration(os.Args[2])
	case "migate", "migrate":
		migrateDB()
	case "loader":
		runLoader()
	case "make":
		if len(os.Args) < 4 || os.Args[2] != "service" {
			log.Fatal("Usage: go run . make service <name>")
		}
		makeService(os.Args[3])
	case "gen":
		makeGen()
	default:
		help()
	}
}

func help() {
	fmt.Println("Available commands:")
	fmt.Println("  make:migration <name>")
	fmt.Println("  migate (runs auto-migration)")
	fmt.Println("  make service <name>")
	fmt.Println("  gen (runs sqlc generate)")
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Println("Warning: could not open .env file")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}
