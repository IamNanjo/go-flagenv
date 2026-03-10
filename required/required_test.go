package required_test

import (
	"testing"

	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/required"
	"github.com/IamNanjo/go-flagenv/testdata"
)

func TestRequired(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	config := new(testdata.AllTypes)
	f, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed %v", err)
	}

	err = required.Check(config, f)
	if err == nil {
		t.Fatalf("Required field check passed when expected to fail")
	}

	config.IntSlicePtr = new(testdata.IntSlice)
	err = config.IntSlicePtr.UnmarshalText([]byte("1,2,3"))
	if err != nil {
		t.Fatalf("Failed to create IntSlice from string %v", err)
	}

	config.IntPtr = new(-10)
	config.String = "Required string"
	config.CustomStructPtr = new(testdata.CustomStruct{Bool: true})

	f, err = fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed %v", err)
	}

	err = required.Check(config, f)
	if err != nil {
		t.Fatalf("Required field check failed with fields set: %+v", config)
	}
}
