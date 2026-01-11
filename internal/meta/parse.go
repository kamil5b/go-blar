package meta

import (
	"reflect"
	"strings"
)

// registry is a global cache of parsed entity metadata.
var registry = make(map[string]*EntityMeta)

// Parse parses a struct and returns its EntityMeta.
// The struct must be a GORM entity with struct tags.
func Parse(entity any) (*EntityMeta, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check cache
	key := t.String()
	if meta, ok := registry[key]; ok {
		return meta, nil
	}

	// Parse new metadata
	meta := &EntityMeta{
		Type:       t,
		Name:       t.Name(),
		TableName:  toSnakeCase(t.Name()) + "s", // Default pluralization
		Fields:     make([]*FieldMeta, 0),
		Nested:     make([]*NestedMeta, 0),
		Aggregates: make([]*AggregateMeta, 0),
	}

	// Extract table name from gorm tag if present
	if gormTag, ok := t.FieldByName("gorm"); ok {
		if tableName := parseGormTag(gormTag.Tag.Get("gorm")); tableName != "" {
			meta.TableName = tableName
		}
	}

	// Parse fields
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// Skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		fm := parseField(sf)
		if fm != nil {
			meta.Fields = append(meta.Fields, fm)

			// Track primary key
			if fm.IsPK {
				meta.PKField = fm
			}
		}
	}

	// Cache it
	registry[key] = meta

	return meta, nil
}

// parseField extracts metadata from a struct field.
func parseField(sf reflect.StructField) *FieldMeta {
	blarTag := sf.Tag.Get("go-blar")
	gormTag := sf.Tag.Get("gorm")

	if blarTag == "" && gormTag == "" {
		// Include fields without tags by default
	}

	fm := &FieldMeta{
		Name:  sf.Name,
		Type:  sf.Type,
		Index: sf.Index,
	}

	// Parse go-blar tags
	if blarTag != "" {
		parts := strings.Split(blarTag, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)

			switch {
			case part == "pk":
				fm.IsPK = true
			case part == "nested":
				fm.Nested = true
			case part == "list":
				fm.List = true
			case part == "hidden":
				fm.Hidden = true
			case part == "readonly":
				fm.ReadOnly = true
			case strings.HasPrefix(part, "fk:"):
				fkTable := strings.TrimPrefix(part, "fk:")
				fm.FK = &ForeignKey{TableName: fkTable}
			case strings.HasPrefix(part, "m2m:"):
				m2mTable := strings.TrimPrefix(part, "m2m:")
				fm.M2M = &ManyToMany{TableName: m2mTable}
			case strings.HasPrefix(part, "count:"):
				fieldPath := strings.TrimPrefix(part, "count:")
				meta := &AggregateMeta{
					Name:  fm.Name,
					Type:  "count",
					Field: fieldPath,
				}
				// This will be added to EntityMeta later
				_ = meta
			case strings.HasPrefix(part, "sum:"):
				fieldPath := strings.TrimPrefix(part, "sum:")
				meta := &AggregateMeta{
					Name:  fm.Name,
					Type:  "sum",
					Field: fieldPath,
				}
				// This will be added to EntityMeta later
				_ = meta
			}
		}
	}

	// Parse gorm tags for primary key detection
	if strings.Contains(gormTag, "primaryKey") {
		fm.IsPK = true
	}

	return fm
}

// parseGormTag extracts the table name from a gorm tag.
func parseGormTag(tag string) string {
	if tag == "" {
		return ""
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "table:") {
			return strings.TrimPrefix(part, "table:")
		}
	}

	return ""
}

// toSnakeCase converts CamelCase to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if i == 0 {
			result.WriteRune(r)
			continue
		}

		// Avoid adding extra underscores if current or previous char is '_'
		if r == '_' {
			result.WriteRune(r)
			continue
		}

		if r >= 'A' && r <= 'Z' {
			if runes[i-1] != '_' && runes[i-1] >= 'a' && runes[i-1] <= 'z' {
				result.WriteRune('_')
			} else if runes[i-1] != '_' && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ClearRegistry clears the metadata cache (useful for testing).
func ClearRegistry() {
	registry = make(map[string]*EntityMeta)
}
