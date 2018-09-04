package builder

import (
	"time"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/logger"
)

const (
	nullDuration = time.Duration(0)
)

// builder Client options: logger, kube client config, etc.
type clientOptions struct {
	log            logger.Logger
	kubeConfigPath string
	configProvider kube.ConfigProvider
	timeout        time.Duration
}

func (options clientOptions) Ptr() *clientOptions {
	return &options
}

func DefaultClientOptionsPtr() *clientOptions {
	return &clientOptions{
		log:            logger.StdLogger(),
		configProvider: kube.AutoConfig,
		timeout:        60 * time.Second,
	}
}

func DefaultClientOptions() clientOptions {
	return *DefaultClientOptionsPtr()
}

func WithTimeout(timeout time.Duration) clientOptions {
	return clientOptions{
		timeout: timeout,
	}
}

func Debug() clientOptions {
	return clientOptions{
		log: logger.DebugLogger(),
	}
}

func KubeConfigPath(configpath string) clientOptions {
	return clientOptions{
		configProvider: kube.FileConfigProvider(configpath),
	}
}

func KubeConfigProvider(provider kube.ConfigProvider) clientOptions {
	return clientOptions{
		configProvider: provider,
	}
}

func WithLog(log logger.Logger) clientOptions {
	return clientOptions{
		log: log,
	}
}

func (options clientOptions) kubeConfig() (kube.Config, error) {
	var provider kube.ConfigProvider = kube.AutoConfig
	switch {
	case options.kubeConfigPath != "":
		provider = kube.FileConfigProvider(options.kubeConfigPath)
	case options.configProvider != nil:
		provider = options.configProvider
	}
	var kubeConfig, loadClientConfigErr = provider()
	if options.timeout != nullDuration {
		kubeConfig.Timeout = options.timeout
	}
	return kubeConfig, loadClientConfigErr
}

func (options *clientOptions) Merge(another ...clientOptions) *clientOptions {
	for _, anotherOptions := range another {
		if anotherOptions.log != nil {
			options.log = anotherOptions.log
		}

		if anotherOptions.kubeConfigPath != "" {
			options.kubeConfigPath = anotherOptions.kubeConfigPath
		}

		if anotherOptions.configProvider != nil {
			options.configProvider = anotherOptions.configProvider
		}

		if anotherOptions.timeout != nullDuration {
			options.timeout = anotherOptions.timeout
		}
	}
	return options
}
