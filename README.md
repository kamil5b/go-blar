# go-blar

**A minimal, embeddable, struct-driven ORM framework for Go.**

go-blar provides a single API surface that auto-generates REST endpoints from GORM models, complete with lifecycle hooks, relationships, and aggregates—all without reflection in public code or code generation.

## Design Philosophy

- **One entrypoint**: `goblar.Run()` or `goblar.New()`
- **Struct-only**: No code generation, no interfaces, no reflection in public APIs
- **Minimal**: All reflection and internals live in `internal/`
- **Embeddable**: Works in any Go application, produces one binary
- **Extensible**: Hooks are the intentional escape hatch

---

## Installation

```bash
go get github.com/kamil5b/go-blar
```

---

## Quick Start

### The simplest way: `Run()`

```go
package main

import (
	"github.com/kamil5b/go-blar/goblar"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price float64
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic(err)
	}

	app := goblar.New(goblar.WithDB(db), goblar.WithAddress(":8080"))
	app.Register(&Product{})
	app.Start()
}
```

Automatically generates:
```
POST   /product
GET    /product
GET    /product/{id}
PUT    /product/{id}
DELETE /product/{id}
```

---

## API Reference

### `goblar.Run(models ...any) error`

Zero-brain entrypoint. Creates an app, registers models, and starts the server on `:8080`.

```go
goblar.Run(&Product{}, &User{})
```

### `goblar.New(opts ...Option) *App`

Create an app with options.

```go
app := goblar.New(
	goblar.WithDB(db),
	goblar.WithAddress(":3000"),
	goblar.WithMiddleware(myMiddleware),
)
```

### `app.Register(models ...any) error`

Register GORM entities.

```go
app.Register(&Product{}, &User{})
```

### `app.Start() error`

Start the HTTP server.

```go
if err := app.Start(); err != nil {
	log.Fatal(err)
}
```

---

## Configuration Options

### `WithDB(db *gorm.DB)`

Set the database connection.

```go
app := goblar.New(goblar.WithDB(db))
```

### `WithAddress(addr string)`

Set the server address (default: `:8080`).

```go
app := goblar.New(goblar.WithAddress(":3000"))
```

### `WithMiddleware(m func(http.Handler) http.Handler)`

Add HTTP middleware.

```go
app := goblar.New(
	goblar.WithMiddleware(cors.Default().Handler),
)
```

---

## Struct Tags (DSL)

Use `go-blar` tags to configure entities and fields.

### Entity Tags

(Currently parsed but reserved for future use.)

### Field Tags

```go
type Product struct {
	// Primary key
	ID uint `go-blar:"pk" gorm:"primaryKey"`

	// Basic fields
	Name     string
	Price    float64
	Quantity int

	// Hidden from API
	Secret string `go-blar:"hidden"`

	// Read-only (ignored on Create/Update)
	CreatedAt time.Time `go-blar:"readonly"`

	// Foreign key
	UserID uint   `go-blar:"fk:User"`
	User   *User  `gorm:"foreignKey:UserID"`

	// Many-to-many
	Tags []Tag `go-blar:"m2m:product_tags" gorm:"many2many:product_tags"`

	// Aggregates
	ItemCount int `go-blar:"count:Items"`
	ItemSum   float64 `go-blar:"sum:Items.Price*Items.Quantity"`

	// List endpoint
	Items []Item `go-blar:"list"`
}
```

---

## Lifecycle Hooks

Implement optional interfaces to hook into entity lifecycle:

```go
func (p *Product) BeforeCreate(ctx context.Context, tx *gorm.DB) error {
	// Validate before creation
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func (p *Product) AfterCreate(ctx context.Context, tx *gorm.DB) error {
	// Log, cache invalidation, etc.
	return nil
}

func (p *Product) BeforeUpdate(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (p *Product) AfterUpdate(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (p *Product) BeforeDelete(ctx context.Context, tx *gorm.DB) error {
	return nil
}

func (p *Product) AfterDelete(ctx context.Context, tx *gorm.DB) error {
	return nil
}
```

All hook interfaces are optional. Just implement what you need.

---

## Hook Interfaces (Public)

```go
package goblar

type BeforeCreate interface {
	BeforeCreate(ctx context.Context, tx *gorm.DB) error
}

type AfterCreate interface {
	AfterCreate(ctx context.Context, tx *gorm.DB) error
}

type BeforeUpdate interface {
	BeforeUpdate(ctx context.Context, tx *gorm.DB) error
}

type AfterUpdate interface {
	AfterUpdate(ctx context.Context, tx *gorm.DB) error
}

type BeforeDelete interface {
	BeforeDelete(ctx context.Context, tx *gorm.DB) error
}

type AfterDelete interface {
	AfterDelete(ctx context.Context, tx *gorm.DB) error
}
```

---

## Project Structure

```
go-blar/
├── go.mod                          // Module definition
├── goblar/                         // PUBLIC API
│   ├── app.go                      // App, New()
│   ├── hooks.go                    // Hook interfaces
│   ├── options.go                  // Option pattern
│   └── run.go                      // Run()
│
└── internal/                       // HIDDEN
    ├── meta/
    │   ├── parse.go                // Struct parsing & tag extraction
    │   └── entity.go               // EntityMeta, FieldMeta structures
    │
    ├── repo/
    │   ├── repository.go           // Generic Repository[T]
    │   └── save.go                 // Save/create operations
    │
    ├── aggregate/
    │   └── compute.go              // Aggregate computation (count, sum, etc)
    │
    ├── hooks/
    │   └── hooks.go                // Hook invocation helpers
    │
    ├── http/
    │   ├── router.go               // Router wrapper, route registration
    │   └── handlers.go             // Generic HTTP handlers
    │
    └── util/
        └── reflect.go              // Reflection helpers (hidden)
```

---

## Implementation Status

### ✅ Complete
- Public API (app, hooks, options, run)
- Metadata parsing
- Hook interfaces & execution
- Generic repository (CRUD)
- HTTP router scaffolding
- Struct tag parsing

### 🔄 In Progress
- HTTP handlers (CRUD handlers)
- Route registration
- Middleware application

### 📋 Future
- Aggregate computation
- Foreign key loading
- Many-to-many loading
- Nested struct handling
- Validation framework
- GraphQL layer (optional)

---

## Versioning

- **v0.x**: DSL may change, API may shift
- **v1.0**: Tags frozen, hooks & options can expand without breaking

---

## Contributing

This is a research project exploring minimal ORMs in Go. PRs welcome.

---

## License

MIT
