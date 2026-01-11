package goblar

import (
	"context"

	"gorm.io/gorm"
)

// BeforeCreate is called before an entity is created.
type BeforeCreate interface {
	BeforeCreate(ctx context.Context, tx *gorm.DB) error
}

// AfterCreate is called after an entity is created.
type AfterCreate interface {
	AfterCreate(ctx context.Context, tx *gorm.DB) error
}

// BeforeUpdate is called before an entity is updated.
type BeforeUpdate interface {
	BeforeUpdate(ctx context.Context, tx *gorm.DB) error
}

// AfterUpdate is called after an entity is updated.
type AfterUpdate interface {
	AfterUpdate(ctx context.Context, tx *gorm.DB) error
}

// BeforeDelete is called before an entity is deleted.
type BeforeDelete interface {
	BeforeDelete(ctx context.Context, tx *gorm.DB) error
}

// AfterDelete is called after an entity is deleted.
type AfterDelete interface {
	AfterDelete(ctx context.Context, tx *gorm.DB) error
}
