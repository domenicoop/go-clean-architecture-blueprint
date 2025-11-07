package service

import (
	"context"
	"fmt"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/apperror"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"

	"github.com/google/uuid"
)

// entityService is a concrete implementation of the EntityService interface.
type entityService struct {
	repo EntityRepository
}

// NewEntityService creates a new entityService instance.
func NewEntityService(repo EntityRepository) EntityService {
	return &entityService{
		repo: repo,
	}
}

// Create creates a new entity.
func (s *entityService) Create(ctx context.Context, entity *domain.Entity) error {
	if entity.Name == "" {
		return apperror.ErrInvalidInput
	}

	entity.ID = uuid.New().String()

	if err := s.repo.Create(ctx, entity); err != nil {
		return fmt.Errorf("service: failed to create entity: %w", err)
	}
	return nil
}

// GetByID retrieves an entity by its ID.
func (s *entityService) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to find entity with id %s: %w", id, err)
	}
	return entity, nil
}

// Update updates an existing entity.
func (s *entityService) Update(ctx context.Context, entity *domain.Entity) error {
	if entity.Name == "" {
		return apperror.ErrInvalidInput
	}

	if err := s.repo.Update(ctx, entity); err != nil {
		return fmt.Errorf("service: failed to update entity with id %s: %w", entity.ID, err)
	}
	return nil
}

// Delete deletes an entity by its ID.
func (s *entityService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: failed to delete entity with id %s: %w", id, err)
	}
	return nil
}

// List retrieves a list of all entities.
func (s *entityService) List(ctx context.Context) ([]*domain.Entity, error) {
	entities, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list entities: %w", err)
	}
	return entities, nil
}
