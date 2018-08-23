package cgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph_Execute(test *testing.T) {
	for i := 0; i < 100 && !test.Failed(); i++ {
		var graph = make(Graph)
		var history = []string{}
		graph.AddNode("A", []string{"B", "C"}, func() error {
			history = append(history, "A")
			return nil
		})
		graph.AddNode("B", []string{"C"}, func() error {
			history = append(history, "B")
			return nil
		})
		graph.AddNode("C", []string{}, func() error {
			history = append(history, "C")
			return nil
		})
		graph.Execute("A")
		assert.Equal(test, []string{"C", "B", "A"}, history)
	}
}
