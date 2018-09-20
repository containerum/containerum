package emberr

import (
	"bytes"

	"github.com/containerum/containerum/embark/pkg/utils/why"
)

var (
	_ Error = Chain{}
)

type Chain struct {
	head error
	tail []error
}

func NewChain(head error, tail ...error) Chain {
	return Chain{
		head: head,
		tail: tail,
	}
}

func (err Chain) Head() error {
	return err.head
}

func (err Chain) Error() string {
	var buf = &bytes.Buffer{}
	why.FprintFromIter(buf,
		err.head.Error(),
		len(err.tail),
		func(i int) (string, error) {
			return err.tail[i].Error(), nil
		})
	return buf.String()
}

func (err Chain) Unwrap() error {
	switch len(err.tail) {
	case 0:
		return err.head
	case 1:
		return Chain{
			head: err.tail[0],
		}
	default:
		return Chain{
			head: err.tail[0],
			tail: append([]error{}, err.tail[1:]...),
		}
	}
}
