package goblar

import (
	"net/http"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWithDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	cfg := newConfig()
	opt := WithDB(db)
	opt(cfg)

	if cfg.db != db {
		t.Fatal("expected db to be set")
	}
}

func TestWithAddress(t *testing.T) {
	cfg := newConfig()
	opt := WithAddress(":5000")
	opt(cfg)

	if cfg.addr != ":5000" {
		t.Fatalf("expected address :5000, got %s", cfg.addr)
	}
}

func TestWithMiddleware(t *testing.T) {
	cfg := newConfig()

	middleware1 := func(h http.Handler) http.Handler { return h }
	middleware2 := func(h http.Handler) http.Handler { return h }

	WithMiddleware(middleware1)(cfg)
	WithMiddleware(middleware2)(cfg)

	if len(cfg.middleware) != 2 {
		t.Fatalf("expected 2 middlewares, got %d", len(cfg.middleware))
	}
}

func TestConfig_Apply(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatal(err)
	}

	cfg := newConfig()
	cfg.apply(
		WithDB(db),
		WithAddress(":7000"),
	)

	if cfg.db != db {
		t.Fatal("expected db to be set")
	}

	if cfg.addr != ":7000" {
		t.Fatalf("expected address :7000, got %s", cfg.addr)
	}
}

func TestNewConfig_Defaults(t *testing.T) {
	cfg := newConfig()

	if cfg.addr != ":8080" {
		t.Fatalf("expected default address :8080, got %s", cfg.addr)
	}

	if cfg.db != nil {
		t.Fatal("expected db to be nil by default")
	}

	if len(cfg.middleware) != 0 {
		t.Fatal("expected no middleware by default")
	}
}
