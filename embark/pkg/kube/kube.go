package kube

import (
	"io"
	"time"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kubeClientAPI "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/genericclioptions/resource"
)

type Kube struct {
	Config kubeClientAPI.Config
	*kubernetes.Clientset
	config  *rest.Config
	factory cmdutil.Factory
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
	kube.factory = cmdutil.NewFactory(configFlags())
	kube.Clientset = kubecli
	kube.config = restConfig
	return &kube, nil
}

func configFlags() *genericclioptions.ConfigFlags {
	var config = genericclioptions.NewConfigFlags()
	return config
}

func (kube *Kube) builder(namespace string, reader io.Reader) *resource.Result {
	return kube.factory.NewBuilder().
		ContinueOnError().
		NamespaceParam(namespace).
		Stream(reader, "").
		DefaultNamespace().
		Flatten().
		Do()
}

func (kube *Kube) InstallFromReader(namespace string, re io.Reader) error {
	var infos, builderErr = kube.builder(namespace, re).Infos()
	if builderErr != nil {
		return builderErr
	}
	for _, info := range infos {
		kube.RESTClient().
			Post().
			Body(info.Object).
			Do()
		var obj, createErr = resource.
			NewHelper(kube.RESTClient(), info.Mapping).
			Create(namespace, true, info.Object)
		if createErr != nil {
			return createErr
		}
		if err := info.Refresh(obj, false); err != nil {
			return err
		}
	}
	return nil
}
