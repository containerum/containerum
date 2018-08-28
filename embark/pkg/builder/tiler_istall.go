package builder

import (
	"github.com/go-yaml/yaml"
	kubeAppsV1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube struct {
	*kubernetes.Clientset
	config *rest.Config
}

// goto init()
var defaulTiler = kubeAppsV1.Deployment{}

func NewKubeClient(cfgpath string) (*Kube, error) {
	var kube Kube
	var config *rest.Config
	var err error

	if cfgpath == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", cfgpath)
	}
	if err != nil {
		return nil, err
	}
	kubecli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kube.Clientset = kubecli
	kube.config = config
	return &kube, nil
}

func (client *Client) InstallTiller(kubeConfigFile string) error {
	var kube, err = NewKubeClient(kubeConfigFile)
	if err != nil {
		return err
	}
	var _, createTilerErr = kube.AppsV1().
		Deployments(defaulTiler.Namespace).
		Create(&defaulTiler)
	return createTilerErr
}

func init() {
	const paragonTiller = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: helm
    name: tiller
  name: tiller-deploy
  namespace: kube-system
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: helm
        name: tiller
    spec:
      containers:
      - env:
        - name: TILLER_NAMESPACE
          value: kube-system
        - name: TILLER_HISTORY_MAX
          value: "0"
        image: gcr.io/kubernetes-helm/tiller:v2.9.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 44135
          initialDelaySeconds: 1
          timeoutSeconds: 1
        name: tiller
        ports:
        - containerPort: 44134
          name: tiller
        - containerPort: 44135
          name: http
        readinessProbe:
          httpGet:
            path: /readiness
            port: 44135
          initialDelaySeconds: 1
          timeoutSeconds: 1
`
	if err := yaml.Unmarshal([]byte(paragonTiller), &defaulTiler); err != nil {
		panic(err)
	}
}

func newInt32(i int32) *int32 {
	return &i
}
