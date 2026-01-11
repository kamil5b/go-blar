package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kamil5b/go-blar/internal/meta"
)

// Router wraps chi.Router and provides auto-route registration.
type Router struct {
	*chi.Mux
	middlewares []func(http.Handler) http.Handler
}

// New creates a new Router.
func New() *Router {
	return &Router{
		Mux:         chi.NewMux(),
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// AddMiddleware adds a middleware to the router.
func (r *Router) AddMiddleware(m func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, m)
}

// ApplyMiddleware applies all registered middlewares to the router.
func (r *Router) ApplyMiddleware() {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		r.Use(r.middlewares[i])
	}
}

// RegisterEntityRoutes registers REST routes for an entity type.
// This is an internal method called by the app.
func RegisterEntityRoutes(router *Router, meta *meta.EntityMeta, handlers *Handlers) {
	// Construct the resource name from the entity name
	resourceName := toURLPath(meta.Name)

	// POST /resource
	router.Post("/"+resourceName, handlers.CreateHandler(meta))

	// GET /resource
	router.Get("/"+resourceName, handlers.ListHandler(meta))

	// GET /resource/{id}
	router.Get("/"+resourceName+"/{id}", handlers.GetHandler(meta))

	// PUT /resource/{id}
	router.Put("/"+resourceName+"/{id}", handlers.UpdateHandler(meta))

	// DELETE /resource/{id}
	router.Delete("/"+resourceName+"/{id}", handlers.DeleteHandler(meta))
}

// toURLPath converts a CamelCase entity name to a kebab-case URL path.
func toURLPath(s string) string {
	var result []byte
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '-')
		}
		result = append(result, byte(r+'a'-'A'+1)) // Simple lowercase
	}
	return string(result)
}
