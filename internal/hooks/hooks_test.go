package hooks

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

// MockEntityWithHooks implements all hook interfaces.
type MockEntityWithHooks struct {
	ID                 uint
	Name               string
	BeforeCreateCalled bool
	AfterCreateCalled  bool
	BeforeUpdateCalled bool
	AfterUpdateCalled  bool
	BeforeDeleteCalled bool
	AfterDeleteCalled  bool
	ShouldFail         bool
}

func (m *MockEntityWithHooks) BeforeCreate(ctx context.Context, tx *gorm.DB) error {
	m.BeforeCreateCalled = true
	if m.ShouldFail {
		return errors.New("before create failed")
	}
	return nil
}

func (m *MockEntityWithHooks) AfterCreate(ctx context.Context, tx *gorm.DB) error {
	m.AfterCreateCalled = true
	if m.ShouldFail {
		return errors.New("after create failed")
	}
	return nil
}

func (m *MockEntityWithHooks) BeforeUpdate(ctx context.Context, tx *gorm.DB) error {
	m.BeforeUpdateCalled = true
	if m.ShouldFail {
		return errors.New("before update failed")
	}
	return nil
}

func (m *MockEntityWithHooks) AfterUpdate(ctx context.Context, tx *gorm.DB) error {
	m.AfterUpdateCalled = true
	if m.ShouldFail {
		return errors.New("after update failed")
	}
	return nil
}

func (m *MockEntityWithHooks) BeforeDelete(ctx context.Context, tx *gorm.DB) error {
	m.BeforeDeleteCalled = true
	if m.ShouldFail {
		return errors.New("before delete failed")
	}
	return nil
}

func (m *MockEntityWithHooks) AfterDelete(ctx context.Context, tx *gorm.DB) error {
	m.AfterDeleteCalled = true
	if m.ShouldFail {
		return errors.New("after delete failed")
	}
	return nil
}

// SimpleEntity does not implement any hooks.
type SimpleEntity struct {
	ID   uint
	Name string
}

func TestCallBeforeCreate(t *testing.T) {
	ctx := context.Background()

	// Test with entity that implements BeforeCreate
	entity := &MockEntityWithHooks{}
	err := CallBeforeCreate(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.BeforeCreateCalled {
		t.Fatal("expected BeforeCreate to be called")
	}
}

func TestCallBeforeCreateWithError(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{ShouldFail: true}
	err := CallBeforeCreate(ctx, entity, nil)

	if err == nil {
		t.Fatal("expected error from BeforeCreate")
	}
}

func TestCallBeforeCreateNoImplementation(t *testing.T) {
	ctx := context.Background()

	entity := &SimpleEntity{}
	err := CallBeforeCreate(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error for entity without hook, got %v", err)
	}
}

func TestCallAfterCreate(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{}
	err := CallAfterCreate(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.AfterCreateCalled {
		t.Fatal("expected AfterCreate to be called")
	}
}

func TestCallBeforeUpdate(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{}
	err := CallBeforeUpdate(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.BeforeUpdateCalled {
		t.Fatal("expected BeforeUpdate to be called")
	}
}

func TestCallAfterUpdate(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{}
	err := CallAfterUpdate(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.AfterUpdateCalled {
		t.Fatal("expected AfterUpdate to be called")
	}
}

func TestCallBeforeDelete(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{}
	err := CallBeforeDelete(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.BeforeDeleteCalled {
		t.Fatal("expected BeforeDelete to be called")
	}
}

func TestCallAfterDelete(t *testing.T) {
	ctx := context.Background()

	entity := &MockEntityWithHooks{}
	err := CallAfterDelete(ctx, entity, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !entity.AfterDeleteCalled {
		t.Fatal("expected AfterDelete to be called")
	}
}

func TestMultipleHooks(t *testing.T) {
	ctx := context.Background()
	entity := &MockEntityWithHooks{}

	// Call all hooks in sequence
	if err := CallBeforeCreate(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}
	if err := CallAfterCreate(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}
	if err := CallBeforeUpdate(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}
	if err := CallAfterUpdate(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}
	if err := CallBeforeDelete(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}
	if err := CallAfterDelete(ctx, entity, nil); err != nil {
		t.Fatal(err)
	}

	// All hooks should have been called
	if !entity.BeforeCreateCalled || !entity.AfterCreateCalled ||
		!entity.BeforeUpdateCalled || !entity.AfterUpdateCalled ||
		!entity.BeforeDeleteCalled || !entity.AfterDeleteCalled {
		t.Fatal("expected all hooks to be called")
	}
}
