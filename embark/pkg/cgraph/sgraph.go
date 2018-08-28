package cgraph

type Walker = func(node string, path []string, children []string)

type SGraph map[string][]string

func (graph SGraph) Walk(start string, walker Walker) {
	graph.walk(start, []string{}, walker)
}

func (graph SGraph) walk(node string, path []string, walker Walker) {
	var children = graph.ChildrenOf(node)
	walker(node, copyStrSlice(path), copyStrSlice(children))
	for _, child := range children {
		var nodePath = append(copyStrSlice(path), node)
		graph.walk(child, nodePath, walker)
	}
}

func (graph SGraph) ChildrenOf(name string) []string {
	return copyStrSlice(graph[name])
}

func (graph SGraph) Nodes() []string {
	var nodes = make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}
	return nodes
}

func (graph SGraph) AddNode(name string, children ...string) {
	graph[name] = children
}

func copyStrSlice(s []string) []string {
	return append([]string{}, s...)
}
