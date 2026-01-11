package hooks

import (
	"context"

	"github.com/kamil5b/go-blar/goblar"
	"gorm.io/gorm"
)

// CallBeforeCreate calls the BeforeCreate hook on the entity if it implements it.
func CallBeforeCreate(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.BeforeCreate); ok {
		return h.BeforeCreate(ctx, tx)
	}
	return nil
}

// CallAfterCreate calls the AfterCreate hook on the entity if it implements it.
func CallAfterCreate(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.AfterCreate); ok {
		return h.AfterCreate(ctx, tx)
	}
	return nil
}

// CallBeforeUpdate calls the BeforeUpdate hook on the entity if it implements it.
func CallBeforeUpdate(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.BeforeUpdate); ok {
		return h.BeforeUpdate(ctx, tx)
	}
	return nil
}

// CallAfterUpdate calls the AfterUpdate hook on the entity if it implements it.
func CallAfterUpdate(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.AfterUpdate); ok {
		return h.AfterUpdate(ctx, tx)
	}
	return nil
}

// CallBeforeDelete calls the BeforeDelete hook on the entity if it implements it.
func CallBeforeDelete(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.BeforeDelete); ok {
		return h.BeforeDelete(ctx, tx)
	}
	return nil
}

// CallAfterDelete calls the AfterDelete hook on the entity if it implements it.
func CallAfterDelete(ctx context.Context, entity any, tx *gorm.DB) error {
	if h, ok := entity.(goblar.AfterDelete); ok {
		return h.AfterDelete(ctx, tx)
	}
	return nil
}
