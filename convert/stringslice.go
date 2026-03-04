package convert

import (
	"reflect"
	"strings"
)

var StringSliceType = reflect.TypeFor[[]string]()

// String to string slice conversion.
// Expects comma separated list as string and trims space around each list item.
func ToStringSlice(input string) []string {
	result := make([]string, 0)
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 || !strings.Contains(trimmed, ",") {
		return result
	}
	split := strings.SplitSeq(trimmed, ",")
	for val := range split {
		result = append(result, strings.TrimSpace(val))
	}
	return result
}

// Opposite of ToStringSlice. Joins items with ", ".
func FromStringSlice(input []string) string {
	if len(input) == 0 {
		return ""
	}
	return strings.Join(input, ", ")
}
