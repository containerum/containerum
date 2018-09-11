package builder

import (
	"bytes"
	"context"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/containerum/containerum/embark/pkg/models/release"

	"github.com/containerum/containerum/embark/pkg/models/chart"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/models/containerum"
	"github.com/containerum/containerum/embark/pkg/ogetter"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

func DowloadComponents(baseDir string, cont containerum.Containerum) error {
	const timeout = 60 * time.Second
	var ctx, done = context.WithTimeout(context.Background(), timeout)
	defer done()
	var downloader, _ = errgroup.WithContext(ctx)
	var client = &http.Client{
		Timeout: timeout,
	}
	for _, component := range cont.Components() {
		downloader.Go(downloadDependency(client, baseDir, component.URL()))
	}
	return downloader.Wait()
}

type RenderedComponent struct {
	containerum.ComponentWithName
	Objects map[string]*bytes.Buffer
}

func MergeValues(values ...map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for _, mixin := range values {
		for k, v := range mixin {
			result[k] = v
		}
	}
	return result
}

type m = map[string]interface{}

type renderConfig struct {
	mixinValues map[string]interface{}
	release     *release.Release
}

func RenderComponents(baseDir string, cont containerum.Containerum, configs ...renderConfig) ([]RenderedComponent, error) {

	var config = renderConfig{
		release: &release.Release{},
	}
	for _, layer := range configs {
		if layer.release != nil {
			config.release = layer.release
		}
		if layer.mixinValues != nil {
			config.mixinValues = layer.mixinValues
		}
	}

	var components = make([]RenderedComponent, 0, len(cont))
	for _, component := range cont.Components() {
		var componentPath = path.Join(baseDir, component.Name)
		var objectGetter = ogetter.NewFSObjectGetter(componentPath)

		var valuesFromFile = map[string]interface{}{}
		{
			var serializedValues = &bytes.Buffer{}
			if err := objectGetter.Object("values", serializedValues); err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(serializedValues.Bytes(), &valuesFromFile); err != nil {
				return nil, err
			}
		}

		var ch chart.Chart
		{
			var serializedChart = &bytes.Buffer{}
			if err := objectGetter.Object("Chart", serializedChart); err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(serializedChart.Bytes(), &ch); err != nil {
				return nil, err
			}
		}

		var componentTemplatePath = path.Join(componentPath, "templates")
		var objectComponentsGetter = ogetter.NewFSObjectGetter(componentTemplatePath)
		var rawSerializedComponents, retrieveObjectErr = ogetter.RetrieveObjects(
			objectComponentsGetter,
			component.ObjectNames()...)
		if retrieveObjectErr != nil {
			return nil, retrieveObjectErr
		}

		var tmpl = template.New(component.Name)
		var result = make(map[string]*bytes.Buffer)

		var values = MergeValues(
			m{
				"Chart":   ch,
				"Release": config.release,
			},
			m{
				component.Name: m{
					"name":    component.Name,
					"version": component.Version,
				},
			},
			component.Values,
			config.mixinValues)
		if err := rawSerializedComponents.Render(tmpl, result, values); err != nil {
			return nil, err
		}
		components = append(components, RenderedComponent{
			ComponentWithName: component.Copy(),
			Objects:           result,
		})
	}
	return components, nil
}

func BuildGraph(baseDir string, components []RenderedComponent) (cgraph.Graph, error) {
	var dependencyGraph = cgraph.NewGraph()
	for _, component := range components {
		component := component
		dependencyGraph.AddNode(component.Name, component.DependsOn, func() error {
			return nil
		})
	}
	return dependencyGraph, nil
}
