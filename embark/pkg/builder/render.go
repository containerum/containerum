package builder

import (
	"fmt"
	"io"
	"path"
	"strings"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const notesFileSuffix = "NOTES.txt"

type renderOptions struct {
	Values map[string]interface{}
}

func RenderWithValues(values map[string]interface{}) renderOptions {
	return renderOptions{
		Values: values,
	}
}

func RenderChart(chart *chart.Chart, output io.Writer, options ...renderOptions) error {
	var chartValues, readChartValuesErr = chartutil.ReadValues([]byte(chart.GetValues().Raw))
	if readChartValuesErr != nil {
		return readChartValuesErr
	}
	var renderConfig = renderOptions{
		Values: chartValues,
	}

	for _, option := range options {
		if option.Values != nil {
			renderConfig.Values = option.Values
		}
	}

	var renderEngine = engine.New()
	var targets, renderErr = renderEngine.Render(chart, renderConfig.Values)
	if renderErr != nil {
		return renderErr
	}

	var notes = make([]string, 0)
	for k, v := range targets {
		if strings.HasSuffix(k, notesFileSuffix) {
			// Only apply the notes if it belongs to the parent chart
			// Note: Do not use filePath.Join since it creates a path with \ which is not expected
			if k == path.Join(chart.Metadata.Name, "templates", notesFileSuffix) {
				notes = append(notes, v)
			}
			delete(targets, k)
		} else {
			var _, writeTargetErr = fmt.Fprintf(output, "\n# %s\n%s\n", k, v)
			if writeTargetErr != nil {
				return writeTargetErr
			}
		}
	}
	return nil
}
