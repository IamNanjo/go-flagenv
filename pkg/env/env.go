package env

import (
	"os"
	"reflect"

	"github.com/IamNanjo/go-flagenv/pkg/convert"
	"github.com/IamNanjo/go-flagenv/pkg/fields"

	"github.com/IamNanjo/go-logging/pkg/format"
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

		parsed, err := convert.AutoFromBytes(field.StructField.Type, []byte(val))
		if err != nil {
			return format.Err("Failed to parse field %q with value %q %w", field.StructField.Name, val, err)
		}

		field.Value.Set(reflect.ValueOf(parsed))
	}

	return nil
}
