package testdata

import (
	"testing"
)

func VerifyAllTypes(t *testing.T, config *AllTypes) {
	if v := config.Bool; !v {
		t.Errorf("Expected Bool to be true. Got %v", v)
	}
	if v, cv := config.Int, -10; v != cv {
		t.Errorf("Expected Int to be %v. Got %v", cv, v)
	}
	if v, cv := config.Int64, int64(-50); v != cv {
		t.Errorf("Expected Int64 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Uint, uint(30); v != cv {
		t.Errorf("Expected Uint to be %v. Got %v", cv, v)
	}
	if v, cv := config.Uint64, uint64(70); v != cv {
		t.Errorf("Expected Uint64 to be %v. Got %v", cv, v)
	}
	if v, cv := config.Float64, 80.90; v != cv {
		t.Errorf("Expected Float64 to be %v. Got %v", cv, v)
	}
	if v, cv := config.String, "string"; v != cv {
		t.Errorf(`Expected String to be %q. Got %v`, cv, v)
	}

	if v := config.StringSlice; v == nil {
		t.Errorf("Expected StringSlice to be non-nil. Got nil")
	} else if v, cv := len(config.StringSlice), 3; v != cv {
		t.Errorf("Expected StringSlice to have %v strings. Got %+v", cv, v)
	} else if v, cv := config.StringSlice[0], "value1"; v != cv {
		t.Errorf(`Expected StringSlice[0] to be %v. Got %v`, cv, v)
	} else if v, cv := config.StringSlice[1], "value2"; v != cv {
		t.Errorf(`Expected StringSlice[1] to be %v. Got %v`, cv, v)
	} else if v, cv := config.StringSlice[2], "value3"; v != cv {
		t.Errorf(`Expected StringSlice[2] to be %v. Got %v`, cv, v)
	}

	if v := config.IntSlice; v == nil {
		t.Errorf(`Expected IntSlicePtr to be non-nil. Got nil`)
	} else if v, cv := len(config.IntSlice), 7; v != cv {
		t.Errorf("Expected IntSlicePtr to have %v integers. Got %+v", cv, v)
	} else if v, cv := (config.IntSlice)[0], -3; v != cv {
		t.Errorf(`Expected IntSlicePtr[0] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[1], -2; v != cv {
		t.Errorf(`Expected IntSlicePtr[1] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[2], -1; v != cv {
		t.Errorf(`Expected IntSlicePtr[2] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[3], 0; v != cv {
		t.Errorf(`Expected IntSlicePtr[3] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[4], 1; v != cv {
		t.Errorf(`Expected IntSlicePtr[4] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[5], 2; v != cv {
		t.Errorf(`Expected IntSlicePtr[5] to be %v. Got %v`, cv, v)
	} else if v, cv := (config.IntSlice)[6], 3; v != cv {
		t.Errorf(`Expected IntSlicePtr[6] to be %v. Got %v`, cv, v)
	}

	if v := config.CustomStructPtr; v == nil {
		t.Errorf(`Expected CustomStructPtr to be non-nil. Got nil`)
	} else if v := config.CustomStructPtr.Bool; !v {
		t.Errorf(`Expected CustomStructPtr.Bool to be true. Got %v`, v)
	}

	if v, cv := config.NestedContent.Int, -10; v != cv {
		t.Errorf("Expected Int to be %v. Got %v", cv, v)
	}
	if v, cv := config.NestedContent.String, "string"; v != cv {
		t.Errorf("Expected String to be %v. Got %v", cv, v)
	}
}
