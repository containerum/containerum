// Package fer provides formatting print functions to stderr
package fer

import (
	"fmt"
	"os"
	"sync"
)

var (
	mu sync.Mutex
)

func Println(args ...interface{}) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	return fmt.Fprintln(os.Stderr, args...)
}

func Printf(ff string, args ...interface{}) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	return fmt.Fprintf(os.Stdout, ff, args...)
}

func Print(args ...interface{}) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	return fmt.Fprint(os.Stdout, args...)
}

func Fatal(ff string, args ...interface{}) {
	Printf(ff, args...)
	os.Exit(1)
}
