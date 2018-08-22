package cgraph

import "fmt"

type Future chan error

func (future Future) Await() error {
	return <-future
}

func (future Future) Ok() {
	close(future)
}

func (future Future) Err(err error) {
	future <- err
}

type Node struct {
	Name   string
	Deps   []Node
	Action func() error
}

type Graph struct {
	nodes map[string]Node
}

func (graph Graph) AddNode(name string, action func() error, dependsOn ...string) {
	var deps = make([]Node, 0, len(dependsOn))
	for _, depName := range dependsOn {
		var dep, exists = graph.nodes[depName]
		if !exists {
			panic(fmt.Sprintf("[cgrpah.Graph.AddNode] dependency %q is not defined", depName))
		}
		deps = append(deps, dep)
	}
	graph.nodes[name] = Node{
		Name:   name,
		Action: action,
		Deps:   deps,
	}
}

type Nodes []Node

func (nodes Nodes) Sources() []Node {
	var sources = make([]Node, 0)
	for _, node := range nodes {
		if len(node.Deps) == 0 {
			sources = append(sources, node)
		}
	}
	return sources
}

func (nodes Nodes) Sinks() []Node {
	var deps = make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		for _, dep := range node.Deps {
			deps[dep.Name] = struct{}{}
		}
	}
	var sinks = make(Nodes, 0, len(nodes))
	for _, node := range nodes {
		var _, isDep = deps[node.Name]
		if !isDep {
			sinks = append(sinks, node)
		}
	}
	return sinks
}
