# Go flagenv

[![Go](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml/badge.svg)](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/IamNanjo/go-flagenv)](https://goreportcard.com/report/github.com/IamNanjo/go-flagenv)

Uses reflection and struct tags to parse flags and environment variables (also from .env file) with
optional defaults and support for required variables that will cause the parser to panic if they are missing or empty.

## Usage

```go
// Example struct with all supported types (all of the native types can be pointers or slices).
// CustomStruct implements pkg/convert.CustomParser which allows using any type as long as it's a pointer.
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


// Parse with automatic handling of -help flag and look for .env file in current working directory.
// Passes CLI flags provided to program for flag parsing.
func main() {
    config := new(AllTypes)
    if err := flagenv.Parse(config); err != nil {
	// Parsing failed. Handle error here
    }
}

// Pass custom CLI flags and .env file path
func main() {
    config := new(AllTypes)
    if err := flagenv.ParseCustom(config, os.Args[1:], ".env"); err != nil {
	// If err is flag.ErrHelp then help menu was called and program should typically exit after printing it.
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Handle error
    }
}
```

<details>
<summary>CustomStruct definition</summary>

```go
type CustomStruct struct{ Bool bool }

func (cs *CustomStruct) UnmarshalText(input []byte) error {
	cs.Bool = strings.ToLower(string(input)) == "true"
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

<details>
<summary>Help menu generated from AllTypes</summary>

```plaintext
Usage of flagenv:

    -bool · bool
        ENV VARIABLE: BOOL
         DESCRIPTION: Boolean
             DEFAULT: true

    -byteSlice · []uint8 · [REQUIRED]
        ENV VARIABLE: BYTE_SLICE
         DESCRIPTION: Byte slice

    -customStructPtr · testdata.CustomStruct · [REQUIRED]
        ENV VARIABLE: CUSTOM_STRUCT_PTR
         DESCRIPTION: Custom struct pointer

    -duration · time.Duration
        ENV VARIABLE: DURATION
         DESCRIPTION: Time duration
             DEFAULT: 1h1m1s

    -float64 · float64
        ENV VARIABLE: FLOAT64
         DESCRIPTION: 64-bit float
             DEFAULT: 1000.00

    -int · int
        ENV VARIABLE: INT
         DESCRIPTION: Int
             DEFAULT: 1000

    -int16 · int16
        ENV VARIABLE: INT16
         DESCRIPTION: 16-bit int
             DEFAULT: 128

    -int32 · int32
        ENV VARIABLE: INT32
         DESCRIPTION: 32-bit int
             DEFAULT: 256

    -int64 · int64
        ENV VARIABLE: INT64
         DESCRIPTION: 64-bit int
             DEFAULT: 512

    -int8 · int8
        ENV VARIABLE: INT8
         DESCRIPTION: 8-bit int
             DEFAULT: 64

    -intPtr · int · [REQUIRED]
        ENV VARIABLE: INT_PTR
         DESCRIPTION: Int pointer

    -intPtrSlicePtr · []*int
        ENV VARIABLE: INT_PTR_SLICE_PTR
         DESCRIPTION: Pointer to slice of int pointers
             DEFAULT: 100, 101

    -intSlice · []int
        ENV VARIABLE: INT_SLICE
         DESCRIPTION: Int slice
             DEFAULT: 100, 101

    -intSlicePtr · []int
        ENV VARIABLE: INT_SLICE_PTR
         DESCRIPTION: Int slice pointer
             DEFAULT: 100, 101

    -nestedInt · int
        ENV VARIABLE: NESTED_INT
             DEFAULT: 0

    -nestedString · string
        ENV VARIABLE: NESTED_STRING

    -string · string · [REQUIRED]
        ENV VARIABLE: STRING
         DESCRIPTION: String

    -stringSlice · []string
        ENV VARIABLE: STRING_SLICE
         DESCRIPTION: String slice
             DEFAULT: val1, val2

    -uint · uint
        ENV VARIABLE: UINT
         DESCRIPTION: Uint
             DEFAULT: 1000

    -uint16 · uint16
        ENV VARIABLE: UINT16
         DESCRIPTION: 16-bit uint
             DEFAULT: 128

    -uint32 · uint32
        ENV VARIABLE: UINT32
         DESCRIPTION: 32-bit uint
             DEFAULT: 256

    -uint64 · uint64
        ENV VARIABLE: UINT64
         DESCRIPTION: 64-bit uint
             DEFAULT: 512

    -uint8 · uint8
        ENV VARIABLE: UINT8
         DESCRIPTION: 8-bit uint
             DEFAULT: 64
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
