package builder

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestBuilder_Build(test *testing.T) {
	var target, errTarget = ioutil.ReadFile("testdata/target.yaml")
	if errTarget != nil {
		test.Fatal(errTarget)
	}
	var values, errValues = ioutil.ReadFile("testdata/values.yaml")
	if errValues != nil {
		test.Fatal(errValues)
	}
	var helpers, errHelpers = ioutil.ReadFile("testdata/helpers.tpl")
	if errHelpers != nil {
		test.Fatal(errHelpers)
	}
	var output = &bytes.Buffer{}
	var errBuild = Builder{
		Template: string(target),
		Values:   string(values),
		Helpers:  string(helpers),
		Output:   output,
	}.Build()
	if errBuild != nil {
		test.Fatal(errBuild)
	}
	test.Log(output)
}
