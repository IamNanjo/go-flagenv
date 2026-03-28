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

	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/internal/convert"

	"github.com/IamNanjo/go-logging/pkg/ansi"
	"github.com/IamNanjo/go-logging/pkg/format"
)

const usageIndent = 4
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

func Parse[T any](c *T, f *fields.Fields, args []string) error {
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
		var hasEnv bool

		for envKey, envField := range f.Env {
			if envField == field {
				envTag = envKey
				hasEnv = true
				break
			}
		}

		if hasEnv {
			description = append(description,
				fmt.Sprintf("ENV VARIABLE: %s", (&ansi.ColoredText{
					Color: ansi.Green,
					Text:  envTag,
				}).String()),
			)
		}

		if field.Description != nil && *field.Description != "" {
			description = append(description, fmt.Sprintf(" DESCRIPTION: %s", *field.Description))
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

		// Pointer to value received from flag
		var flagValue any = nil

		defaultString := convert.AutoToString(Default)
		if defaultString != "" && !slices.Contains(f.Required, field) {
			description = append(description, "     DEFAULT: "+defaultString)
		}
		finalDescription := strings.Join(description, "\n")

		if reflect.TypeOf(Default).Implements(convert.CustomParserType) {
			flagValue = flagSet.String(flagName, "", finalDescription)
			postProcess[flagName] = func() error {
				parsed, err := convert.AutoFromBytes(field.StructField.Type, []byte(*flagValue.(*string)))
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
			flagValue = flagSet.Bool(flagName, false, finalDescription)
			postProcess[flagName] = func() error {
				val := (flagValue.(*bool))
				newValue := reflect.ValueOf(val)
				if isPointer {
					field.Value.Set(newValue)
				} else {
					field.Value.Set(newValue.Elem())
				}
				return nil
			}
		default:
			flagValue = flagSet.String(flagName, "", finalDescription)
			postProcess[flagName] = func() error {
				val := *(flagValue.(*string))
				if val == "" {
					return nil
				}
				result, err := convert.AutoFromBytes(field.StructField.Type, []byte(val))
				if err != nil {
					return format.Err("Failed to convert from bytes %w", err)
				}

				field.Value.Set(reflect.ValueOf(result))

				return nil
			}
		}
	}

	flagSet.SetOutput(&HelpError.Message)
	flagSet.Usage = func() {
		writer := flagSet.Output()

		fmt.Fprintf(writer, "Usage of %s:", flagSet.Name())

		flagSet.VisitAll(func(flag *flag.Flag) {
			flagField := f.Flags[flag.Name]
			flagType := flagField.StructField.Type
			if flagType.Kind() == reflect.Pointer {
				flagType = flagType.Elem()
			}

			var output strings.Builder
			output.Write([]byte("\n\n"))
			output.WriteString((&ansi.ColoredText{Color: ansi.Green, Text: "-" + flag.Name}).String())

			flagTypeString := flagType.String()
			output.WriteString(usageSeparator)
			output.WriteString((&ansi.ColoredText{Color: ansi.Yellow, Text: flagTypeString}).String())

			if slices.Contains(f.Required, flagField) {
				output.WriteString(usageSeparator)
				output.WriteString((&ansi.ColoredText{Color: ansi.Red, Text: "[REQUIRED]"}).String())
			}

			output.WriteByte('\n')
			output.WriteString(indentAllLines(flag.Usage, usageIndent))

			fmt.Fprint(writer, output.String())
		})
		fmt.Fprintln(writer)
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

func indentAllLines(input string, indentWidth int) string {
	indent := strings.Repeat(" ", indentWidth)

	var result strings.Builder
	for l := range strings.Lines(input) {
		result.WriteString(indent)
		result.WriteString(l)
	}

	return result.String()
}
