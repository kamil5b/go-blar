package meta

import "reflect"

// FieldMeta holds metadata about a single field in an entity.
type FieldMeta struct {
	Name     string
	Type     reflect.Type
	Index    []int // NestedIndex for embedded structs
	IsPK     bool
	FK       *ForeignKey
	Nested   bool
	M2M      *ManyToMany
	List     bool
	Hidden   bool
	ReadOnly bool
}

// ForeignKey holds metadata for a foreign key relationship.
type ForeignKey struct {
	TableName string
	FieldName string
}

// ManyToMany holds metadata for a many-to-many relationship.
type ManyToMany struct {
	TableName string
	FieldName string
}

// NestedMeta holds metadata about a nested/embedded struct.
type NestedMeta struct {
	Name  string
	Type  reflect.Type
	Index []int
}

// AggregateMeta holds metadata about an aggregate field.
type AggregateMeta struct {
	Name  string
	Type  string // "count", "sum", etc.
	Field string // nested field path
}

// EntityMeta holds all metadata about an entity type.
type EntityMeta struct {
	Type       reflect.Type
	Name       string
	TableName  string
	PKField    *FieldMeta
	Fields     []*FieldMeta
	Nested     []*NestedMeta
	Aggregates []*AggregateMeta
}

// GetFieldByName returns a field by its name.
func (em *EntityMeta) GetFieldByName(name string) *FieldMeta {
	for _, f := range em.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// GetAggregateByName returns an aggregate by its name.
func (em *EntityMeta) GetAggregateByName(name string) *AggregateMeta {
	for _, a := range em.Aggregates {
		if a.Name == name {
			return a
		}
	}
	return nil
}
