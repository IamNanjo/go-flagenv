package convert

import (
	"reflect"
	"strings"

	"github.com/IamNanjo/go-logging/pkg/format"
)

// Automatically convert to supported type
func AutoFromBytes(t reflect.Type, input []byte) (any, error) {
	if parser := CustomParserFromType(t); parser != nil {
		err := parser.UnmarshalText(input)
		if err != nil {
			return nil, format.Err("CustomParser.UnmarshalText failed %w", err)
		}
		return parser, nil
	}

	isPointer := t.Kind() == reflect.Pointer
	actualType := t
	if isPointer {
		actualType = t.Elem()
	}

	if actualType.Kind() == reflect.Slice {
		sliceItemType := actualType.Elem()
		result := reflect.MakeSlice(actualType, 0, 8)

		for s := range strings.SplitSeq(string(input), ",") {
			s = strings.TrimSpace(s)
			converted, err := AutoFromBytes(sliceItemType, []byte(s))
			if err != nil {
				return nil, format.Err("Failed to convert slice element %w", err)
			}
			result = reflect.Append(result, reflect.ValueOf(converted))
		}

		resultVal := reflect.New(result.Type())
		resultVal.Elem().Set(result)
		if isPointer {
			return resultVal.Interface(), nil
		} else {
			return resultVal.Elem().Interface(), nil
		}
	}

	c, compatible := FromBytes[actualType]
	if !compatible {
		return nil, format.Err("Unsupported type %q", t)
	}

	parsed, err := c(input)
	if err != nil {
		return nil, err
	}
	parsedValue := reflect.ValueOf(parsed)

	resultVal := reflect.New(parsedValue.Type())
	resultVal.Elem().Set(parsedValue)
	if isPointer {
		return resultVal.Interface(), nil
	} else {
		return resultVal.Elem().Interface(), nil
	}
}
