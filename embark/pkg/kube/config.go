package kube

import (
	"time"

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
