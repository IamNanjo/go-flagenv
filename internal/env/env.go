package env

import (
	"os"
	"reflect"

	"github.com/IamNanjo/go-flagenv/convert"
	"github.com/IamNanjo/go-flagenv/internal/fields"
	"github.com/IamNanjo/go-flagenv/internal/format"
)

func Parse[T any](c *T, f *fields.Fields) error {
	for key, field := range f.Env {
		// Already set
		if field.Value.IsValid() && !field.Value.IsZero() {
			continue
		}

		val, exists := os.LookupEnv(key)
		if !exists {
			continue
		}

		parsed, err := convert.AutoFromString(field.StructField.Type, field.Value, val)
		if err != nil {
			return format.Err("Failed to parse field %q with value %q %w", field.StructField.Name, val, err)
		}

		field.Value.Set(reflect.ValueOf(parsed))
	}

	return nil
}
