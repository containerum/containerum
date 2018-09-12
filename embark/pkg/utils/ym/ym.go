package ym

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadYAML(path string, to interface{}) error {
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, to)
}
