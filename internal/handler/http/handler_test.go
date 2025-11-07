package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
	"github.com/go-chi/chi/v5"
)

// mockEntityService is a mock implementation of the EntityService interface.
type mockEntityService struct {
	CreateFunc  func(ctx context.Context, entity *domain.Entity) error
	GetByIDFunc func(ctx context.Context, id string) (*domain.Entity, error)
	UpdateFunc  func(ctx context.Context, entity *domain.Entity) error
	DeleteFunc  func(ctx context.Context, id string) error
	ListFunc    func(ctx context.Context) ([]*domain.Entity, error)
}

func (m *mockEntityService) Create(ctx context.Context, entity *domain.Entity) error {
	return m.CreateFunc(ctx, entity)
}

func (m *mockEntityService) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *mockEntityService) Update(ctx context.Context, entity *domain.Entity) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *mockEntityService) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *mockEntityService) List(ctx context.Context) ([]*domain.Entity, error) {
	return m.ListFunc(ctx)
}

func TestEntityHandler(t *testing.T) {
	mockService := &mockEntityService{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewEntityHandler(mockService, logger)

	t.Run("CreateEntity", func(t *testing.T) {
		mockService.CreateFunc = func(ctx context.Context, entity *domain.Entity) error {
			return nil
		}

		body, _ := json.Marshal(CreateEntityRequest{Name: "Test"})
		req := httptest.NewRequest("POST", "/entities", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.CreateEntity(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("GetEntity", func(t *testing.T) {
		expectedEntity := &domain.Entity{ID: "1", Name: "Test"}
		mockService.GetByIDFunc = func(ctx context.Context, id string) (*domain.Entity, error) {
			return expectedEntity, nil
		}

		req := httptest.NewRequest("GET", "/entities/1", nil)
		rr := httptest.NewRecorder()

		// Add chi context to the request
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.GetEntity(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var response EntityResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if response.ID != "1" {
			t.Errorf("expected entity with ID 1, got %s", response.ID)
		}
	})
}
