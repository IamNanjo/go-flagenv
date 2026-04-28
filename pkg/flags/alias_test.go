package flags_test

import (
	"testing"

	"github.com/IamNanjo/go-flagenv/pkg/fields"
	"github.com/IamNanjo/go-flagenv/pkg/flags"
	"github.com/IamNanjo/go-flagenv/testdata"
)

var aliasArgs = []string{
	"-i32", "-40",
	"-i64", "-50",
	"-u32", "60",
	"-u64", "70",
	"-f64", "80.90",
}

func TestFlagAliases(t *testing.T) {
	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err = flags.Parse(config, fields, aliasArgs); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if v, cv := config.Int32, int32(-40); v != cv {
		t.Errorf("Expected Int32 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Int64, int64(-50); v != cv {
		t.Errorf("Expected Int64 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Uint32, uint32(60); v != cv {
		t.Errorf("Expected Uint32 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Uint64, uint64(70); v != cv {
		t.Errorf("Expected Uint64 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Float64, 80.90; v != cv {
		t.Errorf("Expected Float64 to be %v. Got %v", cv, v)
	}
}
