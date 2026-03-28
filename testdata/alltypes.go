package testdata

import "time"

// Uses all the supported types. Should parse correctly.
type AllTypes struct {
	Bool            bool          `flag:"bool"            env:"BOOL"              default:"true"      desc:"Boolean"`
	Int             int           `flag:"int"             env:"INT"               default:"1000"      desc:"Int"`
	IntPtr          *int          `flag:"intPtr"          env:"INT_PTR"           required:"true"     desc:"Int pointer"`
	Int8            int8          `flag:"int8"            env:"INT8"              default:"64"        desc:"8-bit int"`
	Int16           int16         `flag:"int16"           env:"INT16"             default:"128"       desc:"16-bit int"`
	Int32           int32         `flag:"int32"           env:"INT32"             default:"256"       desc:"32-bit int"`
	Int64           int64         `flag:"int64"           env:"INT64"             default:"512"       desc:"64-bit int"`
	Uint            uint          `flag:"uint"            env:"UINT"              default:"1000"      desc:"Uint"`
	Uint8           uint8         `flag:"uint8"           env:"UINT8"             default:"64"        desc:"8-bit uint"`
	Uint16          uint16        `flag:"uint16"          env:"UINT16"            default:"128"       desc:"16-bit uint"`
	Uint32          uint32        `flag:"uint32"          env:"UINT32"            default:"256"       desc:"32-bit uint"`
	Uint64          uint64        `flag:"uint64"          env:"UINT64"            default:"512"       desc:"64-bit uint"`
	Float64         float64       `flag:"float64"         env:"FLOAT64"           default:"1000"      desc:"64-bit float"`
	String          string        `flag:"string"          env:"STRING"            required:"true"     desc:"String"`
	ByteSlice       []byte        `flag:"byteSlice"       env:"BYTE_SLICE"        required:"true"     desc:"Byte slice"`
	StringSlice     []string      `flag:"stringSlice"     env:"STRING_SLICE"      default:"val1,val2" desc:"String slice"`
	IntSlice        []int         `flag:"intSlice"        env:"INT_SLICE"         default:"100,101"   desc:"Int slice"`
	IntSlicePtr     *[]int        `flag:"intSlicePtr"     env:"INT_SLICE_PTR"     default:"100,101"   desc:"Int slice pointer"`
	IntPtrSlicePtr  *[]*int       `flag:"intPtrSlicePtr"  env:"INT_PTR_SLICE_PTR" default:"100,101"   desc:"Pointer to slice of int pointers"`
	Duration        time.Duration `flag:"duration"        env:"DURATION"          default:"1h1m1s"    desc:"Time duration"`
	CustomStructPtr *CustomStruct `flag:"customStructPtr" env:"CUSTOM_STRUCT_PTR" required:"true"     desc:"Custom struct pointer"`
	NestedContent   NestedContent `flag:"nested"          env:"NESTED_"`
}

type NestedContent struct {
	Int    int    `flag:"Int" env:"INT"`
	String string `flag:"String" env:"STRING"`
}
