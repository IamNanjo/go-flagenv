package testdata

// CustomParser types must be pointers
type InvalidFields struct {
	CustomStruct CustomStruct `flag:"customStruct" env:"CUSTOM_STRUCT" desc:"Custom struct" required:"true" default:"true"`
}
