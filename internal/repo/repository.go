package repo

import (
	"context"

	"github.com/kamil5b/go-blar/internal/meta"
	"gorm.io/gorm"
)

// Repository is a generic CRUD repository for any entity type.
type Repository[T any] struct {
	db   *gorm.DB
	meta *meta.EntityMeta
}

// New creates a new repository for the given entity type.
func New[T any](db *gorm.DB, meta *meta.EntityMeta) *Repository[T] {
	return &Repository[T]{
		db:   db,
		meta: meta,
	}
}

// Create creates a new entity in the database.
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}

	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves an entity by its primary key.
func (r *Repository[T]) GetByID(ctx context.Context, id any) (*T, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// GetAll retrieves all entities.
func (r *Repository[T]) GetAll(ctx context.Context) ([]T, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// Update updates an entity in the database.
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}

	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete deletes an entity from the database.
func (r *Repository[T]) Delete(ctx context.Context, id any) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}

	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// Count returns the total number of entities.
func (r *Repository[T]) Count(ctx context.Context) (int64, error) {
	if r.db == nil {
		return 0, gorm.ErrInvalidDB
	}

	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error
	return count, err
}
