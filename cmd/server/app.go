package main

import (
	"log/slog"

	httpHandler "github.com/domenicoop/go-clean-architecture-blueprint/internal/handler/http"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/repository/inmemory"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/service"
)

// application struct holds the dependencies for our application.
// This is a common pattern for dependency injection in Go.
type application struct {
	config config
	logger *slog.Logger

	// handlers
	entityHandler *httpHandler.EntityHandler
}

func newApplication(cfg config, logger *slog.Logger) *application {
	// Wire up dependencies: repository -> service -> handler
	entityRepo := inmemory.NewEntityRepository()
	entityService := service.NewEntityService(entityRepo)
	entityHandler := httpHandler.NewEntityHandler(entityService, logger)

	return &application{
		config:        cfg,
		logger:        logger,
		entityHandler: entityHandler,
	}
}
