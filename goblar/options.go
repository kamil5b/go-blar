package goblar

import (
	"net/http"

	"gorm.io/gorm"
)

// config holds the internal configuration for the app.
type config struct {
	db         *gorm.DB
	addr       string
	middleware []func(http.Handler) http.Handler
}

// Option is a functional option for configuring the App.
type Option func(*config)

// WithDB sets the GORM database connection.
func WithDB(db *gorm.DB) Option {
	return func(c *config) {
		c.db = db
	}
}

// WithAddress sets the HTTP server address (e.g., ":8080").
func WithAddress(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

// WithMiddleware adds an HTTP middleware to the router.
func WithMiddleware(m func(http.Handler) http.Handler) Option {
	return func(c *config) {
		c.middleware = append(c.middleware, m)
	}
}

// newConfig creates a new config with sensible defaults.
func newConfig() *config {
	return &config{
		addr: ":8080",
	}
}

// apply applies all options to the config.
func (c *config) apply(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}
