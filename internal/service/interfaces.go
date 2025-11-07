package service

import (
	"context"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
)

// EntityRepository defines the contract for data persistence operations for Entities.
type EntityRepository interface {
	Create(ctx context.Context, entity *domain.Entity) error
	FindByID(ctx context.Context, id string) (*domain.Entity, error)
	Update(ctx context.Context, entity *domain.Entity) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Entity, error)
}

// EntityService defines the contract for business logic operations for Entities.
type EntityService interface {
	Create(ctx context.Context, entity *domain.Entity) error
	GetByID(ctx context.Context, id string) (*domain.Entity, error)
	Update(ctx context.Context, entity *domain.Entity) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Entity, error)
}
