package containerum

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCopyTree(test *testing.T) {
	var x = map[string]interface{}{
		"1": 1,
		"object": map[string]interface{}{
			"ints":       []int{1, 2, 3},
			"strings":    []string{"foo", "bar"},
			"interfaces": []interface{}{1, 2.0, "three"},
		},
		"string": "string",
		"component": Component{
			Version:   "42",
			Objects:   []string{"UFO"},
			DependsOn: []string{"nothing"},
			Values: map[string]interface{}{
				"data": "ball",
			},
		},
	}
	var cp = copyTree(x)
	assert.Equal(test, cp, x)
}
