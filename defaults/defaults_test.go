package defaults_test

import (
	"testing"
	"time"

	"github.com/IamNanjo/go-flagenv/defaults"
	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/testdata"
)

func TestDefault(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed %v", err)
	}

	defaults.Set(config, fields)

	if v, cv := config.Bool, true; v != cv {
		t.Fatalf("Expected Bool to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Int, 1000; v != cv {
		t.Fatalf("Expected Int to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Int8, int8(64); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Int16, int16(128); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Int32, int32(256); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Int64, int64(512); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint, uint(1000); v != cv {
		t.Fatalf("Expected Uint to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint8, uint8(64); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint16, uint16(128); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint32, uint32(256); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint64, uint64(512); v != cv {
		t.Fatalf("Expected Uint64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Float64, float64(1000); v != cv {
		t.Fatalf("Expected Float64 to be %+v. Got %+v", cv, v)
	}

	if v := config.StringSlice; v == nil {
		t.Fatalf("Expected StringSlice to be non-nil. Got nil")
	} else if v, cv := len(config.StringSlice), 2; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := config.StringSlice[0], "val1"; v != cv {
		t.Fatalf("Expected StringSlice[0] to be %+v. Got %+v", cv, v)
	} else if v, cv := config.StringSlice[1], "val2"; v != cv {
		t.Fatalf("Expected StringSlice[1] to be %+v strings. Got %+v", cv, v)
	}

	if v, cv := config.Duration, time.Duration(time.Hour+time.Minute+time.Second); v != cv {
		t.Fatalf("Expected Duration to be %+v. Got %+v", cv, v)
	}

	if v := config.IntSlice; v == nil {
		t.Fatalf("Expected IntSlice to be non-nil. Got nil")
	} else if v, cv := len(config.IntSlice), 2; v != cv {
		t.Fatalf("Expected IntSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := (config.IntSlice)[0], 100; v != cv {
		t.Fatalf("Expected IntSlice[0] to be %+v. Got %+v", cv, v)
	} else if v, cv := (config.IntSlice)[1], 101; v != cv {
		t.Fatalf("Expected IntSlice[1] to be %+v. Got %+v", cv, v)
	}
}
