package kube

import (
	"io"
	"os"

	"github.com/containerum/containerum/embark/pkg/utils/kubeconf"
)

var (
	_ KubectlConfigProvider = StdKubectConfig
)

func StdKubectConfig() (KubectlConfig, error) {
	var config, err = kubeconf.LoadFromFile(autoKubectlConfigPath())
	return KubectlConfig(config), err
}

func LoadKubectlConfigFromPath(configPath string) KubectlConfigProvider {
	return func() (KubectlConfig, error) {
		var config, err = kubeconf.LoadFromFile(configPath)
		if err != nil {
			return KubectlConfig{}, err
		}
		return KubectlConfig(config), err
	}
}

func KubectlConfigFromReader(re io.Reader) KubectlConfigProvider {
	var config, err = kubeconf.LoadFromReader(re)
	return KubectlConfig(config).AsProviderWithErr(err)
}

func DecodeConfig(data []byte) (KubectlConfig, error) {
	var config, err = kubeconf.Load(data)
	return KubectlConfig(config), err
}

func autoKubectlConfigPath() string {
	var configPathFromEnv, configPathFromEnvExists = os.LookupEnv("KUBECONFIG")
	if configPathFromEnvExists {
		return configPathFromEnv
	}
	return os.ExpandEnv("$HOME/.kube/config")
}
