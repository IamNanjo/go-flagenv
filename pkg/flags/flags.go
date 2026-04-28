package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/IamNanjo/go-flagenv/pkg/convert"
	"github.com/IamNanjo/go-flagenv/pkg/fields"

	"github.com/IamNanjo/go-logging/pkg/ansi"
	"github.com/IamNanjo/go-logging/pkg/format"
)

const usageIndent = "    "
const usageSeparator = " · "

type helpError struct {
	Message strings.Builder
	Inner   error
}

func (e helpError) Error() string {
	return e.Message.String()
}
func (e helpError) Unwrap() error {
	return e.Inner
}

var HelpError = helpError{Inner: flag.ErrHelp}

func Parse(f *fields.Fields, args []string) error {
	// Functions to be called after flags are parsed (value has been set in the pointer)
	var postProcess = make(map[string]func() error, 0)

	flagSet := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	for flagName, field := range f.Flags {
		// Flag description (help flag) will show these values if they are set
		//  1. flag name
		//  2. environment variable for the flag
		//  3. variable description
		//  4. default value
		description := make([]string, 0, 3)

		var envTag string
		var envAliases []string
		var hasEnv bool

		for envKey, envField := range f.Env {
			if envField.Field == field.Field {
				envTag = envKey
				envAliases = envField.Aliases
				hasEnv = true
				break
			}
		}

		if hasEnv {
			var envVars strings.Builder
			envVars.WriteString(usageIndent)
			envVars.WriteString("ENV VARIABLES: ")
			envVars.WriteString((&ansi.ColoredText{
				Color: ansi.Green,
				Text:  envTag,
			}).String())
			for _, e := range envAliases {
				envVars.WriteString(", ")
				envVars.WriteString((&ansi.ColoredText{
					Color: ansi.Green,
					Text:  e,
				}).String())
			}
			description = append(description, envVars.String())
		}

		if field.Description != nil && *field.Description != "" {
			description = append(description, fmt.Sprintf("%s  DESCRIPTION: %s", usageIndent, *field.Description))
		}

		// If pointer, get type that is pointed to
		isPointer := field.StructField.Type.Kind() == reflect.Pointer
		actualType := field.StructField.Type
		if isPointer {
			actualType = actualType.Elem()
		}

		// Get default value that is being pointed at if field is a pointer
		Default := field.Default
		if isPointer && Default != nil {
			defaultType := reflect.TypeOf(Default)
			pointer := reflect.ValueOf(Default)
			if defaultType.Implements(convert.CustomParserType) && pointer.IsNil() {
				Default = reflect.New(reflect.TypeOf(Default).Elem()).Interface()
			} else if pointer.IsNil() {
				Default = reflect.Zero(reflect.TypeOf(Default).Elem()).Interface()
			}
		}

		defaultString := convert.AutoToString(Default)
		if defaultString != "" && !slices.Contains(f.Required, field.Field) {
			description = append(description,
				fmt.Sprintf(
					"%s      DEFAULT: %s",
					usageIndent,
					(&ansi.ColoredText{
						Color: ansi.Blue,
						Text:  defaultString,
					}).String(),
				),
			)
		}
		finalDescription := strings.Join(description, "\n")

		if reflect.TypeOf(Default).Implements(convert.CustomParserType) {
			var flagValue string
			for i := len(field.Aliases) - 1; i >= 0; i-- {
				flagSet.StringVar(&flagValue, field.Aliases[i], "", "")
			}
			flagSet.StringVar(&flagValue, flagName, "", finalDescription)
			postProcess[flagName] = func() error {
				if flagValue == "" {
					return nil
				}
				parsed, err := convert.AutoFromBytes(field.StructField.Type, []byte(flagValue))
				if err != nil {
					return format.Err("CustomParser.UnmarshalText failed %w", err)
				}

				field.Value.Set(reflect.ValueOf(parsed))
				return nil
			}
			continue
		}

		switch actualType {
		case reflect.TypeFor[bool]():
			var flagValue bool
			for _, a := range field.Aliases {
				flagSet.BoolVar(&flagValue, a, false, "")
			}
			flagSet.BoolVar(&flagValue, flagName, false, finalDescription)
			postProcess[flagName] = func() error {
				newValue := reflect.ValueOf(&flagValue)
				if isPointer {
					field.Value.Set(newValue)
				} else {
					field.Value.Set(newValue.Elem())
				}
				return nil
			}
		default:
			var flagValue string
			for _, a := range field.Aliases {
				flagSet.StringVar(&flagValue, a, "", "")
			}
			flagSet.StringVar(&flagValue, flagName, "", finalDescription)
			postProcess[flagName] = func() error {
				if flagValue == "" {
					return nil
				}
				result, err := convert.AutoFromBytes(field.StructField.Type, []byte(flagValue))
				if err != nil {
					return format.Err("Failed to convert from bytes %w", err)
				}

				field.Value.Set(reflect.ValueOf(result))

				return nil
			}
		}
	}

	HelpError.Message.Reset()
	flagSet.SetOutput(&HelpError.Message)
	flagSet.Usage = func() {
		writer := flagSet.Output()

		fmt.Fprintf(writer, "Usage of %s:", flagSet.Name())

		flagSet.VisitAll(func(flag *flag.Flag) {
			flagField, isFlag := f.Flags[flag.Name]
			if !isFlag {
				return
			}
			flagType := flagField.StructField.Type
			if flagType.Kind() == reflect.Pointer {
				flagType = flagType.Elem()
			}

			var output strings.Builder
			output.Write([]byte("\n\n"))
			output.WriteString((&ansi.ColoredText{Color: ansi.Green, Text: "-" + flag.Name}).String())

			for _, a := range flagField.Aliases {
				output.WriteByte(' ')
				output.WriteString((&ansi.ColoredText{Color: ansi.Green, Text: "-" + a}).String())
			}

			flagTypeString := flagType.String()
			output.WriteString(usageSeparator)
			output.WriteString((&ansi.ColoredText{Color: ansi.Yellow, Text: flagTypeString}).String())

			if slices.Contains(f.Required, flagField.Field) {
				output.WriteString(usageSeparator)
				output.WriteString((&ansi.ColoredText{Color: ansi.Red, Text: "[REQUIRED]"}).String())
			}

			output.WriteByte('\n')
			output.WriteString(flag.Usage)

			fmt.Fprint(writer, indentAllLines(output.String()))
		})
	}

	err := flagSet.Parse(args)
	if errors.Is(err, flag.ErrHelp) {
		return HelpError
	} else if err != nil {
		return err
	}

	for flagName, f := range postProcess {
		if err := f(); err != nil {
			return format.Err("Post process failed for flag -%s %w", flagName, err)
		}
	}

	return nil
}

func indentAllLines(input string) string {
	var result strings.Builder
	for l := range strings.Lines(input) {
		result.WriteString(usageIndent)
		result.WriteString(l)
	}
	return result.String()
}
