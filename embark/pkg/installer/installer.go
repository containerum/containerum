package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/containerum/containerum/embark/pkg/cgraph"

	"github.com/containerum/containerum/embark/pkg/object"
	"github.com/containerum/containerum/embark/pkg/ogetter"

	"github.com/containerum/containerum/embark/pkg/renderer"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/depsearch"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/models/components"
	"github.com/containerum/containerum/embark/pkg/utils/why"
)

const Containerum = "containerum"

type Installer struct {
	Static                bool
	ContainerumConfigPath string
	TempDir               string
	KubectlConfigPath     string
}

type GraphExecuter = func() error

func (installer Installer) Install() error {
	if err := installer.SetupTempDir(); err != nil {
		return err
	}
	var containerumComponents, loadContDataErr = installer.LoadContainerumConfig()
	if loadContDataErr != nil {
		return loadContDataErr
	}
	var componentSearcher, loadDependenciesErr = installer.LoadDependencies(containerumComponents)
	if loadDependenciesErr != nil {
		return loadDependenciesErr
	}
	var renderedComponents, renderComponentsErr = installer.RenderComponents(componentSearcher, containerumComponents)
	if renderComponentsErr != nil {
		return renderComponentsErr
	}
	var install, buildInstallationGraphErr = installer.BuildInstallationGraph(renderedComponents)
	if buildInstallationGraphErr != nil {
		return buildInstallationGraphErr
	}
	return install()
}

func (installer Installer) BuildInstallationGraph(renderedComponents []renderer.RenderedComponent) (GraphExecuter, error) {
	var kubeClient, newKubeClientErr = kube.NewKube()
	if newKubeClientErr != nil {
		return nil, newKubeClientErr
	}
	var gr = cgraph.NewGraph()
	for _, component := range renderedComponents {
		component := component
		var dependencies = component.DependsOn()
		gr.AddNode(component.Name(), dependencies, func() error {
			return component.ForEachObject(func(obj kube.Object) error {
				return kubeClient.Create(obj)
			})
		})
	}
	var sinks = gr.Sinks()
	if len(sinks) == 0 {
		return nil, fmt.Errorf("unable to exectude totally cycled graph")
	}
	return func() error { return gr.Execute(sinks...) }, nil
}

func (installer Installer) SetupTempDir() error {
	if installer.TempDir == "" {
		installer.TempDir = path.Join(os.TempDir(), "embark")
	}
	if err := os.MkdirAll(installer.TempDir, os.ModePerm|os.ModeDir); err != nil && !os.IsExist(err) {
		return emberr.ErrUnableToCreateTempDir{
			Path:   installer.TempDir,
			Reason: err,
		}
	}
	return nil
}

func (installer Installer) LoadContainerumConfig() (components.Components, error) {
	return loadContainerumConfig(installer.ContainerumConfigPath)
}

func (installer Installer) DownloadContainerumChart(contComponents components.Components) error {
	if contComponents.Contains(Containerum) {
		var getter = &http.Client{
			Timeout: 10 * time.Second,
		}
		var downloadContainerumErr = builder.DownloadComponent(getter,
			installer.TempDir,
			contComponents.MustGet(Containerum).URL())
		if downloadContainerumErr != nil {
			return downloadContainerumErr
		}
	}
	return nil
}

func (installer Installer) DownloadComponents(contComponents components.Components) error {
	var chartIndex, buildingChartIndexErr = depsearch.FS(installer.TempDir)
	if buildingChartIndexErr != nil {
		return buildingChartIndexErr
	}
	var notDownloadedComponents = contComponents.
		Filter(func(component components.ComponentWithName) bool {
			return !chartIndex.Contains(component.Name)
		})
	if notDownloadedComponents.Len() > 0 {
		why.Print("Components to download", notDownloadedComponents.Names()...)
		if err := builder.DownloadComponents(installer.TempDir, notDownloadedComponents); err != nil {
			return err
		}
	}
	return nil
}

func (installer Installer) RenderComponents(componentSearcher depsearch.Searcher, containerumComponents components.Components) ([]renderer.RenderedComponent, error) {
	var renderedComponents = make([]renderer.RenderedComponent, 0, len(containerumComponents))
	for _, component := range containerumComponents.Slice() {
		var componentPath, searchComponentErr = componentSearcher.ResolveVersion(component.Name, component.Version)
		if searchComponentErr != nil {
			return nil, searchComponentErr
		}
		var gettter ogetter.ObjectGetter
		if installer.Static {
			gettter = ogetter.NewEmbeddedFSObjectGetter(path.Join(componentPath, "templates"))
		} else {
			gettter = ogetter.NewFSObjectGetter(componentPath)
		}
		var renderedComponent, renderErr = renderer.Renderer{
			Name:            component.Name,
			ObjectsToRender: component.Objects,
			DependsOn:       component.DependsOn,
			ObjectGetter:    gettter,
			Constructor: func(reader io.Reader) (kube.Object, error) {
				return object.ObjectFromYAML(reader)
			},
		}.RenderComponent()
		if renderErr != nil {
			return nil, renderErr
		}
		renderedComponents = append(renderedComponents, renderedComponent)
	}
	return renderedComponents, nil
}

func (installer Installer) LoadDependencies(containerumComponents components.Components) (depsearch.Searcher, error) {
	var componentSearcher depsearch.Searcher
	// ? maybe add more cases in the future
	switch installer.Static {
	case true:
		componentSearcher = depsearch.Static()
	default:
		if err := installer.DownloadContainerumChart(containerumComponents); err != nil {
			return componentSearcher, err
		}
		if err := installer.DownloadComponents(containerumComponents); err != nil {
			return componentSearcher, err
		}
		var buildIndexErr error
		componentSearcher, buildIndexErr = depsearch.FS(installer.TempDir)
		if buildIndexErr != nil {
			return componentSearcher, buildIndexErr
		}
	}
	return componentSearcher, nil
}
