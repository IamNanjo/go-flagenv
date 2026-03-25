package required_test

import (
	"testing"

	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/required"
	"github.com/IamNanjo/go-flagenv/testdata"
)

func TestRequired(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	// Check before filling required fields.
	config := new(testdata.AllTypes)
	f, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed %v", err)
	}

	err = required.Check(config, f)
	if err == nil {
		t.Fatalf("Required field check passed when expected to fail")
	}

	// Fill required fields.
	config.IntPtr = new(-10)
	config.String = "Required string"
	config.CustomStructPtr = new(testdata.CustomStruct{Bool: true})
	config.ByteSlice = []byte("Required byte slice")

	// Ensure required field check passes with required fields set.
	f, err = fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed %v", err)
	}

	err = required.Check(config, f)
	if err != nil {
		t.Fatalf("Required field check failed with fields set: %v", err)
	}
}
