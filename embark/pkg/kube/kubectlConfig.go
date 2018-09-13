package kube

import (
	"github.com/containerum/containerum/embark/pkg/utils/kubeconf"
	weirdKubeClient "github.com/ericchiang/k8s"
	goyaml "github.com/go-yaml/yaml"
)

type KubectlConfigProvider func() (KubectlConfig, error)

type KubectlConfig weirdKubeClient.Config

func (config KubectlConfig) Copy() KubectlConfig {
	return KubectlConfig{
		Kind:           config.Kind,
		APIVersion:     config.APIVersion,
		Preferences:    config.Preferences,
		Clusters:       append([]weirdKubeClient.NamedCluster{}, config.Clusters...),
		AuthInfos:      append([]weirdKubeClient.NamedAuthInfo{}, config.AuthInfos...),
		Contexts:       append([]weirdKubeClient.NamedContext{}, config.Contexts...),
		CurrentContext: config.CurrentContext,
		Extensions:     append([]weirdKubeClient.NamedExtension{}, config.Extensions...),
	}
}

func (config KubectlConfig) beWeird() *weirdKubeClient.Config {
	var weirdConfig = weirdKubeClient.Config(config.Copy())
	return &weirdConfig
}

var (
	_ KubectlConfigProvider = KubectlConfig{}.Provider
	_                       = KubectlConfig{}.AsProviderWithErr(nil)
)

func (config KubectlConfig) Provider() (KubectlConfig, error) {
	return config.Copy(), nil
}

func (config KubectlConfig) AsProviderWithErr(err error) KubectlConfigProvider {
	if err != nil {
		return func() (KubectlConfig, error) {
			return KubectlConfig{}, err
		}
	}
	return config.Provider
}

var (
	_ goyaml.Unmarshaler = new(KubectlConfig)
)

func (config *KubectlConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var adapter kubeconf.Config
	if err := unmarshal(&adapter); err != nil {
		return err
	}
	var convertedConfig, err = adapter.ToK8S()
	if err != nil {
		return err
	}
	*config = KubectlConfig(convertedConfig)
	return nil
}
