package convert

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Map of supported type to FromString function
var FromString = map[reflect.Type]func(input []byte) (any, error){
	reflect.TypeFor[bool](): func(input []byte) (any, error) {
		return strings.ToLower(string(input)) == "true", nil
	},
	reflect.TypeFor[int](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseInt(string(input), 0, 0)
		return int(parsed), err
	},
	reflect.TypeFor[int8](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseInt(string(input), 0, 8)
		return int8(parsed), err
	},
	reflect.TypeFor[int16](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseInt(string(input), 0, 16)
		return int16(parsed), err
	},
	reflect.TypeFor[int32](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseInt(string(input), 0, 32)
		return int32(parsed), err
	},
	reflect.TypeFor[int64](): func(input []byte) (any, error) {
		return strconv.ParseInt(string(input), 0, 64)
	},
	reflect.TypeFor[uint](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseUint(string(input), 0, 0)
		return uint(parsed), err
	},
	reflect.TypeFor[uint8](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseUint(string(input), 0, 8)
		return uint8(parsed), err
	},
	reflect.TypeFor[uint16](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseUint(string(input), 0, 16)
		return uint16(parsed), err
	},
	reflect.TypeFor[uint32](): func(input []byte) (any, error) {
		parsed, err := strconv.ParseUint(string(input), 0, 32)
		return uint32(parsed), err
	},
	reflect.TypeFor[uint64](): func(input []byte) (any, error) {
		return strconv.ParseUint(string(input), 0, 64)
	},
	reflect.TypeFor[float32](): func(input []byte) (any, error) {
		return strconv.ParseFloat(string(input), 32)
	},
	reflect.TypeFor[float64](): func(input []byte) (any, error) {
		return strconv.ParseFloat(string(input), 64)
	},
	reflect.TypeFor[[]byte](): func(input []byte) (any, error) {
		return input, nil
	},
	reflect.TypeFor[string](): func(input []byte) (any, error) {
		return string(input), nil
	},
	reflect.TypeFor[time.Duration](): func(input []byte) (any, error) {
		return time.ParseDuration(string(input))
	},
}
