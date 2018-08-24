package builder

import (
	"io/ioutil"

	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"github.com/go-yaml/yaml"
	"os"
	"unicode/utf8"
)

func LoadYAML(path string, to interface{}) error {
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, to)
}

func MkDirIfNotExists(dir string) error {
	var err = os.MkdirAll(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func RequirementsNameWidth(req requirements.Requirements) int {
	var width = 0
	for _, dep := range req.Dependencies {
		var w = utf8.RuneCountInString(dep.Name)
		if w > width {
			width = w
		}
	}
	return width
}
