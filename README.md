# Go flagenv

[![Go](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml/badge.svg)](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/IamNanjo/go-flagenv)](https://goreportcard.com/report/github.com/IamNanjo/go-flagenv)

Uses reflection and struct tags to parse flags and environment variables (also from .env file) with
optional defaults and support for required variables that will cause the parser to panic if they are missing or empty.

## Usage

```go
// Example struct with all supported types (all of the native types can be pointers or slices).
// CustomStruct implements convert.CustomParser which allows using any type as long as it's a pointer.
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
	Duration        time.Duration `flag:"duration"        env:"DURATION"          default:"1h1m1s"    desc:"Time duration"`
	IntSlicePtr     *IntSlice     `flag:"intSlicePtr"     env:"INT_SLICE_PTR"     default:"100,101"   desc:"Int slice pointer"`
	CustomStructPtr *CustomStruct `flag:"customStructPtr" env:"CUSTOM_STRUCT_PTR" required:"true"     desc:"Custom struct pointer"`
}

// Parse with automatic handling of -help flag and look for .env file in current working directory.
// Passes CLI flags provided to program for flag parsing.
func main() {
    config := new(AllTypes)
    if err := flagenv.Parse(config); err != nil {
	// TODO: Parsing failed. Handle error here
    }
}

func main() {
    config := new(AllTypes)
    if err := flagenv.ParseCustom(); err != nil {
	// TODO: Parsing failed or help menu was requested.
	// If err is flag.HelpError then help menu was called and program should typically exit after printing it.
    }
}
```

<details>
<summary>CustomStruct definition</summary>

```go
type CustomStruct struct{ Bool bool }

func (cs *CustomStruct) FromString(input string) error {
	cs.Bool = strings.ToLower(input) == "true"
	return nil
}

func (cs *CustomStruct) String() string {
	if cs.Bool {
		return "true"
	}
	return "false"
}
```
</details>

## Development

Issues and pull requests are welcome

**Requires at least Go 1.26.0**

1. Install dependencies
    ```sh
    go mod download
    ```

1. Run tests
    ```sh
    go test -count=1 -failfast ./...
    ```
