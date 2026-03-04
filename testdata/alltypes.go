package testdata

import "time"

// Uses all the supported types. Should parse correctly.
type AllTypes struct {
	Bool            bool          `flag:"bool"            env:"BOOL"              desc:"Boolean"               default:"true"`
	Int             int           `flag:"int"             env:"INT"               desc:"Int"                   default:"1000"`
	IntPtr          *int          `flag:"intPtr"          env:"INT_PTR"           desc:"Int pointer"           required:"true"`
	Int64           int64         `flag:"int64"           env:"INT64"             desc:"64-bit int"            default:"1000"`
	Uint            uint          `flag:"uint"            env:"UINT"              desc:"Uint"                  default:"1000"`
	Uint64          uint64        `flag:"uint64"          env:"UINT64"            desc:"64-bit uint"           default:"1000"`
	Float64         float64       `flag:"float64"         env:"FLOAT64"           desc:"64-bit float"          default:"1000"`
	String          string        `flag:"string"          env:"STRING"            desc:"String"                required:"true"`
	StringSlice     []string      `flag:"stringSlice"     env:"STRING_SLICE"      desc:"String slice"          default:"default1,default2"`
	Duration        time.Duration `flag:"duration"        env:"DURATION"          desc:"Time duration"         default:"1h1m1s"`
	IntSlicePtr     *IntSlice     `flag:"intSlicePtr"     env:"INT_SLICE_PTR"     desc:"Int slice pointer"     default:"1000,1001,1002"`
	CustomStructPtr *CustomStruct `flag:"customStructPtr" env:"CUSTOM_STRUCT_PTR" desc:"Custom struct pointer" required:"true"`
}
