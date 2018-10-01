package installer

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"

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

func (installer Installer) Install() error {
	if err := installer.SetupTempDir(); err != nil {
		return err
	}
	var containerumComponents, loadContDataErr = installer.LoadContainerumConfig()
	if loadContDataErr != nil {
		return loadContDataErr
	}

	var componentSearcher depsearch.Searcher
	// ? maybe add more cases in the future
	switch installer.Static {
	case true:
		componentSearcher = depsearch.Static()
	default:
		if err := installer.DownloadContainerumChart(containerumComponents); err != nil {
			return err
		}
		if err := installer.DownloadComponents(containerumComponents); err != nil {
			return err
		}
		var buildIndexErr error
		componentSearcher, buildIndexErr = depsearch.FS(installer.TempDir)
		if buildIndexErr != nil {
			return buildIndexErr
		}
	}
	var renderedComponents = make([]renderer.RenderedComponent, 0, len(containerumComponents))
	for _, component := range containerumComponents.Slice() {
		var componentPath, searchComponentErr = componentSearcher.ResolveVersion(component.Name, component.Version)
		if searchComponentErr != nil {
			return searchComponentErr
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
			ObjectGetter:    gettter,
			Constructor: func(reader io.Reader) (kube.Object, error) {
				return object.ObjectFromYAML(reader)
			},
		}.RenderComponent()
		if renderErr != nil {
			return renderErr
		}
		renderedComponents = append(renderedComponents, renderedComponent)
	}
	var kubeClient, newKubeClientErr = kube.NewKube()
	if newKubeClientErr != nil {
		return newKubeClientErr
	}
	_ = kubeClient
	return nil
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
