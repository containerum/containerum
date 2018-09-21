package containerum

import (
	"gopkg.in/yaml.v2"
)

type Containerum map[string]Component

func (containerum Containerum) Len() int {
	return len(containerum)
}

func (containerum Containerum) Copy() Containerum {
	var cp = make(Containerum, containerum.Len())
	for name, component := range containerum {
		cp[name] = component.Copy()
	}
	return cp
}

func (containerum Containerum) Components() []ComponentWithName {
	var components = make([]ComponentWithName, 0, len(containerum))
	for name, component := range containerum {
		components = append(components, ComponentWithName{
			Name:      name,
			Component: component.Copy(),
		})
	}
	return components
}

func (containerum Containerum) String() string {
	var data, _ = yaml.Marshal(containerum)
	return string(data)
}

func (containerum Containerum) New() Containerum {
	return make(Containerum, containerum.Len())
}

func (containerum Containerum) Filter(pred func(component ComponentWithName) bool) Containerum {
	var filtered = containerum.New()
	for name, component := range containerum {
		if pred(ComponentWithName{
			Name:      name,
			Component: component.Copy(),
		}) {
			filtered[name] = component.Copy()
		}
	}
	return filtered
}
