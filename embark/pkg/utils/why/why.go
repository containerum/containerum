package why

import (
	"fmt"
	"io"
	"os"
)

func Fprint(wr io.Writer, head string, items ...string) error {
	var _, err = fmt.Fprintf(wr, "%s\n", head)
	if err != nil {
		return err
	}
	var itemsLastIndex = len(items) - 1
	for i, item := range items {
		var prefix = "╠"
		if i == itemsLastIndex {
			prefix = "╚"
		}
		var _, err = fmt.Fprintf(wr, "\t%s═ %s\n", prefix, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func Print(head string, items ...string) error {
	return Fprint(os.Stdout, head, items...)
}

func FprintFromIter(wr io.Writer, head string, n int, source func(i int) (string, error)) error {
	var _, err = fmt.Fprintf(wr, "%s\n", head)
	if err != nil {
		return err
	}
	var itemsLastIndex = n - 1
	for i := 0; i < n; i++ {
		var prefix = "╠"
		if i == itemsLastIndex {
			prefix = "╚"
		}
		var item, getItemErr = source(i)
		if getItemErr != nil {
			return getItemErr
		}
		var _, err = fmt.Fprintf(wr, "\t%s═ %s\n", prefix, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintFromIter(head string, n int, source func(i int) (string, error)) error {
	return FprintFromIter(os.Stdout, head, n, source)
}
