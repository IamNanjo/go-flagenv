package required

import (
	"fmt"

	"github.com/IamNanjo/go-flagenv/fields"

	"github.com/IamNanjo/go-logging/pkg/format"
)

// Will panic if required fields are not provided
func Check[T any](c *T, f *fields.Fields) error {
	for _, field := range f.Required {
		if field.Value.IsZero() {
			flagTag, hasFlag := field.StructField.Tag.Lookup("flag")
			envTag, hasEnv := field.StructField.Tag.Lookup("env")

			optionName := ""

			if hasFlag && hasEnv {
				optionName = fmt.Sprintf("-%s (%s)", flagTag, envTag)
			} else {
				switch {
				case hasFlag:
					optionName = flagTag
				case hasEnv:
					optionName = envTag
				default:
					return format.Err("Required field %q does not have flag or env struct tags")
				}
			}

			return format.Err("Required option %s not set\n", optionName)
		}
	}
	return nil
}
