package convert

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IamNanjo/go-flagenv/convert"
	"github.com/IamNanjo/go-flagenv/internal/format"
)

var CustomParserType = reflect.TypeFor[convert.CustomParser]()

// Map of supported type to FromString function
var FromString = map[reflect.Type]func(input []byte) (any, error){
	reflect.TypeFor[bool](): func(input []byte) (any, error) {
		return strings.ToLower(string(input)) == "true", nil
	},
	reflect.TypeFor[int](): func(input []byte) (any, error) {
		return strconv.Atoi(string(input))
	},
	reflect.TypeFor[int64](): func(input []byte) (any, error) {
		return strconv.ParseInt(string(input), 10, 64)
	},
	reflect.TypeFor[uint](): func(input []byte) (any, error) {
		u64, err := strconv.ParseUint(string(input), 10, 0)
		return uint(u64), err
	},
	reflect.TypeFor[uint64](): func(input []byte) (any, error) {
		return strconv.ParseUint(string(input), 10, 64)
	},
	reflect.TypeFor[float64](): func(input []byte) (any, error) {
		return strconv.ParseFloat(string(input), 64)
	},
	reflect.TypeFor[string](): func(input []byte) (any, error) {
		return string(input), nil
	},
	reflect.TypeFor[[]int](): func(input []byte) (any, error) {
		return ToStringSlice(string(input)), nil
	},
	reflect.TypeFor[[]string](): func(input []byte) (any, error) {
		return ToStringSlice(string(input)), nil
	},
	reflect.TypeFor[time.Duration](): func(input []byte) (any, error) {
		return time.ParseDuration(string(input))
	},
}

// Returns error if unsupported
func FieldIsSupported(t reflect.Type) error {
	if t.Implements(CustomParserType) {
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
func CustomParserFromType(t reflect.Type) convert.CustomParser {
	if t.Implements(CustomParserType) {
		parser, ok := reflect.New(t.Elem()).Interface().(convert.CustomParser)
		if !ok {
			return nil
		}
		return parser
	}
	return nil
}

// Automatically convert to supported type
func AutoFromBytes(t reflect.Type, v reflect.Value, input []byte) (any, error) {
	if parser := CustomParserFromType(t); parser != nil {
		err := parser.UnmarshalText(input)
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
