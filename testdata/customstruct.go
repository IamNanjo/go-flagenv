package testdata

import (
	"strings"
)

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
