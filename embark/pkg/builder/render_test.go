package builder

import (
	"testing"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
)

func TestRenderChart(test *testing.T) {
	var ch, loadChartErr = chartutil.Load("testdata/postgresql")
	if loadChartErr != nil {
		test.Fatal(loadChartErr)
	}
	var rendered, renderChartErr = RenderChart(ch)
	if renderChartErr != nil {
		test.Fatal(renderChartErr)
	}
	var result, _ = yaml.Marshal(rendered)
	test.Log(string(result))
}
