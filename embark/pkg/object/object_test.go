package object

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestObjectEncodeDecode(test *testing.T) {
	var testDeplData, loadTestDataErr = ioutil.ReadFile("testdata/depl.yaml")
	if loadTestDataErr != nil {
		test.Fatal(testDeplData)
	}
	var _, err = ObjectFromYAML(bytes.NewReader(testDeplData))
	if err != nil {
		test.Fatal(err)
	}
}
