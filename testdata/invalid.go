package testdata

// CustomParser types must be pointers
type InvalidFields struct {
	IntSlice IntSlice `flag:"intSlice" env:"INT_SLICE" desc:"Int slice" required:"true" default:"1000,1001,1002"`
}
