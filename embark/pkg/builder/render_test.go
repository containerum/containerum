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
	var _, renderChartErr = RenderChart(ch)
	if renderChartErr != nil {
		test.Fatal(renderChartErr)
	}
	test.Log(buf)
}
