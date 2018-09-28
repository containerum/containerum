package depsearch

import (
	"testing"
)

func TestStatic(test *testing.T) {
	var searcher = Static()
	for _, componentName := range []string{"auth"} {
		if !searcher.Contains(componentName) {
			test.Fatalf("unable to find component %q in static FS", componentName)
		}
	}
}
