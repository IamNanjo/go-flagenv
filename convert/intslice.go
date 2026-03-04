package convert

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/IamNanjo/go-flagenv/internal/format"
)

var IntSliceType = reflect.TypeFor[[]int]()

// String to integer slice conversion.
// Expects comma separated list as string and trims space around each list item.
func ToIntSlice(input string) ([]int, error) {
	result := make([]int, 0)
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 || !strings.Contains(trimmed, ",") {
		return result, nil
	}
	split := strings.SplitSeq(trimmed, ",")
	for val := range split {
		parsed, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return result, format.Err("Value %v cannot be parsed to int", val)
		}
		result = append(result, parsed)
	}
	return result, nil
}

// Opposite of ToIntSlice. Joins items with ", ".
func FromIntSlice(input []int) string {
	if len(input) == 0 {
		return ""
	}

	var result strings.Builder

	last := len(input) - 1
	for i, num := range input {
		result.WriteString(strconv.Itoa(num))
		if i < last {
			result.WriteString(", ")
		}
	}

	return result.String()
}
