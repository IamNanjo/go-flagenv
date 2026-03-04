package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IamNanjo/go-flagenv/convert"
	"github.com/IamNanjo/go-flagenv/internal/fields"
	"github.com/IamNanjo/go-flagenv/internal/format"
)

const usageIndent = 6
const usageSeparator = "  ·  "

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

var HelpError = new(helpError)

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

		envTag, hasEnv := field.StructField.Tag.Lookup("env")
		if hasEnv {
			description = append(description, fmt.Sprintf("%s(ENV: %s)", usageSeparator, envTag))
		} else {
			description = append(description, "")
		}

		if field.Description != nil && *field.Description != "" {
			description = append(description, fmt.Sprintf("  DESCRIPTION: %s", *field.Description))
		}

		// If pointer, get type that is pointed to
		isPointer := field.StructField.Type.Kind() == reflect.Pointer
		actualType := field.StructField.Type
		if isPointer {
			actualType = actualType.Elem()
		}
		kind := actualType.Kind()
		shouldDereference := !isPointer && kind != reflect.Slice && kind != reflect.Map

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

		// Pointer to any supported value received from flag
		var flagValue any = nil

		var defaultStringBuilder strings.Builder

		defaultType := reflect.TypeOf(Default)
		if defaultType.Implements(convert.CustomParserType) {
			parser, ok := Default.(convert.CustomParser)
			if !ok {
				return format.Err("Could not use field %q as CustomParser", field.StructField.Name)
			}
			parsedString := parser.String()
			if len(parsedString) == 0 {
				defaultStringBuilder.WriteString("<empty>")
			} else {
				defaultStringBuilder.WriteString(parsedString)
			}
			description = append(description, "DEFAULT VALUE: "+defaultStringBuilder.String())

			finalDescription := strings.Join(description, "\n")

			flagValue = flagSet.String(flagName, parser.String(), finalDescription)
			postProcess[flagName] = func() error {
				err := parser.FromString(*flagValue.(*string))
				resultType := reflect.TypeOf(parser)
				if resultType != field.StructField.Type {
					return format.Err("CustomParser.FromString returned %q. Expected %q", resultType, field.StructField.Type)
				}
				if err != nil {
					return format.Err("CustomParser.FromString failed %w", err)
				}
				field.Value.Set(reflect.ValueOf(parser))
				return nil
			}

			continue
		}

		switch field.StructField.Type {
		case convert.IntSliceType:
			d := Default.([]int)
			if len(d) == 0 {
				defaultStringBuilder.WriteString("<empty>")
			}
			for i, s := range d {
				defaultStringBuilder.WriteString(strconv.Itoa(s))

				if i < len(d)-1 {
					defaultStringBuilder.WriteString(", ")
				}
			}
			defaultStringBuilder.WriteString(" (%+v)")
		case convert.StringSliceType:
			d := Default.([]string)
			if len(d) == 0 {
				defaultStringBuilder.WriteString("<empty>")
			}
			for i, s := range d {
				defaultStringBuilder.WriteString(s)

				if i < len(d)-1 {
					defaultStringBuilder.WriteString(", ")
				}
			}
			defaultStringBuilder.WriteString(" (%+v)")
		case reflect.TypeFor[string]():
			defaultStringBuilder.WriteString("%q")
		default:
			defaultStringBuilder.WriteString("%+v")
		}

		defaultString := fmt.Sprintf(defaultStringBuilder.String(), Default)
		description = append(description, "DEFAULT VALUE: "+defaultString)

		finalDescription := strings.Join(description, "\n")

		switch actualType {
		case reflect.TypeFor[bool]():
			flagValue = flagSet.Bool(flagName, false, finalDescription)
		case reflect.TypeFor[int]():
			flagValue = flagSet.Int(flagName, 0, finalDescription)
		case reflect.TypeFor[int64]():
			flagValue = flagSet.Int64(flagName, 0, finalDescription)
		case reflect.TypeFor[uint]():
			flagValue = flagSet.Uint(flagName, 0, finalDescription)
		case reflect.TypeFor[uint64]():
			flagValue = flagSet.Uint64(flagName, 0, finalDescription)
		case reflect.TypeFor[float64]():
			flagValue = flagSet.Float64(flagName, 0, finalDescription)
		case reflect.TypeFor[string]():
			flagValue = flagSet.String(flagName, "", finalDescription)
		case reflect.TypeFor[time.Duration]():
			flagValue = flagSet.Duration(flagName, time.Duration(0), finalDescription)
		case convert.StringSliceType:
			flagValue := flagSet.String(flagName, convert.FromStringSlice(Default.([]string)), finalDescription)
			postProcess[flagName] = func() error {
				if *flagValue == "" {
					return nil
				}
				result := convert.ToStringSlice(*flagValue)
				newValue := reflect.ValueOf(result)
				if !shouldDereference {
					field.Value.Set(newValue)
				} else {
					field.Value.Set(newValue.Elem())
				}
				return nil
			}
			continue
		default:
			return format.Err("Unsupported type %q for flag %q", field.StructField.Type, flagName)
		}

		postProcess[flagName] = func() error {
			newValue := reflect.ValueOf(flagValue).Elem()
			if !shouldDereference {
				field.Value.Set(reflect.New(newValue.Type()))
			} else {
				field.Value.Set(newValue)
			}
			return nil
		}
	}

	flagSet.SetOutput(&HelpError.Message)
	flagSet.Usage = func() {
		writer := flagSet.Output()

		fmt.Fprintf(writer, "Usage of %s:", flagSet.Name())
		flagSet.VisitAll(func(f *flag.Flag) {
			fType := fmt.Sprintf("%T", f.Value)
			fType = strings.TrimPrefix(fType, "*flag.")
			fType = strings.TrimSuffix(fType, "Value")

			fmt.Fprintf(writer, "\n\n  -%s%s%s%s", f.Name, usageSeparator, fType, format.IndentAllLines(f.Usage, usageIndent)[usageIndent:])
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
