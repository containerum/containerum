package builder

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"path"
	"sync"
	"text/template"
	"time"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/models/chart"
	"github.com/containerum/containerum/embark/pkg/models/components"
	"github.com/containerum/containerum/embark/pkg/models/release"
	"github.com/containerum/containerum/embark/pkg/ogetter"
	"github.com/containerum/containerum/embark/pkg/utils/spin"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

func DownloadComponents(baseDir string, cont components.Components) error {
	const timeout = 60 * time.Second
	var ctx, done = context.WithTimeout(context.Background(), timeout)
	defer done()
	var downloader, _ = errgroup.WithContext(ctx)
	var client = &http.Client{
		Timeout: timeout,
	}
	var components = cont.Slice()
	var errors = make([]error, 0, len(components))
	var mu sync.Mutex
	var downloadNotify = make(chan struct{}, len(components))
	var spinStop = make(chan struct{})
	go func() {
		var spinner = spin.Loop{
			Frames: []string{".", "o", "0", "@", "*"},
			Prefix: fmt.Sprintf("downloading %d components ", len(components)),
		}
		fmt.Print(spinner.Next())
		var downloaded float64 = 0
		var total = float64(len(components))
		for range downloadNotify {
			downloaded++
			spinner.Prefix = fmt.Sprintf("downloaded %2d ", int(math.Round(100*downloaded/total)))
			fmt.Print(spinner.Next())
		}
		spinner.Erase()
		fmt.Println("")
		close(spinStop)
	}()
	for _, component := range components {
		component := component.Copy()
		downloader.Go(func() error {
			var err = DownloadComponent(client, baseDir, component.URL())
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			}
			downloadNotify <- struct{}{}
			return nil
		})
	}
	downloader.Wait()
	close(downloadNotify)
	<-spinStop

	if len(errors) != 0 {
		return emberr.NewChain(fmt.Errorf("unable to download components"), errors...)
	}
	return nil
}

type RenderedComponent struct {
	components.ComponentWithName
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

type RenderOption func(config *renderConfig)

func RenderWithValues(values map[string]interface{}) RenderOption {
	return func(config *renderConfig) {
		config.mixinValues = values
	}
}

func RenderComponents(baseDir string, cont components.Components, options ...RenderOption) ([]RenderedComponent, error) {
	var config = &renderConfig{
		release: &release.Release{},
	}
	for _, option := range options {
		option(config)
	}
	var components = make([]RenderedComponent, 0, len(cont))
	for _, component := range cont.Slice() {
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
