package env_test

import (
	"testing"

	"github.com/IamNanjo/go-flagenv/pkg/env"
	"github.com/IamNanjo/go-flagenv/pkg/fields"
	"github.com/IamNanjo/go-flagenv/testdata"
)

func TestEnv(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	t.Setenv("BOOL", "true")
	t.Setenv("INT", "-10")
	t.Setenv("INT64", "-50")
	t.Setenv("UINT", "30")
	t.Setenv("UINT64", "70")
	t.Setenv("FLOAT32", "80.90")
	t.Setenv("FLOAT64", "80.90")
	t.Setenv("STRING", "string")
	t.Setenv("STRING_SLICE", "value1,value2,value3")
	t.Setenv("INT_SLICE", "-3, -2, -1, 0, 1, 2, 3")
	t.Setenv("INT_SLICE_PTR", "-3, -2, -1, 0, 1, 2, 3")
	t.Setenv("INT_PTR_SLICE_PTR", "-3, -2, -1, 0, 1, 2, 3")
	t.Setenv("CUSTOM_STRUCT_PTR", "true")
	t.Setenv("NESTED_INT", "-10")
	t.Setenv("NESTED_STRING", "string")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err := env.Parse(fields); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	testdata.VerifyAllTypes(t, config)
}
