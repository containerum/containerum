package installer

import (
	"path"
	"strings"

	"github.com/containerum/containerum/embark/pkg/models/components"
	"github.com/containerum/containerum/embark/pkg/static"

	"gopkg.in/yaml.v2"
)

func loadContainerumConfig(configpath string) (components.Components, error) {
	var containerumConfig components.Components
	var data []byte
	var loadContDataErr error
	switch strings.TrimSpace(configpath) {
	case "", "static":
		data, loadContDataErr = static.ReadFile("containerum.yaml")
	default:
		data, loadContDataErr = static.ReadFile(path.Clean(configpath))
	}
	if loadContDataErr != nil {
		return containerumConfig, loadContDataErr
	}
	if err := yaml.Unmarshal(data, &containerumConfig); err != nil {
		return containerumConfig, err
	}
	return containerumConfig, nil
}
