package cgraph

import (
	"fmt"
)

type Node struct {
	name   string
	action func() error
	deps   []string

	executed bool
}

func (node *Node) Execute() error {
	if !node.executed {
		defer func() {
			//	log.Printf("node %q executed", node.name)
			node.executed = true
		}()
		return node.action()
	}
	return nil
}

type Graph map[string]*Node

func (graph Graph) GetNode(name string) *Node {
	var node, ok = graph[name]
	if !ok {
		panic(fmt.Sprintf("[cgraph.NGraph.GetNode] node %q is not defined", name))
	}
	return node
}

func (graph Graph) Execute(name string) error {
	var node = graph.GetNode(name)
	for _, depName := range node.deps {
		if err := graph.Execute(depName); err != nil {
			return err
		}
	}
	return node.Execute()
}

func (graph Graph) AddNode(name string, deps []string, action func() error) Graph {
	graph[name] = &Node{
		name:   name,
		deps:   append([]string{}, deps...),
		action: action,
	}
	return graph
}
