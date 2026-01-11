package repo

import (
	"context"
	"testing"

	"github.com/kamil5b/go-blar/internal/meta"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestUser struct {
	ID   uint
	Name string
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("failed to migrate test schema: %v", err)
	}

	return db
}

func TestNewRepository(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)

	if repo == nil {
		t.Fatal("expected repository to be non-nil")
	}

	if repo.db != db {
		t.Fatal("expected db to be set")
	}

	if repo.meta != entityMeta {
		t.Fatal("expected meta to be set")
	}
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	user := &TestUser{Name: "John"}
	err = repo.Create(ctx, user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected ID to be assigned")
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	// Create a user first
	user := &TestUser{Name: "Jane"}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatal(err)
	}

	// Retrieve by ID
	retrieved, err := repo.GetByID(ctx, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected user to be retrieved")
	}

	if retrieved.Name != "Jane" {
		t.Fatalf("expected name Jane, got %s", retrieved.Name)
	}
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	// Create multiple users
	users := []*TestUser{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}

	for _, u := range users {
		if err := repo.Create(ctx, u); err != nil {
			t.Fatal(err)
		}
	}

	// Retrieve all
	all, err := repo.GetAll(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(all) != 3 {
		t.Fatalf("expected 3 users, got %d", len(all))
	}
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	// Create a user
	user := &TestUser{Name: "Original"}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatal(err)
	}

	// Update the user
	user.Name = "Updated"
	err = repo.Update(ctx, user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if retrieved.Name != "Updated" {
		t.Fatalf("expected name Updated, got %s", retrieved.Name)
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	// Create a user
	user := &TestUser{Name: "ToDelete"}
	if err := repo.Create(ctx, user); err != nil {
		t.Fatal(err)
	}

	userID := user.ID

	// Delete the user
	err = repo.Delete(ctx, userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify deletion
	retrieved, err := repo.GetByID(ctx, userID)
	if err == nil {
		t.Fatal("expected user to be deleted")
	}

	if retrieved != nil {
		t.Fatal("expected nil after deletion")
	}
}

func TestCount(t *testing.T) {
	db := setupTestDB(t)
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](db, entityMeta)
	ctx := context.Background()

	// Create some users
	for i := 0; i < 5; i++ {
		user := &TestUser{Name: "User"}
		if err := repo.Create(ctx, user); err != nil {
			t.Fatal(err)
		}
	}

	// Count
	count, err := repo.Count(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 5 {
		t.Fatalf("expected count 5, got %d", count)
	}
}

func TestCreateWithoutDB(t *testing.T) {
	meta.ClearRegistry()

	entityMeta, err := meta.Parse(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[TestUser](nil, entityMeta)
	ctx := context.Background()

	user := &TestUser{Name: "John"}
	err = repo.Create(ctx, user)

	if err == nil {
		t.Fatal("expected error when creating without database")
	}
}
