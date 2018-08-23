package requirements

import (
	"fmt"
	"net/url"
)

type Requirements struct {
	Dependencies []Dependency
}

type Dependency struct {
	Name       string
	Repository string
	Version    string
	Tags       []string
}

func (dep Dependency) URL() (string, error) {
	var addr, err = url.Parse(dep.Repository)
	if err != nil {
		return "", fmt.Errorf("invalid repository addr of dependency %q: %v", dep.Name, err)
	}
	addr.Path = dep.FileName()
	return addr.String(), nil
}

func (dep Dependency) FileName() string {
	var id = fmt.Sprintf("%s-%s.tgz", dep.Name, dep.Version)
	return url.PathEscape(id)
}
