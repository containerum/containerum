package containerum

import (
	"net/url"
	"path"
	"strings"

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

func (component ComponentWithName) URL() string {
	var repo = strings.TrimPrefix(component.Repo, "http://")
	repo = strings.TrimPrefix(repo, "https://")
	var componentPath = component.Name + "-" + component.Version + ".tgz"
	return "http://" + repo + "/charts/" + url.PathEscape(componentPath)
}

type Component struct {
	Version   string                 `json:"version"`
	Repo      string                 `json:"repo"`
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

func (component Component) ObjectNames() []string {
	var names = make([]string, 0, len(component.Objects))
	for _, object := range component.Objects {
		var nameWithExt = path.Base(object)
		var ext = path.Ext(nameWithExt)
		var name = strings.TrimSuffix(nameWithExt, ext)
		names = append(names, name)
	}
	return names
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
