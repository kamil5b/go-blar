package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kamil5b/go-blar/internal/hooks"
	"github.com/kamil5b/go-blar/internal/meta"
	"gorm.io/gorm"
)

// Handlers provides HTTP handlers for entity CRUD operations.
// It is used internally and should not be exposed in the public API.
type Handlers struct {
	db *gorm.DB
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(db *gorm.DB) *Handlers {
	return &Handlers{db: db}
}

// CreateHandler returns an HTTP handler for creating a new entity.
func (h *Handlers) CreateHandler(entityMeta *meta.EntityMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Decode JSON body
		entity := makeEntityInstance(entityMeta)
		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Call BeforeCreate hook
		if err := hooks.CallBeforeCreate(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save to database
		if err := h.db.WithContext(ctx).Create(entity).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call AfterCreate hook
		if err := hooks.CallAfterCreate(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entity)
	}
}

// ListHandler returns an HTTP handler for listing all entities.
func (h *Handlers) ListHandler(entityMeta *meta.EntityMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Create a slice of the entity type
		entities := makeEntitySlice(entityMeta)

		// Query database
		if err := h.db.WithContext(ctx).Find(entities).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entities)
	}
}

// GetHandler returns an HTTP handler for retrieving a single entity by ID.
func (h *Handlers) GetHandler(entityMeta *meta.EntityMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Query database
		entity := makeEntityInstance(entityMeta)
		if err := h.db.WithContext(ctx).First(entity, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// UpdateHandler returns an HTTP handler for updating an entity.
func (h *Handlers) UpdateHandler(entityMeta *meta.EntityMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Decode JSON body
		entity := makeEntityInstance(entityMeta)
		if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Call BeforeUpdate hook
		if err := hooks.CallBeforeUpdate(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update in database
		if err := h.db.WithContext(ctx).Model(entity).Where("id = ?", id).Updates(entity).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call AfterUpdate hook
		if err := hooks.CallAfterUpdate(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entity)
	}
}

// DeleteHandler returns an HTTP handler for deleting an entity.
func (h *Handlers) DeleteHandler(entityMeta *meta.EntityMeta) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract ID from URL
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Fetch entity first (for hooks)
		entity := makeEntityInstance(entityMeta)
		if err := h.db.WithContext(ctx).First(entity, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Call BeforeDelete hook
		if err := hooks.CallBeforeDelete(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete from database
		if err := h.db.WithContext(ctx).Delete(entity, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call AfterDelete hook
		if err := hooks.CallAfterDelete(ctx, entity, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// makeEntityInstance creates a new instance of the entity type.
func makeEntityInstance(entityMeta *meta.EntityMeta) any {
	return reflect.New(entityMeta.Type).Interface()
}

// makeEntitySlice creates a new slice of the entity type.
func makeEntitySlice(entityMeta *meta.EntityMeta) any {
	sliceType := reflect.SliceOf(entityMeta.Type)
	return reflect.New(sliceType).Interface()
}
