package defaults

import (
	"reflect"

	"github.com/IamNanjo/go-flagenv/fields"
)

// Set defaults if not defined by user
func Set[T any](c *T, f *fields.Fields) error {
	for _, field := range f.Defaults {
		if field.Value.IsZero() {
			field.Value.Set(reflect.ValueOf(field.Default))
		}
	}
	return nil
}
