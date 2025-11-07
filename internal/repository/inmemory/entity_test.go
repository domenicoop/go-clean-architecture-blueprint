package inmemory

import (
	"context"
	"testing"

	"github.com/domenicoop/go-clean-architecture-blueprint/internal/domain"
)

func TestEntityRepository(t *testing.T) {
	repo := NewEntityRepository()
	ctx := context.Background()

	t.Run("Create and FindByID", func(t *testing.T) {
		entity := &domain.Entity{ID: "1", Name: "Test"}
		err := repo.Create(ctx, entity)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found, err := repo.FindByID(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if found.ID != "1" || found.Name != "Test" {
			t.Errorf("expected entity with ID 1 and Name Test, got %+v", found)
		}
	})

	t.Run("Update", func(t *testing.T) {
		entity := &domain.Entity{ID: "1", Name: "Updated Test"}
		err := repo.Update(ctx, entity)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found, err := repo.FindByID(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if found.Name != "Updated Test" {
			t.Errorf("expected updated name, got %s", found.Name)
		}
	})

	t.Run("List", func(t *testing.T) {
		entities, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(entities) != 1 {
			t.Errorf("expected 1 entity, got %d", len(entities))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.FindByID(ctx, "1")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
