package goblar

import (
	"fmt"
	"net/http"

	"github.com/kamil5b/go-blar/internal/meta"
	"gorm.io/gorm"
)

// App is the main runtime for go-blar.
// It holds minimal public state and hides internal implementation details.
type App struct {
	db       *gorm.DB
	router   http.Handler
	cfg      *config
	registry map[string]any
}

// New creates a new App with the given options.
func New(opts ...Option) *App {
	cfg := newConfig()
	cfg.apply(opts...)

	app := &App{
		db:       cfg.db,
		cfg:      cfg,
		registry: make(map[string]any),
	}

	return app
}

// Register registers one or more model structs with the app.
// Models must be valid GORM entities.
func (a *App) Register(models ...any) error {
	if a.db == nil {
		return fmt.Errorf("database not configured: use WithDB option")
	}

	// Parse and validate each model
	for _, model := range models {
		entityMeta, err := meta.Parse(model)
		if err != nil {
			return fmt.Errorf("failed to parse model %T: %w", model, err)
		}

		// Auto-migrate the entity with GORM
		if err := a.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}

		// Store in registry
		a.registry[entityMeta.Name] = entityMeta
	}

	return nil
}

// Start starts the HTTP server and serves the auto-generated routes.
func (a *App) Start() error {
	if a.cfg.addr == "" {
		a.cfg.addr = ":8080"
	}

	// Build HTTP router (stubbed for now - will be implemented in internal/http)
	// TODO: Build router with registered models
	// TODO: Apply middleware
	// TODO: Register routes for all entities

	server := &http.Server{
		Addr:    a.cfg.addr,
		Handler: a.router,
	}

	return server.ListenAndServe()
}
