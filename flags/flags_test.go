package flags_test

import (
	"errors"
	"testing"

	"github.com/IamNanjo/go-flagenv/flags"
	"github.com/IamNanjo/go-flagenv/internal/fields"
	"github.com/IamNanjo/go-flagenv/testdata"
)

var args = []string{
	"-bool",
	"-int", "-10",
	"-intPtr", "-10",
	"-int64", "-20",
	"-uint", "30",
	"-uint64", "40",
	"-float64", "50.60",
	"-string", "string",
	"-stringSlice", "value1,value2,value3",
	"-intSlicePtr", "-3, -2, -1, 0, 1, 2, 3",
	"-customStructPtr", "true",
}

func TestHelp(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	err = flags.Parse(config, fields, []string{"-help"})
	if !errors.Is(err, flags.HelpError) {
		t.Fatalf("Expected help menu. Got %v", err)
	}
}

func TestFlags(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err = flags.Parse(config, fields, args); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	testdata.VerifyAllTypes(t, config)
}
