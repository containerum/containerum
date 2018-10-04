// +build IntegrationTests

package kubeTests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/containerum/containerum/embark/pkg/object"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/ericchiang/k8s/apis/meta/v1"
)

var KubeConfig string = os.ExpandEnv("$HOME/.kube/config")

func TestKubeClient(test *testing.T) {
	var client, newKubeErr = kube.NewKube(kube.WithKubectlConfigProvider(kube.LoadKubectlConfigFromPath(KubeConfig)))
	if newKubeErr != nil {
		test.Fatal(newKubeErr)
	}
	var testDeplFile, openTestDeplFileErr = os.Open("testdata/depl.yaml")
	if openTestDeplFileErr != nil {
		test.Fatal(openTestDeplFileErr)
	}
	defer testDeplFile.Close()

	var obj, objErr = object.ObjectFromYAML(testDeplFile)
	if objErr != nil {
		test.Fatal(objErr)
	}
	obj.PatchMeta(func(meta *v1.ObjectMeta) {
		var namespace = "testnamespace"
		meta.Namespace = &namespace
	})
	var createErr = client.Create(obj)
	if createErr != nil {
		test.Fatal(createErr)
	}
}

func init() {
	if KubeConfig == "" {
		panic("KubeConfig variable must be set by ldflags")
	}
}

func Ctx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
