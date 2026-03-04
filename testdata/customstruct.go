package testdata

import (
	"strings"
)

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
