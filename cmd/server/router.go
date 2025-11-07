package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// newRouter sets up the Chi router with middleware and routes.
func (app *application) newRouter() http.Handler {
	return app.httpRestRouter()
}

func (app *application) httpRestRouter() http.Handler {
	router := chi.NewRouter()

	// Add common middleware.
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Define routes
	router.Route("/entities", func(r chi.Router) {
		r.Post("/", app.entityHandler.CreateEntity)
		r.Get("/", app.entityHandler.ListEntities)
		r.Get("/{id}", app.entityHandler.GetEntity)
		r.Put("/{id}", app.entityHandler.UpdateEntity)
		r.Delete("/{id}", app.entityHandler.DeleteEntity)
	})

	return router
}
