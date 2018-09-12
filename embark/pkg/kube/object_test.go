package kube

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestObjectEncodeDecode(test *testing.T) {
	var testDeplData, loadTestDataErr = ioutil.ReadFile("testdata/depl.yaml")
	if loadTestDataErr != nil {
		test.Fatal(testDeplData)
	}
	var objFromYAML, err = ObjectFromYAML(bytes.NewReader(testDeplData))
	if err != nil {
		test.Fatal(err)
	}

	assert.Equal(test, testDepl(test, testDeplData), objFromYAML.body)
}

func testDepl(test *testing.T, deplData []byte) map[string]interface{} {
	var depl map[string]interface{}
	if err := yaml.Unmarshal(deplData, &depl); err != nil {
		test.Fatal(err)
	}
	return depl
}
