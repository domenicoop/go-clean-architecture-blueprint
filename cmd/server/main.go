package main

import (
	"log/slog"
	"os"
)

// main is the entry point for the application.
func main() {
	// 1. Initialize a structured, production-ready logger.
	// JSON format is great for log aggregators like Datadog or Splunk.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Load configuration from the environment.
	cfg := loadConfig()

	// 3. Set up the application struct, which holds all our dependencies.
	app := newApplication(cfg, logger)

	// 4. Run the server, with graceful shutdown.
	if err := app.serve(); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
