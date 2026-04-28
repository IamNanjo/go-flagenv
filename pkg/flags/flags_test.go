package flags_test

import (
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/IamNanjo/go-flagenv"
	"github.com/IamNanjo/go-flagenv/pkg/fields"
	"github.com/IamNanjo/go-flagenv/pkg/flags"
	"github.com/IamNanjo/go-flagenv/testdata"
)

var args = []string{
	"-bool",
	"-int", "-10",
	"-intPtr", "-10",
	"-int8", "-20",
	"-int16", "-30",
	"-int32", "-40",
	"-int64", "-50",
	"-uint", "30",
	"-uint8", "40",
	"-uint16", "50",
	"-uint32", "60",
	"-uint64", "70",
	"-float32", "80.90",
	"-float64", "80.90",
	"-string", "string",
	"-stringSlice", "value1,value2,value3",
	"-intSlice", "-3, -2, -1, 0, 1, 2, 3",
	"-intSlicePtr", "-3, -2, -1, 0, 1, 2, 3",
	"-intPtrSlicePtr", "-3, -2, -1, 0, 1, 2, 3",
	"-customStructPtr", "true",
	"-nestedInt", "-10",
	"-nestedString", "string",
}

func TestFlags(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err = flags.Parse(fields, args); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	testdata.VerifyAllTypes(t, config)
}

func TestHelp(t *testing.T) {
	err := flagenv.ParseCustom(new(testdata.AllTypes), []string{"-help"}, "")
	fmt.Printf("%v", err)

	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("Flag parsing did not return HelpError")
	}
}
