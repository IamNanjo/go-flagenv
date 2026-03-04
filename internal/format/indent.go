package format

import "strings"

func IndentAllLines(input string, indentWidth int) string {
	indent := strings.Repeat(" ", indentWidth)

	var result strings.Builder
	for l := range strings.Lines(input) {
		result.WriteString(indent)
		result.WriteString(l)
	}

	return result.String()
}
