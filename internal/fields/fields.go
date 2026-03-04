package fields

import (
	"reflect"

	"github.com/IamNanjo/go-flagenv/convert"
	"github.com/IamNanjo/go-flagenv/internal/format"
	"github.com/IamNanjo/go-logging"
)

type Fields struct {
	Flags    map[string]*Field // CLI flags
	Env      map[string]*Field // Environment variables
	Required []*Field          // Required fields
	Defaults []*Field          // All fields with default values
}
type Field struct {
	StructField reflect.StructField // Field type
	Value       reflect.Value       // Value
	Default     any                 // Default value correct type for the field
	Description *string             // Description
}

func Parse[T any](c *T) (*Fields, error) {
	f := &Fields{
		Flags:    map[string]*Field{},
		Env:      map[string]*Field{},
		Required: []*Field{},
		Defaults: []*Field{},
	}

	config := reflect.ValueOf(c).Elem()
	configType := config.Type()

	if configKind := config.Kind(); configKind != reflect.Struct {
		return f, format.Err("Expected config to be a struct. Received %q instead", configKind)
	}

	for i := range configType.NumField() {
		structField := configType.Field(i)
		fieldValue := config.Field(i)

		if !fieldValue.CanSet() {
			logging.Debug("Skipping non-settable field %q\n", structField.Name)
			continue
		}

		if err := convert.FieldIsSupported(structField.Type); err != nil {
			return f, format.Err("Field %q type is unsupported %w", structField.Name, err)
		}

		field := &Field{StructField: structField, Value: fieldValue}

		flagName, flagTagSet := structField.Tag.Lookup("flag")
		if flagTagSet {
			f.Flags[flagName] = field
		}

		envTag, envTagSet := structField.Tag.Lookup("env")
		if envTagSet {
			f.Env[envTag] = field
		}

		_, required := structField.Tag.Lookup("required")
		if required {
			f.Required = append(f.Required, field)
		}

		defaultValueString, hasDefault := structField.Tag.Lookup("default")
		if hasDefault {
			defaultValue, err := convert.AutoFromString(structField.Type, fieldValue, defaultValueString)
			if err != nil {
				return f, format.Err("Field %q default value parsing failed %w", structField.Name, err)
			}
			field.Default = defaultValue
		} else {
			field.Default = reflect.Zero(field.StructField.Type).Interface()
		}

		f.Defaults = append(f.Defaults, field)

		description, hasDescription := structField.Tag.Lookup("desc")
		if hasDescription {
			field.Description = &description
		}
	}

	return f, nil
}
