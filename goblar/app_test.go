package goblar

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestEntity is a sample entity for testing.
type TestEntity struct {
	ID   uint
	Name string
}

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("expected app to be non-nil")
	}

	if app.registry == nil {
		t.Fatal("expected registry to be initialized")
	}

	if app.db != nil {
		t.Fatal("expected db to be nil without WithDB option")
	}
}

func TestNewWithDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	app := New(WithDB(db))
	if app == nil {
		t.Fatal("expected app to be non-nil")
	}

	if app.db != db {
		t.Fatal("expected db to be set")
	}
}

func TestNewWithAddress(t *testing.T) {
	app := New(WithAddress(":3000"))
	if app.cfg.addr != ":3000" {
		t.Fatalf("expected address :3000, got %s", app.cfg.addr)
	}
}

func TestNewWithMultipleOptions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	app := New(
		WithDB(db),
		WithAddress(":9000"),
	)

	if app.db != db {
		t.Fatal("expected db to be set")
	}

	if app.cfg.addr != ":9000" {
		t.Fatal("expected address to be :9000")
	}
}

func TestRegisterWithoutDB(t *testing.T) {
	app := New()
	err := app.Register(&TestEntity{})

	if err == nil {
		t.Fatal("expected error when registering without database")
	}
}

func TestRegisterValid(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	app := New(WithDB(db))
	err = app.Register(&TestEntity{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check registry
	if len(app.registry) == 0 {
		t.Fatal("expected entity in registry")
	}
}

func TestRegisterMultiple(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	type User struct {
		ID   uint
		Name string
	}

	type Product struct {
		ID    uint
		Title string
		Price float64
	}

	app := New(WithDB(db))
	err = app.Register(&User{}, &Product{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(app.registry) != 2 {
		t.Fatalf("expected 2 entities in registry, got %d", len(app.registry))
	}
}
