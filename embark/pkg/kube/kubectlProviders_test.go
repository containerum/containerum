package kube

import (
	"io/ioutil"
	"testing"
)

func TestDecodeConfig(test *testing.T) {
	var data, loadDataErr = ioutil.ReadFile("testdata/test_kube_config.yaml")
	if loadDataErr != nil {
		test.Fatal(loadDataErr)
	}
	var config, decodeErr = DecodeConfig(data)
	if decodeErr != nil {
		test.Fatal(decodeErr)
	}
	test.Logf("%#v", config.Clusters[0])
}
