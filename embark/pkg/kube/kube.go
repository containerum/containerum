package kube

import (
	"context"
	"time"

	weirdKubeClient "github.com/ericchiang/k8s"
)

type Kube struct {
	timeout time.Duration
	*weirdKubeClient.Client
}

func NewKube(options ..._Config) (Kube, error) {
	var config = _Config{
		timeout:               60 * time.Second,
		kubectlConfigProvider: StdKubectConfig,
	}.Merge(options...)

	var kubectlConfig, kubectlConfigErr = config.KubectlConfig()
	if kubectlConfigErr != nil {
		return Kube{}, kubectlConfigErr
	}
	var weirdClient, newWeirdClientErr = weirdKubeClient.NewClient(kubectlConfig.beWeird())
	if newWeirdClientErr != nil {
		return Kube{}, newWeirdClientErr
	}
	var kubeClinent = Kube{
		timeout: config.timeout,
		Client:  weirdClient,
	}
	return kubeClinent, nil
}

func (kube Kube) Create(obj Object) error {
	RegisterObject()
	var ctx, done = context.WithTimeout(context.Background(), kube.timeout)
	defer done()
	return kube.Client.Create(ctx, &obj)
}
