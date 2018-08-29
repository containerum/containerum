package builder

import (
	kubeAppsV1 "k8s.io/api/apps/v1"
	kubeCoreV1 "k8s.io/api/core/v1"
	kubeAPIv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

// goto init()
var defaulTiler = kubeAppsV1.Deployment{
	TypeMeta: kubeAPIv1.TypeMeta{
		Kind:       "deployment",
		APIVersion: "extensions/v1beta1",
	},

	ObjectMeta: kubeAPIv1.ObjectMeta{
		Name:      "tiller-deploy",
		Namespace: "kube-system",
		Labels: map[string]string{
			"app":  "helm",
			"name": "tiller",
		},
	},
	Spec: kubeAppsV1.DeploymentSpec{
		Selector: &kubeAPIv1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "helm",
			},
		},
		Replicas: newInt32(1),
		Template: kubeCoreV1.PodTemplateSpec{
			ObjectMeta: kubeAPIv1.ObjectMeta{
				Labels: map[string]string{
					"app":  "helm",
					"name": "tiller",
				},
			},
			Spec: kubeCoreV1.PodSpec{
				Containers: []kubeCoreV1.Container{
					{
						Image:           "gcr.io/kubernetes-helm/tiller:v2.9.1",
						Name:            "tiller",
						ImagePullPolicy: "IfNotPresent",
						LivenessProbe: &kubeCoreV1.Probe{
							InitialDelaySeconds: 1,
							TimeoutSeconds:      1,
							Handler: kubeCoreV1.Handler{
								HTTPGet: &kubeCoreV1.HTTPGetAction{
									Path: "/liveness",
									Port: intstr.FromInt(44135),
								},
							},
						},
						Env: []kubeCoreV1.EnvVar{
							{Name: "TILLER_NAMESPACE", Value: ""},
							{Name: "TILLER_HISTORY_MAX", Value: "0"},
						},
						Ports: []kubeCoreV1.ContainerPort{
							{ContainerPort: 44134, Name: "tiller"},
							{ContainerPort: 44135, Name: "http"},
						},
						ReadinessProbe: &kubeCoreV1.Probe{
							InitialDelaySeconds: 1,
							TimeoutSeconds:      1,
							Handler: kubeCoreV1.Handler{
								HTTPGet: &kubeCoreV1.HTTPGetAction{
									Path: "/readiness",
									Port: intstr.FromInt(44135),
								},
							},
						},
					},
				},
			},
		},
	},
}

func NewKubeClient(config kubeClientAPI.Config) (*Kube, error) {
	var kube Kube
	var restConfig, err = clientcmd.BuildConfigFromKubeconfigGetter("", func() (*kubeClientAPI.Config, error) {
		return &config, nil
	})
	if err != nil {
		return nil, err
	}
	kubecli, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	kube.Clientset = kubecli
	kube.config = restConfig
	return &kube, nil
}

func (client *Client) InstallTiller(config kubeClientAPI.Config) (port int32, e error) {
	var kube, err = NewKubeClient(config)
	if err != nil {
		return 0, ErrUnableToInstallTiler{Prefix: "unable to init kube client", Reason: err}
	}

	var deployments, fetchDeploymentsListErr = kube.AppsV1().
		Deployments(defaulTiler.Namespace).List(kubeAPIv1.ListOptions{})
	if fetchDeploymentsListErr != nil {
		return 0, fetchDeploymentsListErr
	}
	var depl, deploymentAlreadyExists = findDepl(deployments.Items, defaulTiler.Name)
	if !deploymentAlreadyExists {
		var d, createTilerErr = kube.AppsV1().
			Deployments(defaulTiler.Namespace).
			Create(&defaulTiler)
		if createTilerErr != nil {
			return 0, ErrUnableToInstallTiler{Prefix: "unable to create tiller deploy", Reason: createTilerErr}
		}
		depl = *d
	}

	var services, fetchServicesErr = kube.CoreV1().
		Services(defaulTiler.Namespace).List(kubeAPIv1.ListOptions{})
	if fetchServicesErr != nil {
		return 0, fetchServicesErr
	}
	var serv, serviceAlreadyExists = findServ(services.Items, defaultTillerService.Name)
	if !serviceAlreadyExists {
		if !deploymentAlreadyExists {
			var port, ok = getFirstPort(depl)
			if !ok {
				return 0, ErrUnableToInstallTiler{Prefix: "invalid tiller deployment: no container ports!"}
			}
			defaultTillerService.Spec.Ports[0].Port = port
		}
		var s, createTillerService = kube.CoreV1().
			Services(defaulTiler.Namespace).Create(&defaultTillerService)
		if createTillerService != nil {
			return 0, ErrUnableToInstallTiler{Prefix: "unable to install tiller service", Reason: createTillerService}
		}
		serv = *s
	}
	return serv.Spec.Ports[0].Port, nil
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
	_ = paragonTiller
	//if err := yaml.Unmarshal([]byte(paragonTiller), &defaulTiler); err != nil {
	//	panic(err)
	//}

	//	fmt.Printf("default Tiller deployment name: %+v\n", defaulTiler)
}

func newInt32(i int32) *int32 {
	return &i
}

func findDepl(list []kubeAppsV1.Deployment, name string) (kubeAppsV1.Deployment, bool) {
	for _, depl := range list {
		if depl.Name == name {
			return *depl.DeepCopy(), true
		}
	}
	return kubeAppsV1.Deployment{}, false
}

func findServ(list []kubeCoreV1.Service, name string) (kubeCoreV1.Service, bool) {
	for _, serv := range list {
		if serv.Name == name {
			return *serv.DeepCopy(), true
		}
	}
	return kubeCoreV1.Service{}, false
}

func getFirstPort(depl kubeAppsV1.Deployment) (int32, bool) {
	for _, container := range depl.Spec.Template.Spec.Containers {
		for _, p := range container.Ports {
			return p.ContainerPort, true
		}
	}
	return -1, false
}
