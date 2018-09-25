package renderer

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/ogetter"
	"github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRenderer(test *testing.T) {
	var getter = ogetter.NewFSObjectGetter("testdata/postgresql")
	test.Log(getter.ObjectNames())

	var component, rendererErr = Renderer{
		Name:   "postgresql",
		Getter: getter,
		Values: testValues(test),
		Contstructor: func(reader io.Reader) (kube.Object, error) {
			var buf = &bytes.Buffer{}
			_, err := buf.ReadFrom(reader)
			return mockObject{Buffer: buf}, err
		},
	}.RenderComponent()
	if rendererErr != nil {
		test.Fatal(rendererErr)
	}
	test.Log(component.Objects())
}

var (
	_ kube.Object = mockObject{}
)

type mockObject struct {
	*bytes.Buffer
}

func (mockObject) Kind() string {
	return "mock object"
}

func (mockObject) GetMetadata() *v1.ObjectMeta {
	return &v1.ObjectMeta{}
}

func testValues(test *testing.T) map[string]interface{} {
	var values map[string]interface{}
	var data, loadTestValuesErr = ioutil.ReadFile("testdata/postgresql/values.yaml")
	assert.Nil(test, loadTestValuesErr)
	assert.Nil(test, yaml.Unmarshal(data, &values))
	return values
}
