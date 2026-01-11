package aggregate

import (
	"github.com/kamil5b/go-blar/internal/meta"
	"gorm.io/gorm"
)

// ComputeAggregates computes aggregate values (count, sum, etc.) for entities.
// This is called internally when loading related data.
func ComputeAggregates(db *gorm.DB, entity any, aggregates []*meta.AggregateMeta) error {
	// TODO: Implement aggregate computation
	// This will handle count:, sum:, and other aggregate directives
	return nil
}
