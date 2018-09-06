package kube

import (
	"time"

	"github.com/containerum/containerum/embark/pkg/emberr"
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
	var restConfig, buildingKubeConfigFromGetterErr = clientcmd.BuildConfigFromKubeconfigGetter("", config.Getter())
	if buildingKubeConfigFromGetterErr != nil {
		return nil, emberr.ErrUnableToCreateKubeCLient{Reason: buildingKubeConfigFromGetterErr, Comment: "while building REST config from getter"}
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	restConfig.Timeout = config.Timeout
	var kubecli, newClientsetErr = kubernetes.NewForConfig(restConfig)
	if newClientsetErr != nil {
		return nil, emberr.ErrUnableToCreateKubeCLient{Reason: newClientsetErr, Comment: "while creating clientset from REST config"}
	}
	kube.Clientset = kubecli
	kube.config = restConfig
	return &kube, nil
}
