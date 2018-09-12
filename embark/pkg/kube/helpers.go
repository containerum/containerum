package kube

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/containerum/containerum/embark/pkg/utils/ym"
)

var (
	_ KubectlConfigProvider = StdKubectConfig
)

func StdKubectConfig() (KubectlConfig, error) {
	var config KubectlConfig
	return config, ym.LoadYAML(autoKubectlConfigPath(), &config)
}

func LoadKubectlConfigFromPath(configPath string) KubectlConfigProvider {
	return func() (KubectlConfig, error) {
		var config KubectlConfig
		return config, ym.LoadYAML(configPath, &config)
	}
}

func KubectlConfigFromReader(re io.Reader) KubectlConfigProvider {
	var config KubectlConfig
	var err = yaml.NewDecoder(re).Decode(&config)
	return config.AsProviderWithErr(err)
}

func autoKubectlConfigPath() string {
	var configPathFromEnv, configPathFromEnvExists = os.LookupEnv("KUBECONFIG")
	if configPathFromEnvExists {
		return configPathFromEnv
	}
	return "~/.kube/config"
}
