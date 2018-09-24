package components

import (
	"gopkg.in/yaml.v2"
)

type Components map[string]Component

func (comps Components) Len() int {
	return len(comps)
}

func (comps Components) Copy() Components {
	var cp = make(Components, comps.Len())
	for name, component := range comps {
		cp[name] = component.Copy()
	}
	return cp
}

func (comps Components) Slice() []ComponentWithName {
	var components = make([]ComponentWithName, 0, len(comps))
	for name, component := range comps {
		components = append(components, ComponentWithName{
			Name:      name,
			Component: component.Copy(),
		})
	}
	return components
}

func (comps Components) String() string {
	var data, _ = yaml.Marshal(comps)
	return string(data)
}

func (comps Components) New() Components {
	return make(Components, comps.Len())
}

func (comps Components) Filter(pred func(component ComponentWithName) bool) Components {
	var filtered = comps.New()
	for name, component := range comps {
		if pred(ComponentWithName{
			Name:      name,
			Component: component.Copy(),
		}) {
			filtered[name] = component.Copy()
		}
	}
	return filtered
}

func (comps Components) Contains(name string) bool {
	var _, contains = comps[name]
	return contains
}

func (comps Components) Get(name string) (ComponentWithName, bool) {
	var component, ok = comps[name]
	if !ok {
		return ComponentWithName{}, false
	}
	return ComponentWithName{
		Name:      name,
		Component: component.Copy(),
	}, true
}

func (comps Components) MustGet(name string) ComponentWithName {
	var component, ok = comps[name]
	if !ok {
		return ComponentWithName{}
	}
	return ComponentWithName{
		Name:      name,
		Component: component.Copy(),
	}
}

func (comps Components) Names() []string {
	var names = make([]string, 0, comps.Len())
	for name := range comps {
		names = append(names, name)
	}
	return names
}
