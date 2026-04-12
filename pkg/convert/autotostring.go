package convert

import (
	"reflect"
	"strings"
)

// Automatically convert to supported type
func AutoToString(input any) string {
	v := reflect.ValueOf(input)
	if parser := CustomParserFromType(v.Type()); parser != nil {
		return input.(CustomParser).String()
	}

	isPointer := v.Kind() == reflect.Pointer
	if isPointer {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		vLen := v.Len()
		if vLen == 0 {
			return ""
		}

		var result strings.Builder
		for i := range vLen {
			item := v.Index(i)
			if item.Kind() == reflect.Pointer {
				item = item.Elem()
			}

			result.WriteString(AutoToString(item.Interface()))

			if i < vLen-1 {
				result.WriteByte(',')
				result.WriteByte(' ')
			}
		}
		return result.String()
	}

	return ToString[v.Type()](v.Interface())
}
