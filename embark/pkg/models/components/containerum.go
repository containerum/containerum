package components

import (
	"gopkg.in/yaml.v2"
)

type Components map[string]Component

func (containerum Components) Len() int {
	return len(containerum)
}

func (containerum Components) Copy() Components {
	var cp = make(Components, containerum.Len())
	for name, component := range containerum {
		cp[name] = component.Copy()
	}
	return cp
}

func (containerum Components) Components() []ComponentWithName {
	var components = make([]ComponentWithName, 0, len(containerum))
	for name, component := range containerum {
		components = append(components, ComponentWithName{
			Name:      name,
			Component: component.Copy(),
		})
	}
	return components
}

func (containerum Components) String() string {
	var data, _ = yaml.Marshal(containerum)
	return string(data)
}

func (containerum Components) New() Components {
	return make(Components, containerum.Len())
}

func (containerum Components) Filter(pred func(component ComponentWithName) bool) Components {
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
