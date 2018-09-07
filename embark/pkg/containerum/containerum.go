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

type ComponentWithName struct {
	Name string `json:"name"`
	Component
}

func (component ComponentWithName) Copy() ComponentWithName {
	return ComponentWithName{
		Name:      component.Name,
		Component: component.Component.Copy(),
	}
}

type Component struct {
	Version   string                 `json:"version"`
	Objects   []string               `json:"objects"`
	DependsOn []string               `json:"depends_on"`
	Values    map[string]interface{} `json:"values"`
}

func (component Component) Copy() Component {
	var cp = component
	cp.Objects = append([]string{}, cp.Objects...)
	cp.DependsOn = append([]string{}, cp.DependsOn...)
	component.Values = copyTree(component.Values)
	return component
}

func copyTree(tree map[string]interface{}) map[string]interface{} {
	var cp = make(map[string]interface{})
	for k, v := range tree {
		switch v := v.(type) {
		case nil:
			continue
		case map[string]interface{}:
			cp[k] = copyTree(v)
		case []string:
			cp[k] = append([]string{}, v...)
		case []int:
			cp[k] = append([]int{}, v...)
		case []interface{}:
			cp[k] = append([]interface{}{}, v...)
		case Component:
			cp[k] = v.Copy()
		default:
			cp[k] = v
		}
	}
	return cp
}
