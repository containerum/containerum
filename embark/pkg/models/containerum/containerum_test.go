package containerum

import (
	"io/ioutil"
	"testing"

	"github.com/go-yaml/yaml"
)

func TestContainerum(test *testing.T) {
	var containerum Containerum
	var data, err = ioutil.ReadFile("testdata/containerum.yaml")
	if err != nil {
		test.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &containerum); err != nil {
		test.Fatal(err)
	}
	test.Log("\n", containerum.Copy())
}
