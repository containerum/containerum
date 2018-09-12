package kube

import "time"

func WithTimeout(timeout time.Duration) _Config {
	return _Config{
		timeout: timeout,
	}
}

func WithKubectlConfig(kubectlConfig KubectlConfig) _Config {
	return _Config{
		kubectlConfigProvider: kubectlConfig.Provider,
	}
}

func WithKubectlConfigProvider(provider KubectlConfigProvider) _Config {
	return _Config{
		kubectlConfigProvider: provider,
	}
}
