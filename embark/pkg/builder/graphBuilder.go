package builder

import (
	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/containerum"
)

func BuildGraph(baseDir string, cont containerum.Containerum) (cgraph.Graph, error) {
	var dependencyGraph = cgraph.NewGraph()
	for _, component := range cont.Components() {
		component = component.Copy()
		dependencyGraph.AddNode(component.Name, component.DependsOn, func() error {
			return nil
		})
	}
	return dependencyGraph, nil
}
