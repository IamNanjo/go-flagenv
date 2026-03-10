package dotenv

import (
	"bufio"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/internal/convert"
	"github.com/IamNanjo/go-flagenv/internal/format"
	"github.com/IamNanjo/go-logging"
)

var quoteRunes = []byte{'"', '\''}

func Parse[T any](c *T, f *fields.Fields, path string) error {
	file, err := os.Open(path)
	if err != nil {
		logging.Debug("No .env file, skipping...\n")
		return nil
	}
	defer file.Close()

	// Scan .env one line at a time (default split function)
	scanner := bufio.NewScanner(file)

	envVars := map[string]string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore lines that are either empty, comments or have no =
		if len(line) == 0 || line[0] == '#' || !strings.Contains(line, "=") {
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

		envVars[k] = v
	}

	// Check if Config fields use any of the parsed variables
	for key, field := range f.Env {
		// Already set
		if field.Value.IsValid() && !field.Value.IsZero() {
			continue
		}

		val, exists := envVars[key]
		if !exists {
			continue
		}

		parsed, err := convert.AutoFromBytes(field.StructField.Type, field.Value, []byte(val))
		if err != nil {
			return format.Err("Failed to parse field %q with value %q %w", field.StructField.Name, val, err)
		}

		if convert.IsNormalPointer(field.StructField.Type) && reflect.TypeOf(parsed).Kind() != reflect.Pointer {
			if field.Value.IsNil() {
				field.Value.Set(reflect.New(reflect.TypeOf(parsed)))
			} else {
				field.Value.Elem().Set(reflect.ValueOf(parsed))
			}
		} else {
			field.Value.Set(reflect.ValueOf(parsed))
		}
	}

	return nil
}
