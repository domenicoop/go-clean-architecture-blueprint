package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve starts the HTTP server and handles graceful shutdown.
func (app *application) serve() error {
	// Create a custom HTTP server with timeouts. This is crucial for production
	// to prevent resource exhaustion from slow or malicious clients.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.port),
		Handler:      app.newRouter(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// shutdownError channel will receive any errors from the graceful shutdown.
	shutdownError := make(chan error)

	// Start a goroutine to listen for shutdown signals.
	go func() {
		// Create a quit channel that listens for interrupt signals.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		// Give outstanding requests a deadline to finish.
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Calling Shutdown() gracefully shuts down the server without
		// interrupting any active connections. It works by first closing all
		// open listeners, then closing all idle connections, and then waiting
		// indefinitely for connections to return to idle and then shut down.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	// Calling srv.ListenAndServe() will block until the server is shut down.
	// If it returns an error, it's likely a critical one (e.g., port already in use).
	// We specifically ignore http.ErrServerClosed, which is the expected error
	// on a graceful shutdown.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Block until we receive the result of the shutdown.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("server stopped gracefully")

	return nil
}
