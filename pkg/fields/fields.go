// This package handles parsing supported fields in structs.
package fields

import (
	"reflect"

	"github.com/IamNanjo/go-flagenv/pkg/convert"

	"github.com/IamNanjo/go-logging"
	"github.com/IamNanjo/go-logging/pkg/format"
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
	Default     any                 // Default value (correct type for the field)
	Description *string             // Description
}

func Parse[T any](c *T) (*Fields, error) {
	f := &Fields{
		Flags:    map[string]*Field{},
		Env:      map[string]*Field{},
		Required: []*Field{},
		Defaults: []*Field{},
	}

	return parse(f, c, "", "")
}

func parse[T any](f *Fields, c T, flagPrefix string, envPrefix string) (*Fields, error) {
	config := reflect.ValueOf(c)
	if config.Kind() == reflect.Pointer {
		config = config.Elem()
	}
	configType := config.Type()

	if configKind := config.Kind(); configKind != reflect.Struct {
		return f, format.Err("Expected config to be a struct. Received %q instead", configKind)
	}

	for i := range configType.NumField() {
		structField := configType.Field(i)
		fieldValue := config.Field(i)

		// Recursive parsing for nested structs
		if !structField.Type.Implements(convert.CustomParserType) {
			if structField.Type.Kind() == reflect.Struct {
				fPrefix := flagPrefix + structField.Tag.Get("flag")
				ePrefix := envPrefix + structField.Tag.Get("env")
				fVal := fieldValue.Addr().Interface()
				if _, err := parse(f, fVal, fPrefix, ePrefix); err != nil {
					return f, format.Err("Nested field parsing failed %w", err)
				}
				continue
			}

			if structField.Type.Kind() == reflect.Pointer && structField.Type.Elem().Kind() == reflect.Struct {
				fPrefix := flagPrefix + structField.Tag.Get("flag")
				ePrefix := envPrefix + structField.Tag.Get("env")
				fVal := fieldValue.Interface()
				if _, err := parse(f, fVal, fPrefix, ePrefix); err != nil {
					return f, format.Err("Nested field parsing failed %w", err)
				}
				continue
			}
		}

		if !fieldValue.CanSet() {
			logging.Default.Debug("Skipping non-settable field %q\n", structField.Name)
			continue
		}

		if err := convert.FieldIsSupported(structField.Type); err != nil {
			logging.Default.Debug("Skipping unsupported field %q\n", structField.Name)
			continue
		}

		field := &Field{StructField: structField, Value: fieldValue}

		flagName, flagTagSet := structField.Tag.Lookup("flag")
		if flagTagSet {
			flagName = flagPrefix + flagName
			f.Flags[flagName] = field
		}

		envTag, envTagSet := structField.Tag.Lookup("env")
		if envTagSet {
			envTag = envPrefix + envTag
			f.Env[envTag] = field
		}

		if !flagTagSet && !envTagSet {
			logging.Default.Debug("Skipping field with no flag or env tag: %q\n", structField.Name)
			continue
		}

		_, required := structField.Tag.Lookup("required")
		if required {
			f.Required = append(f.Required, field)
		}

		defaultValueString, hasDefault := structField.Tag.Lookup("default")
		if hasDefault {
			defaultValue, err := convert.AutoFromBytes(structField.Type, []byte(defaultValueString))
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
