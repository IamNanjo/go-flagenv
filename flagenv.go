package flagenv

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/IamNanjo/go-flagenv/defaults"
	"github.com/IamNanjo/go-flagenv/dotenv"
	"github.com/IamNanjo/go-flagenv/env"
	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/flags"
	"github.com/IamNanjo/go-flagenv/internal/format"
	"github.com/IamNanjo/go-flagenv/required"
	"github.com/IamNanjo/go-logging"
)

// Populate Config struct with CLI flags, .env file and environment variables in that priority order.
//
// Supported struct tags:
//
//	flag:"flagName"
//	env:"ENV_KEY"
//	default:"defaultValue"
//	required:"true"
//	desc:"description"
//
// Environment variables and slices will always be parsed from a string.
// Slices expect comma separated values and space around each value will have space around it trimmed.
//
// Supported variable types (any of the normal types can be pointers, slices or pointers to slices of those values):
//
//	bool
//	int(8/16/32/64)
//	uint(8/16/32/64)
//	float32
//	float64
//	string
//	[]byte
//	time.Duration
//	*convert.CustomParser (pointer to any other type that implements convert.CustomParser interface)
func Parse[T any](config *T) error {
	err := ParseCustom(config, os.Args[1:], ".env")
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	return err
}

// Same as Parse but allows a custom args source (CLI flags) and .env file path.
// Passing nil or an empty slice to args will skip parsing flags and empty .env path will skip parsing .env file.
// Error can be flags.HelpError in which case an error has not actually occurred and the user has asked for the help menu.
// If error is flags.HelpError, it should be printed to stderr and program should exit.
func ParseCustom[T any](config *T, args []string, envPath string) error {
	if config == nil {
		return format.Err("Config pointer is nil")
	}

	f, err := fields.Parse(config)
	if err != nil {
		return format.Err("Field parsing failed %w", err)
	}

	if len(args) != 0 {
		if err = flags.Parse(config, f, args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return err
			}
			return format.Err("Flag parsing failed %w", err)
		}
	} else {
		logging.Debug("No CLI arguments provided. Skipping...\n")
	}

	if envPath != "" {
		if err = dotenv.Parse(config, f, envPath); err != nil {
			return format.Err(".env parsing failed %w", err)
		}
	} else {
		logging.Debug("No .env path provided. Skipping...\n")
	}

	if err = env.Parse(config, f); err != nil {
		return format.Err("Environment variable parsing failed %w", err)
	}

	if err = defaults.Set(config, f); err != nil {
		return format.Err("Could not set defaults %w", err)
	}
	if err = required.Check(config, f); err != nil {
		return format.Err("Required variable check failed %w", err)
	}

	return nil
}
