package console

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmycin/goforge/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// initLogger initializes the global logger based on configuration
func initLogger() {
	var writers []io.Writer

	// Setup log type
	switch strings.ToLower(config.Log.Type) {
	case "file":
		if config.Log.Path == "" {
			fmt.Println("Warning: LOG_PATH is empty, falling back to console logging")
			writers = append(writers, os.Stdout)
		} else {
			// Ensure directory exists
			dir := filepath.Dir(config.Log.Path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error creating log directory %s: %v\n", dir, err)
				writers = append(writers, os.Stdout)
			} else {
				file, err := os.OpenFile(config.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Printf("Error opening log file %s: %v\n", config.Log.Path, err)
					writers = append(writers, os.Stdout)
				} else {
					writers = append(writers, file)
				}
			}
		}
	case "console":
		fallthrough
	default:
		writers = append(writers, os.Stdout)
	}

	var output io.Writer
	if len(writers) == 1 {
		output = writers[0]
	} else {
		output = io.MultiWriter(writers...)
	}

	// Setup format
	if strings.ToLower(config.Log.Format) == "text" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}
	}

	// Set global logger
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// Setup log level
	level, err := zerolog.ParseLevel(strings.ToLower(config.Log.Level))
	if err != nil {
		fmt.Printf("Warning: invalid LOG_LEVEL '%s', defaulting to 'info'\n", config.Log.Level)
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Info().Msg("Logger initialized")
}
