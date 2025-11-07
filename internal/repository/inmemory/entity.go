package inmemory

import (
	"context"
	"time"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/apperror"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
	"github.com/domenicoop/go-clean-architecture-blueprint/internal/service"
)

// Entity represents a generic domain entity.
type Entity struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// toDomain converts an Entity to a domain.Entity.
func (e *Entity) toDomain() *domain.Entity {
	return &domain.Entity{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// fromDomain converts a domain.Entity to an Entity.
func fromDomain(e *domain.Entity) *Entity {
	return &Entity{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// EntityRepository is a mock implementation of the service.EntityRepository interface.
type EntityRepository struct {
	entities map[string]*Entity
}

// NewEntityRepository creates a new EntityRepository.
func NewEntityRepository() service.EntityRepository {
	return &EntityRepository{
		entities: make(map[string]*Entity),
	}
}

// Create creates a new entity in the mock repository.
func (r *EntityRepository) Create(ctx context.Context, entity *domain.Entity) error {
	if _, exists := r.entities[entity.ID]; exists {
		return apperror.ErrConflict
	}
	storageEntity := fromDomain(entity)
	storageEntity.CreatedAt = time.Now()
	storageEntity.UpdatedAt = time.Now()
	r.entities[entity.ID] = storageEntity
	return nil
}

// FindByID finds an entity by its ID in the mock repository.
func (r *EntityRepository) FindByID(ctx context.Context, id string) (*domain.Entity, error) {
	if entity, exists := r.entities[id]; exists {
		return entity.toDomain(), nil
	}
	return nil, apperror.ErrNotFound
}

// Update updates an entity in the mock repository.
func (r *EntityRepository) Update(ctx context.Context, entity *domain.Entity) error {
	if _, exists := r.entities[entity.ID]; !exists {
		return apperror.ErrNotFound
	}
	storageEntity := fromDomain(entity)
	storageEntity.UpdatedAt = time.Now()
	r.entities[entity.ID] = storageEntity
	return nil
}

// Delete deletes an entity from the mock repository.
func (r *EntityRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.entities[id]; !exists {
		return apperror.ErrNotFound
	}
	delete(r.entities, id)
	return nil
}

// List lists all entities from the mock repository.
func (r *EntityRepository) List(ctx context.Context) ([]*domain.Entity, error) {
	entities := make([]*domain.Entity, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity.toDomain())
	}
	return entities, nil
}
