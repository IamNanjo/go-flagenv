# Go flagenv

[![Go](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml/badge.svg)](https://github.com/IamNanjo/go-flagenv/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/IamNanjo/go-flagenv)](https://goreportcard.com/report/github.com/IamNanjo/go-flagenv)

Uses reflection and struct tags to parse flags and environment variables (also from .env file) with
optional defaults and support for required variables that will cause the parser to panic if they are missing or empty

## Usage

```go
// Example struct with all supported types (all of the native types can be pointers).
// IntSlice and CustomStruct implement convert.CustomParser which allows using any type as long as it's a pointer.
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
```

<details>
<summary>IntSlice definition</summary>

```go
type IntSlice []int

func (is *IntSlice) FromString(input string) error {
    split := strings.Split(input, ",")
    result := make(IntSlice, 0, len(split))
    for _, s := range split {
        trimmed := strings.TrimSpace(s)
        if trimmed == "" {
            continue
        }
        parsed, err := strconv.Atoi(trimmed)
        if err != nil {
            *is = result
            return format.Err("IntSlice parsing failed %w", err)
        }
        result = append(result, parsed)
    }
    *is = result
    return nil
}

func (is *IntSlice) String() string {
    if is == nil || len(*is) == 0 {
        return ""
    }

    var result strings.Builder
    for i, num := range *is {
        result.WriteString(strconv.Itoa(num))
        if i < len(*is)-1 {
            result.WriteString(", ")
        }
    }

    return result.String()
}
```
</details>

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
    go test -count=1 -failfast -v ./...
    ```
