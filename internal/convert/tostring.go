package convert

import (
	"reflect"
	"strconv"
	"time"
)

// Map of supported type to ToString function
var ToString = map[reflect.Type]func(input any) string{
	reflect.TypeFor[bool](): func(input any) string {
		if input.(bool) {
			return "true"
		}
		return "false"
	},
	reflect.TypeFor[int](): func(input any) string {
		return strconv.FormatInt(int64(input.(int)), 10)
	},
	reflect.TypeFor[int8](): func(input any) string {
		return strconv.FormatInt(int64(input.(int8)), 10)
	},
	reflect.TypeFor[int16](): func(input any) string {
		return strconv.FormatInt(int64(input.(int16)), 10)
	},
	reflect.TypeFor[int32](): func(input any) string {
		return strconv.FormatInt(int64(input.(int32)), 10)
	},
	reflect.TypeFor[int64](): func(input any) string {
		return strconv.FormatInt(input.(int64), 10)
	},
	reflect.TypeFor[uint](): func(input any) string {
		return strconv.FormatUint(uint64(input.(uint)), 10)
	},
	reflect.TypeFor[uint8](): func(input any) string {
		return strconv.FormatUint(uint64(input.(uint8)), 10)
	},
	reflect.TypeFor[uint16](): func(input any) string {
		return strconv.FormatUint(uint64(input.(uint16)), 10)
	},
	reflect.TypeFor[uint32](): func(input any) string {
		return strconv.FormatUint(uint64(input.(uint32)), 10)
	},
	reflect.TypeFor[uint64](): func(input any) string {
		return strconv.FormatUint(input.(uint64), 10)
	},
	reflect.TypeFor[float32](): func(input any) string {
		return strconv.FormatFloat(float64(input.(float32)), 'g', 2, 32)
	},
	reflect.TypeFor[float64](): func(input any) string {
		return strconv.FormatFloat(input.(float64), 'g', 2, 64)
	},
	reflect.TypeFor[[]byte](): func(input any) string {
		return string(input.([]byte))
	},
	reflect.TypeFor[string](): func(input any) string {
		return input.(string)
	},
	reflect.TypeFor[time.Duration](): func(input any) string {
		return input.(time.Duration).String()
	},
}
