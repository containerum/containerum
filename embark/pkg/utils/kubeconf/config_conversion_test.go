package kubeconf

import (
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func TestConfigConversion(test *testing.T) {
	var testConfigData, loadTestConfigErr = ioutil.ReadFile("testdata/test_kube_config.yaml")
	if loadTestConfigErr != nil {
		test.Fatal(loadTestConfigErr)
	}
	var config Config
	if err := yaml.Unmarshal(testConfigData, &config); err != nil {
		test.Fatal(err)
	}
	var k8sConfig, configConversionErr = config.ToK8S()
	if configConversionErr != nil {
		test.Fatal(configConversionErr)
	}

	assert.NotEmpty(test, k8sConfig.Clusters)
	for _, cluster := range k8sConfig.Clusters {
		assert.NotEmpty(test, cluster.Cluster.CertificateAuthorityData)
	}

	assert.NotEmpty(test, k8sConfig.AuthInfos)
	for _, user := range k8sConfig.AuthInfos {
		assert.NotEmpty(test, user.AuthInfo.ClientCertificateData)
		assert.NotEmpty(test, user.AuthInfo.ClientKeyData)
	}
	var pp, _ = yaml.Marshal(k8sConfig)
	test.Logf("\n%s", pp)
}
