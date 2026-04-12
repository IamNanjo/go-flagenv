package convert

import (
	"reflect"

	"github.com/IamNanjo/go-logging/pkg/format"
)

var CustomParserType = reflect.TypeFor[CustomParser]()

// Returns error if unsupported
func FieldIsSupported(t reflect.Type) error {
	isPointer := t.Kind() == reflect.Pointer
	if t.Implements(CustomParserType) {
		if isPointer {
			return nil
		} else {
			return format.Err("Fields that implement CustomParser must be pointers")
		}
	}

	actualType := t
	if isPointer {
		actualType = t.Elem()
	}

	if actualType.Kind() == reflect.Slice {
		actualType = actualType.Elem()

		if actualType.Kind() == reflect.Pointer {
			actualType = actualType.Elem()
		}
	}

	for st := range FromBytes {
		if st == actualType {
			return nil
		}
	}

	return format.Err("Field does not implement CustomParser and is not one of the supported native types (%q)", t)
}

// Returns true if type is a pointer to a non CustomParser interface.
func IsNormalPointer(t reflect.Type) bool {
	pointer := t.Kind() == reflect.Pointer
	return pointer && CustomParserFromType(t) == nil
}

// Make a new CustomParser of type t.
func CustomParserFromType(t reflect.Type) CustomParser {
	if t.Implements(CustomParserType) {
		parser, ok := reflect.New(t.Elem()).Interface().(CustomParser)
		if !ok {
			return nil
		}
		return parser
	}
	return nil
}
