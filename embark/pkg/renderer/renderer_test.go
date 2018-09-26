package renderer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/containerum/containerum/embark/pkg/models/chart"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/ogetter"
	"github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRenderer(test *testing.T) {
	var getter = ogetter.NewFSObjectGetter("testdata/postgresql/templates")
	var component, rendererErr = Renderer{
		Name:         "postgresql",
		ObjectGetter: getter,
		Values:       testValues(test),
		ObjectsToRender: []string{
			"deployment",
		},
		Constructor: func(reader io.Reader) (kube.Object, error) {
			var buf = &bytes.Buffer{}
			_, err := buf.ReadFrom(reader)
			return mockObject{Buffer: buf}, err
		},
	}.RenderComponent()
	if rendererErr != nil {
		test.Fatal(rendererErr)
	}
	var objects = component.Objects()
	//test.Log(objects)
	for i, object := range objects {
		var mock, ok = object.(mockObject)
		if !ok {
			test.Logf("object %d has invalid type %t", i, object)
			continue
		}
		var t interface{}
		if err := yaml.Unmarshal(mock.Bytes(), &t); err != nil {
			test.Fatalf("%v\n\n%s", err, lines(mock.String()))
		}
	}
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

func testValues(test *testing.T) Values {
	var chartValues map[string]interface{}
	{
		var data, loadTestValuesErr = ioutil.ReadFile("testdata/postgresql/values.yaml")
		assert.Nil(test, loadTestValuesErr)
		assert.Nil(test, yaml.Unmarshal(data, &chartValues))
	}
	var ch chart.Chart
	{
		var data, loadTestValuesErr = ioutil.ReadFile("testdata/postgresql/Chart.yaml")
		assert.Nil(test, loadTestValuesErr)
		assert.Nil(test, yaml.Unmarshal(data, &ch))
	}
	chartValues["strategy"] = "restart-always"
	var values = DefaultValues()
	values.Chart = ch
	return values
}

func lines(text string) string {
	var lines = strings.Split(text, "\n")
	var maxLineNumberTextLen = len(strconv.Itoa(len(lines)))
	var lineNumberAligment = "%" + strconv.Itoa(maxLineNumberTextLen) + "d"
	fmt.Printf("%s\n", lineNumberAligment)
	var buf = &bytes.Buffer{}
	for i, line := range lines {
		fmt.Fprintf(buf, lineNumberAligment+" %s\n", i+1, line)
	}
	return buf.String()
}
