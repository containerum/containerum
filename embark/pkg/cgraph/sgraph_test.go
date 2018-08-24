package cgraph

import "testing"

func TestSGraph_Walk(test *testing.T) {
	var graph = SGraph{
		"A": []string{"B", "C"},
		"B": []string{"D"},
		"D": []string{"C"},
	}
	graph.Walk("A", func(node string) {
		test.Log(node)
	})
}
