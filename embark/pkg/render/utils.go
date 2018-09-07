package render

import (
	"path"
	"strings"

	"gopkg.in/yaml.v2"
	kubeApsV1 "k8s.io/api/apps/v1"
	kubeBatchAPIv1 "k8s.io/api/batch/v1"
	kubeCoreV1 "k8s.io/api/core/v1"
	kubeExtensionsV1beta1 "k8s.io/api/extensions/v1beta1"
)

// returns filename without extension
func FileNameWithoutExt(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}

func parseDeployment(data string) (kubeApsV1.Deployment, error) {
	var obj = kubeApsV1.Deployment{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseService(data string) (kubeCoreV1.Service, error) {
	var obj = kubeCoreV1.Service{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseVolume(data string) (kubeCoreV1.Volume, error) {
	var obj = kubeCoreV1.Volume{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseIngress(data string) (kubeExtensionsV1beta1.Ingress, error) {
	var obj = kubeExtensionsV1beta1.Ingress{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseSecret(data string) (kubeCoreV1.Secret, error) {
	var obj = kubeCoreV1.Secret{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseJob(data string) (kubeBatchAPIv1.Job, error) {
	var obj = kubeBatchAPIv1.Job{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}

func parseConfigmap(data string) (kubeCoreV1.ConfigMap, error) {
	var obj = kubeCoreV1.ConfigMap{}
	return obj, yaml.Unmarshal([]byte(data), &obj)
}
