package meta

import (
	"reflect"
	"testing"
)

// TestParseBasicStruct tests parsing a basic struct.
func TestParseBasicStruct(t *testing.T) {
	type Product struct {
		ID    uint   `gorm:"primaryKey" go-blar:"pk"`
		Name  string
		Price float64
	}

	ClearRegistry()
	meta, err := Parse(&Product{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if meta == nil {
		t.Fatal("expected meta to be non-nil")
	}

	if meta.Name != "Product" {
		t.Fatalf("expected name Product, got %s", meta.Name)
	}

	if meta.Type != reflect.TypeOf(Product{}) {
		t.Fatalf("expected correct type")
	}

	if len(meta.Fields) == 0 {
		t.Fatal("expected fields to be parsed")
	}
}

func TestParseWithPrimaryKey(t *testing.T) {
	type Entity struct {
		ID   uint   `gorm:"primaryKey" go-blar:"pk"`
		Name string
	}

	ClearRegistry()
	meta, err := Parse(&Entity{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if meta.PKField == nil {
		t.Fatal("expected primary key field to be detected")
	}

	if meta.PKField.Name != "ID" {
		t.Fatalf("expected PK to be ID, got %s", meta.PKField.Name)
	}
}

func TestParseWithPointerType(t *testing.T) {
	type Product struct {
		ID   uint
		Name string
	}

	ClearRegistry()
	meta, err := Parse(&Product{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should dereference to Product, not *Product
	if meta.Type.Kind() == reflect.Ptr {
		t.Fatal("expected dereferenced type")
	}
}

func TestCaching(t *testing.T) {
	type User struct {
		ID   uint
		Name string
	}

	ClearRegistry()

	meta1, err := Parse(&User{})
	if err != nil {
		t.Fatal(err)
	}

	meta2, err := Parse(&User{})
	if err != nil {
		t.Fatal(err)
	}

	// Should be the same instance
	if meta1 != meta2 {
		t.Fatal("expected cached metadata to be the same instance")
	}
}

func TestParseTableName(t *testing.T) {
	type Article struct {
		ID    uint
		Title string
	}

	ClearRegistry()
	meta, err := Parse(&Article{})

	if err != nil {
		t.Fatal(err)
	}

	// Default pluralization
	if meta.TableName == "" {
		t.Fatal("expected table name to be set")
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserName", "user_name"},
		{"ID", "id"},
		{"Product", "product"},
		{"HTTPServer", "http_server"},
	}

	for _, tt := range tests {
		result := toSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("toSnakeCase(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestParseFieldTags(t *testing.T) {
	type Product struct {
		ID       uint      `go-blar:"pk"`
		Name     string
		Secret   string    `go-blar:"hidden"`
		ReadOnly string    `go-blar:"readonly"`
		Items    []string  `go-blar:"list"`
	}

	ClearRegistry()
	meta, err := Parse(&Product{})

	if err != nil {
		t.Fatal(err)
	}

	// Find fields by name
	idField := meta.GetFieldByName("ID")
	if idField == nil || !idField.IsPK {
		t.Fatal("expected ID to be marked as PK")
	}

	secretField := meta.GetFieldByName("Secret")
	if secretField == nil || !secretField.Hidden {
		t.Fatal("expected Secret to be marked as hidden")
	}

	readOnlyField := meta.GetFieldByName("ReadOnly")
	if readOnlyField == nil || !readOnlyField.ReadOnly {
		t.Fatal("expected ReadOnly to be marked as readonly")
	}

	listField := meta.GetFieldByName("Items")
	if listField == nil || !listField.List {
		t.Fatal("expected Items to be marked as list")
	}
}

func TestGetFieldByName(t *testing.T) {
	type Entity struct {
		ID   uint
		Name string
		Age  int
	}

	ClearRegistry()
	meta, err := Parse(&Entity{})
	if err != nil {
		t.Fatal(err)
	}

	field := meta.GetFieldByName("Name")
	if field == nil {
		t.Fatal("expected to find Name field")
	}

	if field.Name != "Name" {
		t.Fatalf("expected field name Name, got %s", field.Name)
	}

	// Non-existent field
	field = meta.GetFieldByName("NonExistent")
	if field != nil {
		t.Fatal("expected nil for non-existent field")
	}
}
