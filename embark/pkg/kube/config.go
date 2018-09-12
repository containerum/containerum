package kube

import "time"

type _Config struct {
	timeout               time.Duration
	kubectlConfigProvider KubectlConfigProvider
	kubectlConfig         *KubectlConfig
}

func (config _Config) Merge(anothers ..._Config) _Config {
	for _, another := range anothers {
		if another.kubectlConfig != nil {
			config.kubectlConfig = another.kubectlConfig
		}
		if another.kubectlConfigProvider != nil {
			config.kubectlConfigProvider = another.kubectlConfigProvider
		}
		if another.timeout != 0 {
			config.timeout = another.timeout
		}
	}
	return config
}

func (config _Config) KubectlConfig() (KubectlConfig, error) {
	if config.kubectlConfig != nil {
		return (*config.kubectlConfig).Copy(), nil
	}
	if config.kubectlConfigProvider != nil {
		return config.kubectlConfigProvider()
	}
	return StdKubectConfig()
}
