package test

import (
	"testing"

	"github.com/kamil5b/go-blar/goblar"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User is a sample entity.
type User struct {
	ID       uint `gorm:"primaryKey" go-blar:"pk"`
	Name     string
	Email    string
	Password string
}

// Product is a sample entity.
type Product struct {
	ID          uint `gorm:"primaryKey" go-blar:"pk"`
	Name        string
	Description string
	UserID      uint
}

// ProductItem represents items within a product.
type ProductItem struct {
	ID           uint `gorm:"primaryKey" go-blar:"pk"`
	ProductID    uint
	Name         string
	PricePerUnit float64
	Quantity     int
	TotalPrice   float64
}

// Tag is a sample entity for categorization.
type Tag struct {
	ID    uint `gorm:"primaryKey" go-blar:"pk"`
	Label string
	Color string
}

// ProductToPrice is a many-to-many relationship between Product and Tag.
type ProductToPrice struct {
	ProductID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}

func TestAppInitialization(t *testing.T) {
	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create app with configuration
	app := goblar.New(
		goblar.WithDB(db),
		goblar.WithAddress(":8080"),
	)

	// Register models
	if err := app.Register(&User{}, &Product{}, &ProductItem{}, &Tag{}, &ProductToPrice{}); err != nil {
		t.Fatalf("Failed to register models: %v", err)
	}

	t.Log("App initialized successfully with all models registered")
}

func TestUserModel(t *testing.T) {
	user := &User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword",
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got '%s'", user.Name)
	}
}

func TestProductModel(t *testing.T) {
	product := &Product{
		ID:          1,
		Name:        "Laptop",
		Description: "High-performance laptop",
		UserID:      1,
	}

	if product.Name != "Laptop" {
		t.Errorf("Expected product name 'Laptop', got '%s'", product.Name)
	}
}

func TestProductItemModel(t *testing.T) {
	item := &ProductItem{
		ID:           1,
		ProductID:    1,
		Name:         "Item 1",
		PricePerUnit: 99.99,
		Quantity:     5,
		TotalPrice:   499.95,
	}

	if item.TotalPrice != 499.95 {
		t.Errorf("Expected total price 499.95, got %f", item.TotalPrice)
	}
}

func TestTagModel(t *testing.T) {
	tag := &Tag{
		ID:    1,
		Label: "Electronics",
		Color: "#FF5733",
	}

	if tag.Label != "Electronics" {
		t.Errorf("Expected tag label 'Electronics', got '%s'", tag.Label)
	}
}

func TestProductToPriceModel(t *testing.T) {
	p2p := &ProductToPrice{
		ProductID: 1,
		TagID:     1,
	}

	if p2p.ProductID != 1 || p2p.TagID != 1 {
		t.Errorf("Expected ProductID=1 and TagID=1, got ProductID=%d TagID=%d", p2p.ProductID, p2p.TagID)
	}
}

// Generated Routes:
//
// POST   /user                (Create)
// GET    /user                (List all)
// GET    /user/{id}           (Get by ID)
// PUT    /user/{id}           (Update)
// DELETE /user/{id}           (Delete)
//
// POST   /product             (Create)
// GET    /product             (List all)
// GET    /product/{id}        (Get by ID)
// PUT    /product/{id}        (Update)
// DELETE /product/{id}        (Delete)
//
// POST   /productitem         (Create)
// GET    /productitem         (List all)
// GET    /productitem/{id}    (Get by ID)
// PUT    /productitem/{id}    (Update)
// DELETE /productitem/{id}    (Delete)
//
// POST   /tag                 (Create)
// GET    /tag                 (List all)
// GET    /tag/{id}            (Get by ID)
// PUT    /tag/{id}            (Update)
// DELETE /tag/{id}            (Delete)
//
// POST   /productToprice      (Create)
// GET    /productToprice      (List all)
// GET    /productToprice/{id} (Get by ID)
// PUT    /productToprice/{id} (Update)
// DELETE /productToprice/{id} (Delete)
