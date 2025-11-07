package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/apperror"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/service"
	"github.com/go-chi/chi/v5"
)

// EntityHandler is responsible for handling HTTP requests related to entities.
type EntityHandler struct {
	service service.EntityService
	logger  *slog.Logger
}

// NewEntityHandler creates a new EntityHandler.
func NewEntityHandler(service service.EntityService, logger *slog.Logger) *EntityHandler {
	return &EntityHandler{
		service: service,
		logger:  logger,
	}
}

// CreateEntityRequest defines the request body for creating an entity.
type CreateEntityRequest struct {
	Name string `json:"name"`
}

// toDomain converts a CreateEntityRequest to a domain.Entity.
func (r *CreateEntityRequest) toDomain() *domain.Entity {
	return &domain.Entity{
		Name: r.Name,
	}
}

// UpdateEntityRequest defines the request body for updating an entity.
type UpdateEntityRequest struct {
	Name string `json:"name"`
}

// toDomain converts an UpdateEntityRequest to a domain.Entity.
func (r *UpdateEntityRequest) toDomain(id string) *domain.Entity {
	return &domain.Entity{
		ID:   id,
		Name: r.Name,
	}
}

// EntityResponse defines the response body for an entity.
type EntityResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// fromDomain converts a domain.Entity to an EntityResponse.
func fromDomain(entity *domain.Entity) *EntityResponse {
	return &EntityResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// CreateEntity handles the POST /entities endpoint.
func (h *EntityHandler) CreateEntity(w http.ResponseWriter, r *http.Request) {
	var req CreateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperror.ErrInvalidInput)
		return
	}

	entity := req.toDomain()
	if err := h.service.Create(r.Context(), entity); err != nil {
		h.handleError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusCreated, nil)
}

// GetEntity handles the GET /entities/{id} endpoint.
func (h *EntityHandler) GetEntity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entity, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusOK, fromDomain(entity))
}

// UpdateEntity handles the PUT /entities/{id} endpoint.
func (h *EntityHandler) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperror.ErrInvalidInput)
		return
	}

	entity := req.toDomain(id)
	if err := h.service.Update(r.Context(), entity); err != nil {
		h.handleError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusOK, nil)
}

// DeleteEntity handles the DELETE /entities/{id} endpoint.
func (h *EntityHandler) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleError(w, r, err)
		return
	}

	h.writeJSON(w, r, http.StatusOK, nil)
}

// ListEntities handles the GET /entities endpoint.
func (h *EntityHandler) ListEntities(w http.ResponseWriter, r *http.Request) {
	entities, err := h.service.List(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response := make([]*EntityResponse, len(entities))
	for i, entity := range entities {
		response[i] = fromDomain(entity)
	}

	h.writeJSON(w, r, http.StatusOK, response)
}
