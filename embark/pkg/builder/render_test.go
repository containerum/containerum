package builder

import (
	"bytes"
	"testing"

	"k8s.io/helm/pkg/chartutil"
)

func TestRenderChart(test *testing.T) {
	var ch, loadChartErr = chartutil.Load("testdata/volume")
	if loadChartErr != nil {
		test.Fatal(loadChartErr)
	}

	var buf = &bytes.Buffer{}
	var renderChartErr = RenderChart(ch, buf)
	if renderChartErr != nil {
		test.Log(renderChartErr)
	}
	test.Log(buf)
}
