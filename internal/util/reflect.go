package util

import (
	"reflect"
	"strings"
)

// DeepFieldByName finds a field in a struct, including nested fields.
// It follows a dot-separated path like "User.Name".
func DeepFieldByName(v interface{}, fieldPath string) reflect.Value {
	parts := strings.Split(fieldPath, ".")

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, part := range parts {
		if val.Kind() != reflect.Struct {
			return reflect.Value{}
		}

		field := val.FieldByName(part)
		if !field.IsValid() {
			return reflect.Value{}
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		val = field
	}

	return val
}

// IsZeroValue checks if a value is its zero value.
func IsZeroValue(v reflect.Value) bool {
	return v.IsZero()
}
