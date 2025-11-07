package main

import (
	"os"
)

// config holds all the configuration for the application.
// Values are read from environment variables.
type config struct {
	port string // Network port to listen on
	env  string // Current operating environment (e.g., development, production)
}

// loadConfig loads configuration from environment variables.
func loadConfig() config {
	var cfg config

	cfg.port = os.Getenv("PORT")
	if cfg.port == "" {
		cfg.port = "8080" // Default port
	}

	cfg.env = os.Getenv("ENV")
	if cfg.env == "" {
		cfg.env = "development"
	}

	return cfg
}
