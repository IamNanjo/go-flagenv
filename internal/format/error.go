package format

import (
	"fmt"
	"strings"
)

const (
	minLength     = 10 // Minimum line length before wrapping.
	minWrapLength = 70 // input must be this long before text wraps at all.
	wrapCutoff    = 60 // How many characters until text wraps on each line.
	prefixLength  = 6  // Constant prefix length to keep each line aligned. Includes newline character

	errPrefix     = "\n  -> " // Prefix for the error.
	errLinePrefix = "\n     " // Prefix for each line
)

// Consistent error format with wrapping.
func Err(input string, args ...any) error {
	inputLen := len(input)
	if len(input) < minWrapLength {
		return fmt.Errorf(errPrefix+input, args...)
	}

	words := strings.Fields(input)
	if len(words) == 0 {
		return fmt.Errorf(errPrefix + "Error with no message")
	}
	var wrappedInput strings.Builder
	wrappedInput.Grow(inputLen)
	wrappedInput.WriteString(errPrefix)
	wrappedInput.WriteString(words[0])
	lineLength := prefixLength + len(words[0])

	for _, w := range words[1:] {
		wLen := len(w)
		// Very long word -> put it on its own line
		if lineLength > minLength && prefixLength+wLen >= wrapCutoff {
			wrappedInput.WriteString(errLinePrefix)
			wrappedInput.WriteString(w)
			continue
		}

		if lineLength+wLen > wrapCutoff { // Line is long enough to wrap
			wrappedInput.WriteString(errLinePrefix)
			lineLength = prefixLength
		} else { // No need to wrap yet
			wrappedInput.WriteByte(' ')
			lineLength++
		}

		wrappedInput.WriteString(w)
		lineLength += wLen
	}

	return fmt.Errorf(wrappedInput.String(), args...)
}
