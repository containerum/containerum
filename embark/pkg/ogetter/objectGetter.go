package ogetter

import (
	"io"
)

type ObjectGetter interface {
	Object(name string, output io.Writer) error
}
