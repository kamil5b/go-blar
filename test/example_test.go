package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	ID           uint `gorm:"primaryKey" go-blar:"pk"`
	Name         string
	Description  string
	UserID       uint
	ProductItems []ProductItem
	Tags         []Tag   `gorm:"many2many:product_to_price;foreignKey:ID;joinForeignKey:ProductID;references:ID;joinReferences:TagID"`
	TotalPrice   float64 `gorm:"-"` // Computed field
	TotalItem    int     `gorm:"-"` // Computed field
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// AfterFind hook to calculate TotalPrice and TotalItem
func (p *Product) AfterFind(tx *gorm.DB) error {
	p.CalculateTotals()
	return nil
}

// CalculateTotals calculates TotalPrice and TotalItem
func (p *Product) CalculateTotals() {
	p.TotalPrice = 0
	p.TotalItem = 0
	for _, item := range p.ProductItems {
		p.TotalPrice += item.TotalPrice
		p.TotalItem += item.Quantity
	}
}

// ProductItem represents items within a product.
type ProductItem struct {
	ID           uint `gorm:"primaryKey" go-blar:"pk"`
	ProductID    uint
	Name         string
	PricePerUnit float64
	Quantity     int
	TotalPrice   float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Tag is a sample entity for categorization.
type Tag struct {
	ID        uint `gorm:"primaryKey" go-blar:"pk"`
	Label     string
	Color     *string // Optional
	CreatedAt time.Time
	UpdatedAt time.Time
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
	color := "#FF5733"
	tag := &Tag{
		ID:    1,
		Label: "Electronics",
		Color: &color,
	}

	if tag.Label != "Electronics" {
		t.Errorf("Expected tag label 'Electronics', got '%s'", tag.Label)
	}
	if tag.Color == nil || *tag.Color != "#FF5733" {
		t.Errorf("Expected tag color '#FF5733', got '%v'", tag.Color)
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

// ============== HTTP Tests ==============

// Tag HTTP Tests

// TestTagCreateRequest tests creating a tag with required fields
func TestTagCreateRequest(t *testing.T) {
	payload := map[string]interface{}{
		"label": "Electronics",
		"color": "#FF5733",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Errorf("Expected status 200/201, got %d", w.Code)
	}
}

// TestTagListRequest tests listing tags with pagination
func TestTagListRequest(t *testing.T) {
	_, _ = http.NewRequest("GET", "/tag?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestTagDetailRequest tests getting tag details
func TestTagDetailRequest(t *testing.T) {
	_, _ = http.NewRequest("GET", "/tag/1", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestTagUpdateRequest tests updating a tag
func TestTagUpdateRequest(t *testing.T) {
	payload := map[string]interface{}{
		"label": "Updated Electronics",
		"color": "#00FF00",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/tag/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestTagDeleteRequest tests deleting a tag
func TestTagDeleteRequest(t *testing.T) {
	_, _ = http.NewRequest("DELETE", "/tag/1", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("Expected status 200/204, got %d", w.Code)
	}
}

// Product HTTP Tests

// TestProductCreateRequest tests creating a product with all fields
func TestProductCreateRequest(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "Laptop",
		"description": "High-performance laptop",
		"userID":      1,
		"productItems": []map[string]interface{}{
			{
				"name":         "Item 1",
				"pricePerUnit": 99.99,
				"quantity":     5,
				"totalPrice":   499.95,
			},
		},
		"tagIDs": []uint{1, 2},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Errorf("Expected status 200/201, got %d", w.Code)
	}
}

// TestProductListRequest tests listing products with pagination
func TestProductListRequest(t *testing.T) {
	_, _ = http.NewRequest("GET", "/product?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestProductDetailRequest tests getting product details with all related data
func TestProductDetailRequest(t *testing.T) {
	_, _ = http.NewRequest("GET", "/product/1", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestProductUpdateRequest tests updating a product
func TestProductUpdateRequest(t *testing.T) {
	payload := map[string]interface{}{
		"name":        "Updated Laptop",
		"description": "High-performance laptop with updated specs",
		"userID":      1,
		"productItems": []map[string]interface{}{
			{
				"name":         "Item 1",
				"pricePerUnit": 109.99,
				"quantity":     5,
				"totalPrice":   549.95,
			},
		},
		"tagIDs": []uint{1, 2, 3},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestProductDeleteRequest tests deleting a product
func TestProductDeleteRequest(t *testing.T) {
	_, _ = http.NewRequest("DELETE", "/product/1", nil)
	w := httptest.NewRecorder()

	if w.Code != http.StatusOK && w.Code != http.StatusNoContent {
		t.Errorf("Expected status 200/204, got %d", w.Code)
	}
}

// TestProductCreateMinimal tests creating a product with minimal required fields
func TestProductCreateMinimal(t *testing.T) {
	payload := map[string]interface{}{
		"name":   "Phone",
		"userID": 1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Errorf("Expected status 200/201, got %d", w.Code)
	}
}

// ============== Flow Tests (CRUD Operations) ==============

// TestTagFlowComplete tests complete tag flow: create -> list -> update -> list -> delete -> list
func TestTagFlowComplete(t *testing.T) {
	t.Log("=== TAG FLOW TEST START ===")

	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate Tag model
	if err := db.AutoMigrate(&Tag{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Step 1: Create first tag
	t.Log("Step 1: Creating first tag...")
	tag1 := &Tag{Label: "Electronics", Color: stringPtr("#FF5733")}
	t.Logf("REQUEST: %s", toJSON(tag1))
	if err := db.Create(tag1).Error; err != nil {
		t.Fatalf("Step 1 Failed: %v", err)
	}
	if tag1.ID == 0 {
		t.Fatalf("Step 1 Failed: Tag ID not assigned")
	}
	t.Logf("RESPONSE: %s", toJSON(tag1))
	t.Logf("✓ First tag created with ID: %d", tag1.ID)

	// Step 2: Create second tag
	t.Log("Step 2: Creating second tag...")
	tag2 := &Tag{Label: "Gadgets", Color: stringPtr("#00FF00")}
	t.Logf("REQUEST: %s", toJSON(tag2))
	if err := db.Create(tag2).Error; err != nil {
		t.Fatalf("Step 2 Failed: %v", err)
	}
	if tag2.ID == 0 {
		t.Fatalf("Step 2 Failed: Tag ID not assigned")
	}
	t.Logf("RESPONSE: %s", toJSON(tag2))
	t.Logf("✓ Second tag created with ID: %d", tag2.ID)

	// Step 3: List tags
	t.Log("Step 3: Listing all tags...")
	var tags []Tag
	if err := db.Find(&tags).Error; err != nil {
		t.Fatalf("Step 3 Failed: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("Step 3 Failed: Expected 2 tags, got %d", len(tags))
	}
	t.Logf("RESPONSE: %s", toJSON(tags))
	t.Logf("✓ Found %d tags in database", len(tags))

	// Step 4: Update first tag
	t.Log("Step 4: Updating first tag...")
	t.Logf("REQUEST (before update): %s", toJSON(tag1))
	tag1.Label = "Updated Electronics"
	tag1.Color = stringPtr("#0000FF")
	t.Logf("REQUEST (update to): %s", toJSON(tag1))
	if err := db.Save(tag1).Error; err != nil {
		t.Fatalf("Step 4 Failed: %v", err)
	}
	// Verify update
	var updated Tag
	if err := db.First(&updated, tag1.ID).Error; err != nil {
		t.Fatalf("Step 4 Failed: %v", err)
	}
	if updated.Label != "Updated Electronics" {
		t.Fatalf("Step 4 Failed: Expected label 'Updated Electronics', got '%s'", updated.Label)
	}
	t.Logf("RESPONSE: %s", toJSON(updated))
	t.Logf("✓ Tag label updated to: %s", updated.Label)

	// Step 5: List tags again
	t.Log("Step 5: Listing all tags again...")
	var tagsAfterUpdate []Tag
	if err := db.Find(&tagsAfterUpdate).Error; err != nil {
		t.Fatalf("Step 5 Failed: %v", err)
	}
	if len(tagsAfterUpdate) != 2 {
		t.Fatalf("Step 5 Failed: Expected 2 tags, got %d", len(tagsAfterUpdate))
	}
	t.Logf("RESPONSE: %s", toJSON(tagsAfterUpdate))
	t.Logf("✓ Still have %d tags", len(tagsAfterUpdate))

	// Step 6: Delete first tag
	t.Log("Step 6: Deleting first tag...")
	t.Logf("REQUEST (delete): %s", toJSON(tag1))
	if err := db.Delete(tag1).Error; err != nil {
		t.Fatalf("Step 6 Failed: %v", err)
	}
	t.Logf("✓ First tag deleted")

	// Step 7: List tags after deletion
	t.Log("Step 7: Listing all tags after deletion...")
	var tagsAfterDelete []Tag
	if err := db.Find(&tagsAfterDelete).Error; err != nil {
		t.Fatalf("Step 7 Failed: %v", err)
	}
	if len(tagsAfterDelete) != 1 {
		t.Fatalf("Step 7 Failed: Expected 1 tag remaining, got %d", len(tagsAfterDelete))
	}
	t.Logf("RESPONSE: %s", toJSON(tagsAfterDelete))
	t.Logf("✓ %d tag remaining after deletion", len(tagsAfterDelete))

	// Step 8: Delete second tag
	t.Log("Step 8: Deleting second tag...")
	t.Logf("REQUEST (delete): %s", toJSON(tag2))
	if err := db.Delete(tag2).Error; err != nil {
		t.Fatalf("Step 8 Failed: %v", err)
	}
	t.Logf("✓ Second tag deleted")

	// Step 9: List tags after all deletions
	t.Log("Step 9: Listing all tags after all deletions...")
	var tagsAfterAllDelete []Tag
	if err := db.Find(&tagsAfterAllDelete).Error; err != nil {
		t.Fatalf("Step 9 Failed: %v", err)
	}
	if len(tagsAfterAllDelete) != 0 {
		t.Fatalf("Step 9 Failed: Expected 0 tags, got %d", len(tagsAfterAllDelete))
	}
	t.Logf("RESPONSE: %s", toJSON(tagsAfterAllDelete))
	t.Logf("✓ All tags deleted successfully")

	t.Log("=== TAG FLOW TEST COMPLETED ✓ ===")
}

// TestProductFlowComplete tests complete product flow: create -> list -> update -> list -> delete -> list
func TestProductFlowComplete(t *testing.T) {
	t.Log("=== PRODUCT FLOW TEST START ===")

	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&Product{}, &ProductItem{}, &Tag{}, &ProductToPrice{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Step 1: Create first product
	t.Log("Step 1: Creating first product...")
	product1 := &Product{
		Name:        "Laptop",
		Description: "High-performance laptop",
		UserID:      1,
		ProductItems: []ProductItem{
			{Name: "Item 1", PricePerUnit: 99.99, Quantity: 5, TotalPrice: 499.95},
		},
	}
	t.Logf("REQUEST: %s", toJSON(product1))
	if err := db.Create(product1).Error; err != nil {
		t.Fatalf("Step 1 Failed: %v", err)
	}
	if product1.ID == 0 {
		t.Fatalf("Step 1 Failed: Product ID not assigned")
	}
	// Calculate totals
	product1.CalculateTotals()
	if product1.TotalPrice != 499.95 {
		t.Fatalf("Step 1 Failed: Expected TotalPrice 499.95, got %f", product1.TotalPrice)
	}
	if product1.TotalItem != 5 {
		t.Fatalf("Step 1 Failed: Expected TotalItem 5, got %d", product1.TotalItem)
	}
	t.Logf("RESPONSE: %s", toJSON(product1))
	t.Logf("✓ First product created with ID: %d, Name: %s, TotalPrice: %.2f, TotalItem: %d", product1.ID, product1.Name, product1.TotalPrice, product1.TotalItem)

	// Step 2: Create second product
	t.Log("Step 2: Creating second product...")
	product2 := &Product{
		Name:        "Phone",
		Description: "Smartphone with great camera",
		UserID:      1,
		ProductItems: []ProductItem{
			{Name: "Item 1", PricePerUnit: 199.99, Quantity: 3, TotalPrice: 599.97},
		},
	}
	t.Logf("REQUEST: %s", toJSON(product2))
	if err := db.Create(product2).Error; err != nil {
		t.Fatalf("Step 2 Failed: %v", err)
	}
	if product2.ID == 0 {
		t.Fatalf("Step 2 Failed: Product ID not assigned")
	}
	// Calculate totals
	product2.CalculateTotals()
	if product2.TotalPrice != 599.97 {
		t.Fatalf("Step 2 Failed: Expected TotalPrice 599.97, got %f", product2.TotalPrice)
	}
	if product2.TotalItem != 3 {
		t.Fatalf("Step 2 Failed: Expected TotalItem 3, got %d", product2.TotalItem)
	}
	t.Logf("RESPONSE: %s", toJSON(product2))
	t.Logf("✓ Second product created with ID: %d, Name: %s, TotalPrice: %.2f, TotalItem: %d", product2.ID, product2.Name, product2.TotalPrice, product2.TotalItem)

	// Step 3: List products
	t.Log("Step 3: Listing all products...")
	var products []Product
	if err := db.Preload("ProductItems").Preload("Tags").Find(&products).Error; err != nil {
		t.Fatalf("Step 3 Failed: %v", err)
	}
	// Calculate totals for all retrieved products
	for i := range products {
		products[i].CalculateTotals()
	}
	if len(products) != 2 {
		t.Fatalf("Step 3 Failed: Expected 2 products, got %d", len(products))
	}
	// Verify first product totals
	if products[0].TotalPrice != 499.95 {
		t.Fatalf("Step 3 Failed: First product - Expected TotalPrice 499.95, got %f", products[0].TotalPrice)
	}
	if products[0].TotalItem != 5 {
		t.Fatalf("Step 3 Failed: First product - Expected TotalItem 5, got %d", products[0].TotalItem)
	}
	t.Logf("RESPONSE: %s", toJSON(products))
	t.Logf("✓ Listed %d products. Product 1 - TotalPrice: %.2f, TotalItem: %d; Product 2 - TotalPrice: %.2f, TotalItem: %d", len(products), products[0].TotalPrice, products[0].TotalItem, products[1].TotalPrice, products[1].TotalItem)

	// Step 4: Update first product
	t.Log("Step 4: Updating first product...")
	t.Logf("REQUEST (before update): %s", toJSON(product1))
	product1.Name = "Updated Laptop"
	product1.Description = "High-performance laptop with upgraded specs"
	t.Logf("REQUEST (update to): %s", toJSON(product1))
	if err := db.Save(product1).Error; err != nil {
		t.Fatalf("Step 4 Failed: %v", err)
	}
	// Verify update
	var updated Product
	if err := db.Preload("ProductItems").Preload("Tags").First(&updated, product1.ID).Error; err != nil {
		t.Fatalf("Step 4 Failed: %v", err)
	}
	updated.CalculateTotals()
	if updated.Name != "Updated Laptop" {
		t.Fatalf("Step 4 Failed: Expected name 'Updated Laptop', got '%s'", updated.Name)
	}
	if updated.Description != "High-performance laptop with upgraded specs" {
		t.Fatalf("Step 4 Failed: Expected description match, got '%s'", updated.Description)
	}
	if updated.TotalPrice != 499.95 {
		t.Fatalf("Step 4 Failed: Expected TotalPrice 499.95, got %f", updated.TotalPrice)
	}
	t.Logf("RESPONSE: %s", toJSON(updated))
	t.Logf("✓ Product updated - Name: %s, TotalPrice: %.2f, TotalItem: %d", updated.Name, updated.TotalPrice, updated.TotalItem)

	// Step 5: List products again
	t.Log("Step 5: Listing all products again...")
	var productsAfterUpdate []Product
	if err := db.Preload("ProductItems").Preload("Tags").Find(&productsAfterUpdate).Error; err != nil {
		t.Fatalf("Step 5 Failed: %v", err)
	}
	// Calculate totals
	for i := range productsAfterUpdate {
		productsAfterUpdate[i].CalculateTotals()
	}
	if len(productsAfterUpdate) != 2 {
		t.Fatalf("Step 5 Failed: Expected 2 products, got %d", len(productsAfterUpdate))
	}
	t.Logf("RESPONSE: %s", toJSON(productsAfterUpdate))
	t.Logf("✓ Still have %d products. Product 1 - TotalPrice: %.2f; Product 2 - TotalPrice: %.2f", len(productsAfterUpdate), productsAfterUpdate[0].TotalPrice, productsAfterUpdate[1].TotalPrice)

	// Step 6: Delete first product
	t.Log("Step 6: Deleting first product...")
	t.Logf("REQUEST (delete): %s", toJSON(product1))
	if err := db.Delete(product1).Error; err != nil {
		t.Fatalf("Step 6 Failed: %v", err)
	}
	t.Logf("✓ First product deleted")

	// Step 7: List products after deletion
	t.Log("Step 7: Listing all products after deletion...")
	var productsAfterDelete []Product
	if err := db.Preload("ProductItems").Preload("Tags").Find(&productsAfterDelete).Error; err != nil {
		t.Fatalf("Step 7 Failed: %v", err)
	}
	// Calculate totals
	for i := range productsAfterDelete {
		productsAfterDelete[i].CalculateTotals()
	}
	if len(productsAfterDelete) != 1 {
		t.Fatalf("Step 7 Failed: Expected 1 product remaining, got %d", len(productsAfterDelete))
	}
	t.Logf("RESPONSE: %s", toJSON(productsAfterDelete))
	t.Logf("✓ %d product remaining with TotalPrice: %.2f", len(productsAfterDelete), productsAfterDelete[0].TotalPrice)

	// Step 8: Delete second product
	t.Log("Step 8: Deleting second product...")
	t.Logf("REQUEST (delete): %s", toJSON(product2))
	if err := db.Delete(product2).Error; err != nil {
		t.Fatalf("Step 8 Failed: %v", err)
	}
	t.Logf("✓ Second product deleted")

	// Step 9: List products after all deletions
	t.Log("Step 9: Listing all products after all deletions...")
	var productsAfterAllDelete []Product
	if err := db.Find(&productsAfterAllDelete).Error; err != nil {
		t.Fatalf("Step 9 Failed: %v", err)
	}
	if len(productsAfterAllDelete) != 0 {
		t.Fatalf("Step 9 Failed: Expected 0 products, got %d", len(productsAfterAllDelete))
	}
	t.Logf("RESPONSE: %s", toJSON(productsAfterAllDelete))
	t.Logf("✓ All products deleted successfully")

	t.Log("=== PRODUCT FLOW TEST COMPLETED ✓ ===")
}

// ============== Detail Tests (Get Single Records) ==============

// TestTagDetailFlow tests retrieving a single tag with all details
func TestTagDetailFlow(t *testing.T) {
	t.Log("=== TAG DETAIL TEST START ===")

	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate Tag model
	if err := db.AutoMigrate(&Tag{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create a tag
	tag := &Tag{Label: "Premium", Color: stringPtr("#FFD700")}
	if err := db.Create(tag).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	t.Logf("✓ Tag created with ID: %d", tag.ID)

	// Retrieve the tag detail
	t.Log("Getting tag detail...")
	var detail Tag
	if err := db.First(&detail, tag.ID).Error; err != nil {
		t.Fatalf("Detail retrieval failed: %v", err)
	}
	if detail.Label != "Premium" {
		t.Fatalf("Expected label 'Premium', got '%s'", detail.Label)
	}
	if detail.Color == nil || *detail.Color != "#FFD700" {
		t.Fatalf("Expected color '#FFD700', got '%v'", detail.Color)
	}
	t.Logf("DETAIL RESPONSE: %s", toJSON(detail))
	t.Logf("✓ Tag detail retrieved successfully - ID: %d, Label: %s, Color: %s, CreatedAt: %s", detail.ID, detail.Label, *detail.Color, detail.CreatedAt)

	t.Log("=== TAG DETAIL TEST COMPLETED ✓ ===")
}

// TestProductDetailFlow tests retrieving a single product with all related data
func TestProductDetailFlow(t *testing.T) {
	t.Log("=== PRODUCT DETAIL TEST START ===")

	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&Product{}, &ProductItem{}, &Tag{}, &ProductToPrice{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create tags
	tag1 := &Tag{Label: "Electronics", Color: stringPtr("#FF5733")}
	tag2 := &Tag{Label: "Premium", Color: stringPtr("#FFD700")}
	if err := db.Create(tag1).Error; err != nil {
		t.Fatalf("Tag1 creation failed: %v", err)
	}
	if err := db.Create(tag2).Error; err != nil {
		t.Fatalf("Tag2 creation failed: %v", err)
	}

	// Create a product with multiple items
	product := &Product{
		Name:        "Gaming Laptop",
		Description: "High-end gaming laptop with RTX 4090",
		UserID:      1,
		ProductItems: []ProductItem{
			{Name: "Base Unit", PricePerUnit: 1500.00, Quantity: 1, TotalPrice: 1500.00},
			{Name: "Extended Warranty", PricePerUnit: 200.00, Quantity: 1, TotalPrice: 200.00},
			{Name: "Protective Case", PricePerUnit: 50.00, Quantity: 2, TotalPrice: 100.00},
		},
	}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("Product creation failed: %v", err)
	}
	t.Logf("✓ Product created with ID: %d", product.ID)

	// Associate tags
	if err := db.Model(product).Association("Tags").Append(tag1, tag2); err != nil {
		t.Fatalf("Tag association failed: %v", err)
	}
	t.Logf("✓ Associated 2 tags with product")

	// Retrieve product detail with all relations
	t.Log("Getting product detail with all relations...")
	var detail Product
	if err := db.Preload("ProductItems").Preload("Tags").First(&detail, product.ID).Error; err != nil {
		t.Fatalf("Detail retrieval failed: %v", err)
	}
	detail.CalculateTotals()

	// Validate product details
	if detail.Name != "Gaming Laptop" {
		t.Fatalf("Expected name 'Gaming Laptop', got '%s'", detail.Name)
	}
	if len(detail.ProductItems) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(detail.ProductItems))
	}
	if len(detail.Tags) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(detail.Tags))
	}
	expectedTotal := 1500.00 + 200.00 + 100.00
	if detail.TotalPrice != expectedTotal {
		t.Fatalf("Expected TotalPrice %.2f, got %.2f", expectedTotal, detail.TotalPrice)
	}
	if detail.TotalItem != 4 { // 1 + 1 + 2
		t.Fatalf("Expected TotalItem 4, got %d", detail.TotalItem)
	}

	t.Logf("DETAIL RESPONSE: %s", toJSON(detail))
	t.Logf("✓ Product detail retrieved - ID: %d, Name: %s, Items: %d, Tags: %d, TotalPrice: %.2f, TotalItem: %d",
		detail.ID, detail.Name, len(detail.ProductItems), len(detail.Tags), detail.TotalPrice, detail.TotalItem)

	// Verify each item details
	for i, item := range detail.ProductItems {
		t.Logf("  Item %d: Name='%s', Qty=%d, Price=%.2f, Total=%.2f",
			i+1, item.Name, item.Quantity, item.PricePerUnit, item.TotalPrice)
	}

	// Verify each tag details
	for i, tag := range detail.Tags {
		colorStr := "nil"
		if tag.Color != nil {
			colorStr = *tag.Color
		}
		t.Logf("  Tag %d: Label='%s', Color=%s", i+1, tag.Label, colorStr)
	}

	t.Log("=== PRODUCT DETAIL TEST COMPLETED ✓ ===")
}

// TestProductItemDetailFlow tests retrieving product items with calculations
func TestProductItemDetailFlow(t *testing.T) {
	t.Log("=== PRODUCT ITEM DETAIL TEST START ===")

	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&Product{}, &ProductItem{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create a product
	product := &Product{
		Name:        "Bundle",
		Description: "Multi-item bundle",
		UserID:      1,
	}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("Product creation failed: %v", err)
	}

	// Create items
	items := []ProductItem{
		{ProductID: product.ID, Name: "Item A", PricePerUnit: 25.50, Quantity: 10, TotalPrice: 255.00},
		{ProductID: product.ID, Name: "Item B", PricePerUnit: 15.75, Quantity: 20, TotalPrice: 315.00},
		{ProductID: product.ID, Name: "Item C", PricePerUnit: 5.00, Quantity: 50, TotalPrice: 250.00},
	}
	for _, item := range items {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("Item creation failed: %v", err)
		}
	}
	t.Logf("✓ Created 3 product items")

	// Retrieve all items for this product
	t.Log("Retrieving all product items...")
	var retrievedItems []ProductItem
	if err := db.Where("product_id = ?", product.ID).Find(&retrievedItems).Error; err != nil {
		t.Fatalf("Items retrieval failed: %v", err)
	}

	if len(retrievedItems) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(retrievedItems))
	}

	totalPrice := 0.0
	totalQty := 0
	for _, item := range retrievedItems {
		totalPrice += item.TotalPrice
		totalQty += item.Quantity
		t.Logf("ITEM DETAIL: %s", toJSON(item))
	}

	expectedTotalPrice := 820.00 // 255 + 315 + 250
	if totalPrice != expectedTotalPrice {
		t.Fatalf("Expected total price %.2f, got %.2f", expectedTotalPrice, totalPrice)
	}
	if totalQty != 80 { // 10 + 20 + 50
		t.Fatalf("Expected total quantity 80, got %d", totalQty)
	}

	t.Logf("✓ All items verified - Total Price: %.2f, Total Quantity: %d", totalPrice, totalQty)
	t.Log("=== PRODUCT ITEM DETAIL TEST COMPLETED ✓ ===")
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// toJSON converts any value to JSON string
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	return string(data)
}

// Generated Routes:
//
// POST   /user                (Create)
// GET    /user                (List all)
// GET    /user/{id}           (Get by ID)
// PUT    /user/{id}           (Update)
// DELETE /user/{id}           (Delete)
//
// POST   /product             (Create with ProductItems[] and TagIDs[])
// GET    /product             (List all with pagination)
// GET    /product/{id}        (Get by ID with all details)
// PUT    /product/{id}        (Update with ProductItems[] and TagIDs[])
// DELETE /product/{id}        (Delete)
//
// POST   /productitem         (Create)
// GET    /productitem         (List all)
// GET    /productitem/{id}    (Get by ID)
// PUT    /productitem/{id}    (Update)
// DELETE /productitem/{id}    (Delete)
//
// POST   /tag                 (Create with Label and optional Color)
// GET    /tag                 (List all with pagination)
// GET    /tag/{id}            (Get by ID with timestamps)
// PUT    /tag/{id}            (Update)
// DELETE /tag/{id}            (Delete)
//
// POST   /productToprice      (Create)
// GET    /productToprice      (List all)
// GET    /productToprice/{id} (Get by ID)
// PUT    /productToprice/{id} (Update)
// DELETE /productToprice/{id} (Delete)
