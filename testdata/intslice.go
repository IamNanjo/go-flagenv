package testdata

import (
	"strconv"
	"strings"

	"github.com/IamNanjo/go-flagenv/internal/format"
)

type IntSlice []int

func (is *IntSlice) UnmarshalText(input []byte) error {
	split := strings.Split(string(input), ",")
	result := make(IntSlice, 0, len(split))
	for _, s := range split {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}
		parsed, err := strconv.Atoi(trimmed)
		if err != nil {
			*is = result
			return format.Err("IntSlice parsing failed %w", err)
		}
		result = append(result, parsed)
	}
	*is = result
	return nil
}

func (is IntSlice) String() string {
	if is == nil || len(is) == 0 {
		return ""
	}

	var result strings.Builder
	for i, num := range is {
		result.WriteString(strconv.Itoa(num))
		if i < len(is)-1 {
			result.WriteString(", ")
		}
	}

	return result.String()
}
