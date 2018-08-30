package tiller

import (
	kubeAppsV1 "k8s.io/api/apps/v1"
	kubeCoreV1 "k8s.io/api/core/v1"
)

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
