package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kubeClientAPI "k8s.io/client-go/tools/clientcmd/api"
)

type Kube struct {
	Config kubeClientAPI.Config
	*kubernetes.Clientset
	config *rest.Config
}

func NewKubeClient(config Config) (*Kube, error) {
	var kube Kube
	var restConfig, err = clientcmd.BuildConfigFromKubeconfigGetter("", config.Getter())
	if err != nil {
		return nil, err
	}
	restConfig.Timeout = config.Timeout
	kubecli, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	kube.Clientset = kubecli
	kube.config = restConfig
	return &kube, nil
}
