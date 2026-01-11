package main

import (
	"context"
	"log"

	"github.com/kamil5b/go-blar/goblar"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product is a sample entity.
type Product struct {
	ID       uint   `gorm:"primaryKey" go-blar:"pk"`
	Name     string
	Price    float64
	Quantity int    `go-blar:"readonly"`
	Secret   string `go-blar:"hidden"`
}

// BeforeCreate validates the product before creation.
func (p *Product) BeforeCreate(ctx context.Context, tx *gorm.DB) error {
	if p.Price < 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

// AfterCreate logs the product creation.
func (p *Product) AfterCreate(ctx context.Context, tx *gorm.DB) error {
	log.Printf("Product created: %s (ID: %d)", p.Name, p.ID)
	return nil
}

// User is another sample entity.
type User struct {
	ID    uint   `gorm:"primaryKey" go-blar:"pk"`
	Name  string
	Email string
}

func main() {
	// Initialize database
	db, err := gorm.Open(sqlite.Open("app.db"))
	if err != nil {
		log.Fatal(err)
	}

	// Create app with configuration
	app := goblar.New(
		goblar.WithDB(db),
		goblar.WithAddress(":8080"),
	)

	// Register models
	if err := app.Register(&Product{}, &User{}); err != nil {
		log.Fatal(err)
	}

	// Start server
	log.Println("Starting server on :8080")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// Generated Routes:
//
// POST   /product       (Create)
// GET    /product       (List all)
// GET    /product/{id}  (Get by ID)
// PUT    /product/{id}  (Update)
// DELETE /product/{id}  (Delete)
//
// POST   /user          (Create)
// GET    /user          (List all)
// GET    /user/{id}     (Get by ID)
// PUT    /user/{id}     (Update)
// DELETE /user/{id}     (Delete)
