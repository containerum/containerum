package cgraph

import "fmt"

type Node struct {
	Name   string
	Action func() error
	Deps   []string

	executed bool
}

func (node *Node) IsExecuted() bool {
	return node.executed
}

func (node *Node) vals() (name string, deps []string, action func() error) {
	return node.Name, append([]string{}, node.Deps...), node.Action
}

func (node *Node) Copy() *Node {
	return &Node{
		Name:   node.Name,
		Action: node.Action,
		Deps:   append([]string{}, node.Deps...),
	}
}

func (node *Node) Execute() error {
	if !node.executed {
		defer func() {
			//	log.Printf("node %q executed", node.Name)
			node.executed = true
		}()
		return node.Action()
	}
	return nil
}

type Graph map[string]*Node

func NewGraph() Graph {
	return make(Graph)
}

func NewGraphPrealloc(n int) Graph {
	return make(Graph, n)
}

/*
type GraphBuilder = func(node *Node) ([]string, error)

func BuildGraph(start Node, builder GraphBuilder) (Graph, error) {
	var graph = NewGraph()
	graph.AddNode(start.vals())
	var q = make(queue, 16)
	var stop = make(chan struct{})
	defer close(stop)
	q.Push(stop, start.Name)
	q.Push(stop, start.Deps...)
	for {
		var nodeName, ok = <-q
		if !ok {
			break
		}
		var node = graph[nodeName]
		if node != nil {
			continue
		}
		node = &Node{
			Name: nodeName,
		}
		graph[nodeName] = node
		var next, err = builder(node)
		if err != nil {
			return nil, err
		}
		q.Push(stop, next...)
	}
	return graph, nil
}
*/

func (graph Graph) GetNode(name string) *Node {
	var node, ok = graph[name]
	if !ok {
		panic(fmt.Sprintf("[cgraph.NGraph.GetNode] node %q is not defined", name))
	}
	return node
}

func (graph Graph) Execute(names ...string) error {
	for _, name := range names {
		var node = graph.GetNode(name)
		for _, depName := range node.Deps {
			if err := graph.Execute(depName); err != nil {
				return err
			}
		}
		if err := node.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func (graph Graph) AddNode(name string, deps []string, action func() error) Graph {
	graph[name] = &Node{
		Name:   name,
		Deps:   append([]string{}, deps...),
		Action: action,
	}
	return graph
}

func (graph Graph) Nodes() []string {
	var nodes = make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}
	return nodes
}

type queue chan string

func (q queue) Push(stop <-chan struct{}, elems ...string) {
	go func() {
		for _, elem := range elems {
			select {
			case <-stop:
				return
			case q <- elem:
				continue
			}
		}
	}()
}
