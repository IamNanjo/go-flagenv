package dotenv

import (
	"bufio"
	"os"
	"slices"
	"strings"

	"github.com/IamNanjo/go-flagenv/pkg/fields"

	"github.com/IamNanjo/go-logging"
	"github.com/IamNanjo/go-logging/pkg/format"
)

var quoteRunes = []byte{'"', '\''}

// Parse .env file and set all
func Parse[T any](c *T, f *fields.Fields, path string) error {
	file, err := os.Open(path)
	if err != nil {
		logging.Debug("No .env file, skipping...\n")
		return nil
	}
	defer file.Close()

	// Scan .env one line at a time (default split function)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := []rune(strings.TrimSpace(scanner.Text()))

		// Ignore lines that are either empty, comments or have no =
		if len(line) == 0 || line[0] == '#' || !slices.Contains(line, '=') {
			continue
		}

		var key strings.Builder
		var value strings.Builder

		readingKey := true
	charLoop:
		for _, char := range line {
			if readingKey {
				switch char {
				case '#':
					return format.Err("Encountered comment while parsing %q", line)
				case ' ':
					key.Reset()
					continue
				case '=':
					readingKey = false
				default:
					key.WriteRune(char)
				}
			} else {
				switch char {
				case '#':
					break charLoop
				default:
					value.WriteRune(char)
				}
			}
		}

		// Remove quotes from value
		k := strings.TrimSpace(key.String())
		v := strings.TrimSpace(value.String())
		vLen := len(v)
		if vLen >= 2 && slices.Contains(quoteRunes, v[0]) && v[0] == v[vLen-1] {
			v = v[1 : vLen-1]
		}

		os.Setenv(k, v)
	}
	return nil
}
