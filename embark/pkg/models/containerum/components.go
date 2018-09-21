package containerum

import (
	"net/url"
	"path"
	"strings"
)

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
	return "http://" + repo + "/charts/" + url.PathEscape(component.Name+"-"+component.Version+".tgz")
}

type Component struct {
	Version   string                 `json:"version"`
	Repo      string                 `json:"repo"`
	Objects   []string               `json:"objects"`
	DependsOn []string               `json:"depends_on"`
	Values    map[string]interface{} `json:"values"`
}

func (component Component) WithValues(mixins ...map[string]interface{}) map[string]interface{} {
	var values = copyTree(component.Values)
	for _, mixin := range mixins {
		for k, v := range mixin {
			values[k] = v
		}
	}
	return values
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
