package builder

import (
	kubeCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var defaultTillerService = kubeCoreV1.Service{
	TypeMeta: v1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: v1.ObjectMeta{
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
