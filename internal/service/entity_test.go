package service

import (
	"context"
	"testing"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
	"github.com/google/uuid"
)

// mockEntityRepository is a mock implementation of the EntityRepository interface.
type mockEntityRepository struct {
	CreateFunc   func(ctx context.Context, entity *domain.Entity) error
	FindByIDFunc func(ctx context.Context, id string) (*domain.Entity, error)
	UpdateFunc   func(ctx context.Context, entity *domain.Entity) error
	DeleteFunc   func(ctx context.Context, id string) error
	ListFunc     func(ctx context.Context) ([]*domain.Entity, error)
}

func (m *mockEntityRepository) Create(ctx context.Context, entity *domain.Entity) error {
	return m.CreateFunc(ctx, entity)
}

func (m *mockEntityRepository) FindByID(ctx context.Context, id string) (*domain.Entity, error) {
	return m.FindByIDFunc(ctx, id)
}

func (m *mockEntityRepository) Update(ctx context.Context, entity *domain.Entity) error {
	return m.UpdateFunc(ctx, entity)
}

func (m *mockEntityRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *mockEntityRepository) List(ctx context.Context) ([]*domain.Entity, error) {
	return m.ListFunc(ctx)
}

func TestEntityService(t *testing.T) {
	mockRepo := &mockEntityRepository{}
	service := NewEntityService(mockRepo)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		entity := &domain.Entity{Name: "Test"}
		mockRepo.CreateFunc = func(ctx context.Context, e *domain.Entity) error {
			if e.ID == "" {
				t.Error("expected ID to be set")
			}
			_, err := uuid.Parse(e.ID)
			if err != nil {
				t.Errorf("expected ID to be a valid UUID, got %s", e.ID)
			}
			return nil
		}

		err := service.Create(ctx, entity)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		expectedEntity := &domain.Entity{ID: "1", Name: "Test"}
		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*domain.Entity, error) {
			return expectedEntity, nil
		}

		entity, err := service.GetByID(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if entity.ID != "1" {
			t.Errorf("expected entity with ID 1, got %s", entity.ID)
		}
	})
}
