package cgraph

type SGraph map[string][]string

func (graph SGraph) Walk(start string, walker func(node string)) {
	walker(start)
	for _, child := range graph[start] {
		graph.Walk(child, walker)
	}
}
