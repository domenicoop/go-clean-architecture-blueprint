package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/apperror"
)

// handleError is a centralized error handler for the HTTP layer.
// It maps application-specific errors to HTTP status codes and logs unknown errors.
func (h *EntityHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	// Use errors.Is to check for known error types.
	switch {
	case errors.Is(err, apperror.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, apperror.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, apperror.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		// For unknown errors, log the full error and return a generic
		// 500 Internal Server Error to the client.
		h.logger.Error("internal server error", "error", err.Error(), "method", r.Method, "url", r.URL.String())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// writeJSON is a helper for writing JSON responses.
// It marshals the data to JSON first, handling potential errors before writing to the response.
func (h *EntityHandler) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	// If there's no data to send, just write the status code.
	if data == nil {
		w.WriteHeader(status)
		return
	}

	// Marshal the data to JSON. If this fails, it's a server-side problem.
	js, err := json.Marshal(data)
	if err != nil {
		// Log the underlying error and send a generic 500 response.
		err = fmt.Errorf("failed to marshal JSON response: %w", err)
		h.handleError(w, r, err)
		return
	}

	// Add a newline to the JSON output for better readability in terminals.
	js = append(js, '\n')

	// Set the content type and write the status code and response body.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		// If writing fails, the response has already started, so we can't send
		// a new error. We just log it.
		h.logger.Error("failed to write response", "error", err.Error(), "method", r.Method, "url", r.URL.String())
	}
}
