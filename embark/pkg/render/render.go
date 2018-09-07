package render

import (
	"bytes"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const (
	notesFileSuffix = "NOTES.txt"
	Containerum     = "containerum"
	Deployment      = "deployment"
	Service         = "service"
	Ingress         = "ingress"
	Volume          = "volume"
)

type renderOptions struct {
	Values map[string]interface{}
}

func (options renderOptions) Merge(another ...renderOptions) renderOptions {
	for _, anotherOption := range another {
		if anotherOption.Values != nil {
			options.Values = anotherOption.Values
		}
	}
	return options
}

func RenderWithValues(values map[string]interface{}) renderOptions {
	return renderOptions{
		Values: values,
	}
}

type RenderedChart map[string]bytes.Buffer

func RenderChart(ch *chart.Chart, options ...renderOptions) (*RenderedChart, error) {
	var renderConfig = renderOptions{}
	{
		var chartValuesCapsErr error
		renderConfig.Values, chartValuesCapsErr = chartutil.ToRenderValuesCaps(ch,
			&chart.Config{
				Raw: ch.GetValues().Raw,
				//	Values: chartValues,
			},
			chartutil.ReleaseOptions{
				Name:      Containerum,
				Namespace: Containerum,
				IsInstall: true,
			},
			&chartutil.Capabilities{})
		if chartValuesCapsErr != nil {
			return nil, chartValuesCapsErr
		}
	}
	renderConfig.Merge(options...)

	var coalesceErr error

	if false {
		renderConfig.Values, coalesceErr = chartutil.CoalesceValues(ch, &chart.Config{
			Raw: func() string {
				var data, _ = yaml.Marshal(renderConfig.Values)
				return string(data)
			}(),
		})
		if coalesceErr != nil {
			return nil, coalesceErr
		}
	}
	var rendered = RenderedChart{}
	var renderEngine = engine.New()
	var _, renderErr = renderEngine.Render(ch, renderConfig.Values)
	if renderErr != nil {
		return nil, renderErr
	}
	return &rendered, nil
}
