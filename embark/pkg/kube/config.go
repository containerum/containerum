package kube

import (
	"io/ioutil"
	"os"

	"time"

	"gopkg.in/yaml.v2"
	kubeClientAPI "k8s.io/client-go/tools/clientcmd/api"
)

// Kube Client config
type Config struct {
	Timeout time.Duration `json:"timeout. omitempty"`
	kubeClientAPI.Config
}

// Extracts pointer to kube client clientcmd.Config
func (config Config) ToKube() *kubeClientAPI.Config {
	return &config.Config
}

// Helper method for clientcmd.BuildConfigFromKubeconfigGetter
// ```go
// var restConfig, err = clientcmd.BuildConfigFromKubeconfigGetter("", config.Getter())
// ```
func (config Config) Getter() func() (*kubeClientAPI.Config, error) {
	return func() (*kubeClientAPI.Config, error) {
		return config.ToKube(), nil
	}
}

type ConfigProvider func() (Config, error)

var (
	_ ConfigProvider = AutoConfig
)

// Looks for the configuration file in the "~/.kube/config", if it is defined by KUBECONFIG env, then it tries to load from there.
func AutoConfig() (Config, error) {
	var configPath = "~/.kube/config"
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

// Loads kube config from provided file
func LoadConfigFromFile(filename string) (Config, error) {
	var data, readFileErr = ioutil.ReadFile(filename)
	if readFileErr != nil {
		return Config{}, readFileErr
	}
	var config Config
	return config, yaml.Unmarshal(data, &config)
}
