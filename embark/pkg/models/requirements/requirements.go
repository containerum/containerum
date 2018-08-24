package requirements

import (
	"fmt"
	"net/url"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Requirements struct {
	Dependencies []Dependency
}

func FromChart(ch *chart.Chart) (Requirements, error) {
	var hreq, err = chartutil.LoadRequirements(ch)
	if err != nil {
		return Requirements{}, err
	}
	var req = Requirements{
		Dependencies: make([]Dependency, 0, len(hreq.Dependencies)),
	}
	for _, dep := range hreq.Dependencies {
		req.Dependencies = append(req.Dependencies, Dependency(*dep))
	}
	return req, nil
}

func (requirements Requirements) Names() []string {
	var names = make([]string, 0, len(requirements.Dependencies))
	for _, dep := range requirements.Dependencies {
		names = append(names, dep.Name)
	}
	return names
}

func (requirements Requirements) DependencySet() DependencySet {
	var set = make(DependencySet, len(requirements.Dependencies))
	for _, dep := range requirements.Dependencies {
		set[dep.Name] = dep
	}
	return set
}

type Dependency chartutil.Dependency

func (dep Dependency) URL() (string, error) {
	var addr, err = url.Parse(dep.Repository)
	if err != nil {
		return "", fmt.Errorf("invalid repository addr of dependency %q: %v", dep.Name, err)
	}
	addr.Path = "/charts/" + dep.FileName()
	return addr.String(), nil
}

func (dep Dependency) FileName() string {
	var id = fmt.Sprintf("%s-%s.tgz", dep.Name, dep.Version)
	return url.PathEscape(id)
}

func (dep Dependency) String() string {
	return fmt.Sprintf("%s-%s", dep.Name, dep.Version)
}

type DependencySet map[string]Dependency

func (dependency DependencySet) Add(deps ...Dependency) DependencySet {
	for _, dep := range deps {
		dependency[dep.Name] = dep
	}
	return dependency
}
