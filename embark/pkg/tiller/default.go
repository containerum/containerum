package tiller

import (
	kubeAppsV1 "k8s.io/api/apps/v1"
	kubeCoreV1 "k8s.io/api/core/v1"
	kubeAPIv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var defaultTiller = kubeAppsV1.Deployment{
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
		Replicas: func(v int32) *int32 { return &v }(1),
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

var defaultTillerService = kubeCoreV1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "tiller-deploy",
		Namespace: "kube-system",
		Labels: map[string]string{
			"app":  "helm",
			"name": "tiller",
		},
	},
	Spec: kubeCoreV1.ServiceSpec{
		Ports: []kubeCoreV1.ServicePort{
			{Name: "tiller", Port: 44134, Protocol: kubeCoreV1.ProtocolTCP, TargetPort: intstr.FromString("tiller")},
		},
		Type:            kubeCoreV1.ServiceTypeClusterIP,
		SessionAffinity: "None",
		Selector: map[string]string{
			"app":  "helm",
			"name": "tiller",
		},
	},
}
