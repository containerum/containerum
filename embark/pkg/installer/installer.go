package installer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/depsearch"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/models/components"
	"github.com/containerum/containerum/embark/pkg/render"
	"github.com/containerum/containerum/embark/pkg/utils/why"
)

type Installer struct {
	ContainerumConfigPath string
	TempDir               string
	KubectlConfigPath     string
}

func (installer Installer) Install() error {
	if installer.TempDir == "" {
		installer.TempDir = path.Join(os.TempDir(), "embark")
	}
	if err := os.MkdirAll(installer.TempDir, os.ModePerm|os.ModeDir); err != nil && !os.IsExist(err) {
		return emberr.ErrUnableToCreateTempDir{
			Path:   installer.TempDir,
			Reason: err,
		}
	}
	var containerumComponents, loadContDataErr = loadContainerumConfig(installer.ContainerumConfigPath)
	if loadContDataErr != nil {
		return loadContDataErr
	}

	if containerumComponents.Contains(render.Containerum) {
		var getter = &http.Client{
			Timeout: 10 * time.Second,
		}
		var downloadContainerumErr = builder.DownloadComponent(getter,
			installer.TempDir,
			containerumComponents.MustGet(render.Containerum).URL())
		if downloadContainerumErr != nil {
			return downloadContainerumErr
		}
	}

	var chartIndex, buildingChartIndexErr = depsearch.NewSearcher(installer.TempDir)
	if buildingChartIndexErr != nil {
		return buildingChartIndexErr
	}

	var notDownloadedComponents = containerumComponents.
		Filter(func(component components.ComponentWithName) bool {
			return !chartIndex.Contains(component.Name)
		})

	if notDownloadedComponents.Len() > 0 {
		why.Print("Components to download", notDownloadedComponents.Names()...)
		if err := builder.DownloadComponents(installer.TempDir, notDownloadedComponents); err != nil {
			return err
		}
		// ! rebuild index!
		chartIndex, buildingChartIndexErr = depsearch.NewSearcher(installer.TempDir)
		if buildingChartIndexErr != nil {
			return buildingChartIndexErr
		}
	}

	var rendered, renderErr = builder.RenderComponents(installer.TempDir, containerumComponents, builder.RenderWithValues(map[string]interface{}{}))
	if renderErr != nil {
		return renderErr
	}
	var errs []error
	for _, component := range rendered {
		for objectName, object := range component.Objects {
			var fname = path.Join(installer.TempDir, component.Name+"_"+objectName+".yaml")
			if err := ioutil.WriteFile(fname, object.Bytes(), os.ModePerm); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return emberr.NewChain(fmt.Errorf("unable write rendered components"), errs...)
	}
	var kubeClient, newKubeClientErr = kube.NewKube()
	if newKubeClientErr != nil {
		return newKubeClientErr
	}
	_ = kubeClient
	return nil
}
