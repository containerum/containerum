package kube

import (
	"context"
	"time"

	weirdKubeClient "github.com/ericchiang/k8s"
	weirdMetaV1 "github.com/ericchiang/k8s/apis/meta/v1"
)

type Object interface {
	Kind() string
	weirdKubeClient.Resource
}

var (
	_ Object = ObjectMock{}
)

type ObjectMock struct {
	ObjectKind string
	Meta       weirdMetaV1.ObjectMeta
}

func (mock ObjectMock) Kind() string {
	return mock.ObjectKind
}

func (mock ObjectMock) GetMetadata() *weirdMetaV1.ObjectMeta {
	return &mock.Meta
}

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
	var ctx, done = context.WithTimeout(context.Background(), kube.timeout)
	defer done()
	return kube.Client.Create(ctx, obj)
}
