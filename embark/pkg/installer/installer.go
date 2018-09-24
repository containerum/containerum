package installer

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/depsearch"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/models/components"
	"github.com/containerum/containerum/embark/pkg/object"
	"github.com/containerum/containerum/embark/pkg/utils/why"
)

const Containerum = "containerum"

type Installer struct {
	ContainerumConfigPath string
	TempDir               string
	KubectlConfigPath     string
}

func (installer Installer) Install() error {
	if err := installer.setupTempDir(); err != nil {
		return err
	}

	var containerumComponents, loadContDataErr = installer.loadContainerumConfig()
	if loadContDataErr != nil {
		return loadContDataErr
	}

	if err := installer.downloadContainerumIfPresents(containerumComponents); err != nil {
		return err
	}

	if err := installer.downloadUncachedComponents(containerumComponents); err != nil {
		return err
	}

	var rendered, renderErr = builder.RenderComponents(installer.TempDir, containerumComponents, builder.RenderWithValues(map[string]interface{}{}))
	if renderErr != nil {
		return renderErr
	}

	var renderedObjects = make([]object.Object, 0, len(rendered))
	_ = renderedObjects
	var kubeClient, newKubeClientErr = kube.NewKube()
	if newKubeClientErr != nil {
		return newKubeClientErr
	}
	_ = kubeClient
	return nil
}

func (installer Installer) setupTempDir() error {
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

func (installer Installer) loadContainerumConfig() (components.Components, error) {
	return loadContainerumConfig(installer.ContainerumConfigPath)
}

func (installer Installer) downloadContainerumIfPresents(contComponents components.Components) error {
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

func (installer Installer) downloadUncachedComponents(contComponents components.Components) error {
	var chartIndex, buildingChartIndexErr = depsearch.NewSearcher(installer.TempDir)
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
