package defaults_test

import (
	"testing"
	"time"

	"github.com/IamNanjo/go-flagenv/internal/defaults"
	"github.com/IamNanjo/go-flagenv/internal/fields"
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
	if v, cv := config.Int64, int64(1000); v != cv {
		t.Fatalf("Expected Int64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint, uint(1000); v != cv {
		t.Fatalf("Expected Uint to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Uint64, uint64(1000); v != cv {
		t.Fatalf("Expected Uint64 to be %+v. Got %+v", cv, v)
	}
	if v, cv := config.Float64, float64(1000); v != cv {
		t.Fatalf("Expected Float64 to be %+v. Got %+v", cv, v)
	}

	if v := config.StringSlice; v == nil {
		t.Fatalf("Expected StringSlice to be non-nil. Got nil")
	} else if v, cv := len(config.StringSlice), 2; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := config.StringSlice[0], "default1"; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := config.StringSlice[1], "default2"; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	}

	if v, cv := config.Duration, time.Duration(time.Hour+time.Minute+time.Second); v != cv {
		t.Fatalf("Expected Duration to be %+v. Got %+v", cv, v)
	}

	if v := config.IntSlicePtr; v == nil {
		t.Fatalf("Expected StringSlice to be non-nil. Got nil")
	} else if v, cv := len(*config.IntSlicePtr), 3; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := (*config.IntSlicePtr)[0], 1000; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := (*config.IntSlicePtr)[1], 1001; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	} else if v, cv := (*config.IntSlicePtr)[2], 1002; v != cv {
		t.Fatalf("Expected StringSlice to have %+v strings. Got %+v", cv, v)
	}

}
