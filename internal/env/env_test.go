package env_test

import (
	"testing"

	"github.com/IamNanjo/go-flagenv/internal/env"
	"github.com/IamNanjo/go-flagenv/internal/fields"
	"github.com/IamNanjo/go-flagenv/testdata"
)

func TestEnv(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	t.Setenv("BOOL", "true")
	t.Setenv("INT", "-10")
	t.Setenv("INT64", "-20")
	t.Setenv("UINT", "30")
	t.Setenv("UINT64", "40")
	t.Setenv("FLOAT64", "50.60")
	t.Setenv("STRING", "string")
	t.Setenv("STRING_SLICE", "value1,value2,value3")
	t.Setenv("INT_SLICE_PTR", "-3, -2, -1, 0, 1, 2, 3")
	t.Setenv("CUSTOM_STRUCT_PTR", "true")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err := env.Parse(config, fields); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	testdata.VerifyAllTypes(t, config)
}
