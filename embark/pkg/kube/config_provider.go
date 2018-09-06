package kube

import (
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

type ConfigProvider func() (Config, error)

var (
	_ ConfigProvider = AutoConfig
	_                = FileConfigProvider("")
	_                = ConfigFromBytes(nil)
	_ ConfigProvider = func() (Config, error) { return LoadConfigFromFile("") }
)

// Looks for the configuration file in the "~/.kube/config", if it is defined by KUBECONFIG env, then it tries to load from there.
func AutoConfig() (Config, error) {
	var configPath = os.ExpandEnv("$HOME/.kube/config")
	{
		var envConfigPath, envConfigPathDefined = os.LookupEnv("KUBECONFIG")
		if envConfigPathDefined {
			configPath = envConfigPath
		}
	}
	return LoadConfigFromFile(configPath)
}

// Creates config provider, which will load data from provided file
func FileConfigProvider(filename string) ConfigProvider {
	return func() (Config, error) {
		return LoadConfigFromFile(filename)
	}
}

func ConfigFromBytes(data []byte) ConfigProvider {
	return func() (Config, error) {
		var config Config
		var kc, loadKubeConfigErr = clientcmd.Load(data)
		if loadKubeConfigErr != nil {
			return config, loadKubeConfigErr
		}
		config.Config = *kc
		return config, nil
	}
}

// Loads kube config from provided file
func LoadConfigFromFile(filename string) (Config, error) {
	var config Config
	var kc, loadKubeConfigErr = clientcmd.LoadFromFile(filename)
	if loadKubeConfigErr != nil {
		return config, loadKubeConfigErr
	}
	config.Config = *kc
	return config, nil
}
