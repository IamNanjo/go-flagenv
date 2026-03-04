package convert

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IamNanjo/go-flagenv/internal/format"
)

type CustomParser interface {
	// Parse string into CustomParser compatible type.
	FromString(input string) error
	// Stringify CustomParser. Should return same thing that FromString took as input.
	String() string
}

var CustomParserType = reflect.TypeFor[CustomParser]()

// Map of supported type to FromString function
var FromString = map[reflect.Type]func(input string) (any, error){
	reflect.TypeFor[bool](): func(input string) (any, error) {
		return strings.ToLower(input) == "true", nil
	},
	reflect.TypeFor[int](): func(input string) (any, error) {
		return strconv.Atoi(input)
	},
	reflect.TypeFor[int64](): func(input string) (any, error) {
		return strconv.ParseInt(input, 10, 64)
	},
	reflect.TypeFor[uint](): func(input string) (any, error) {
		u64, err := strconv.ParseUint(input, 10, 0)
		return uint(u64), err
	},
	reflect.TypeFor[uint64](): func(input string) (any, error) {
		return strconv.ParseUint(input, 10, 64)
	},
	reflect.TypeFor[float64](): func(input string) (any, error) {
		return strconv.ParseFloat(input, 64)
	},
	reflect.TypeFor[string](): func(input string) (any, error) {
		return input, nil
	},
	reflect.TypeFor[[]int](): func(input string) (any, error) {
		return ToStringSlice(input), nil
	},
	reflect.TypeFor[[]string](): func(input string) (any, error) {
		return ToStringSlice(input), nil
	},
	reflect.TypeFor[time.Duration](): func(input string) (any, error) {
		return time.ParseDuration(input)
	},
}

// Returns error if unsupported
func FieldIsSupported(t reflect.Type) error {
	if t.Implements(reflect.TypeFor[CustomParser]()) {
		if t.Kind() == reflect.Pointer {
			return nil
		} else {
			return format.Err("Fields that implement CustomParser must be pointers")
		}
	}

	isPointer := t.Kind() == reflect.Pointer
	actualType := t
	if isPointer {
		actualType = t.Elem()
	}

	for st := range FromString {
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

// Automatically convert to supported type
func AutoFromString(t reflect.Type, v reflect.Value, input string) (any, error) {
	if parser := CustomParserFromType(t); parser != nil {
		err := parser.FromString(input)
		if err != nil {
			return parser, format.Err("CustomParser.FromString failed %w", err)
		}
		return parser, nil
	}

	isPointer := t.Kind() == reflect.Pointer
	actualType := t
	if isPointer {
		actualType = t.Elem()
	}

	for st, c := range FromString {
		if actualType == st {
			converted, err := c(input)
			return converted, err
		}
	}

	return nil, format.Err("Incompatible type %q", t)
}
