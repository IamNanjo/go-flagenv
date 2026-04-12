package convert

import (
	"encoding"
	"fmt"
)

type CustomParser interface {
	encoding.TextUnmarshaler
	fmt.Stringer
}
