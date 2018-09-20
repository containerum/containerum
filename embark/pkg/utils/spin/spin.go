package spin

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Loop struct {
	Frames  []string
	Prefix  string
	pointer int
	prev    string
}

func (loop *Loop) Next() string {
	loop.Erase()
	var frame = loop.Prefix + loop.Frames[loop.pointer]
	var atStart = strings.Repeat("\b", len(frame))
	loop.prev = frame
	var diff = frame + atStart
	loop.pointer++
	if loop.pointer >= len(loop.Frames) {
		loop.pointer = 0
	}
	return diff
}

func (loop Loop) Erase() {
	var back = strings.Repeat("\b", len(loop.prev))
	var erase = strings.Repeat(" ", utf8.RuneCountInString(loop.prev)+1)
	if len(loop.prev) > 0 {
		back += "\b"
	}
	fmt.Print(back + erase + back)
}
