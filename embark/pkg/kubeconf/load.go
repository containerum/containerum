// Package kubeconf provides adapter types for parsing kubectl config files
// and parsing functions
package kubeconf

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/ericchiang/k8s"
	"github.com/ghodss/yaml"
)

func Load(data []byte) (k8s.Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return k8s.Config{}, err
	}
	return config.ToK8S()
}

func LoadFromFile(filePath string) (k8s.Config, error) {
	var data, readConfigFileErr = ioutil.ReadFile(filePath)
	if readConfigFileErr != nil {
		return k8s.Config{}, readConfigFileErr
	}
	return Load(data)
}

func LoadFromReader(re io.Reader) (k8s.Config, error) {
	var buf = &bytes.Buffer{}
	if _, err := buf.ReadFrom(re); err != nil {
		return k8s.Config{}, err
	}
	return Load(buf.Bytes())
}
