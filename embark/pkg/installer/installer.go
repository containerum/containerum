package installer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/containerum/containerum/embark/pkg/builder"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/kube"
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
	var containerumConfig, loadContDataErr = loadContainerumConfig(installer.ContainerumConfigPath)
	if loadContDataErr != nil {
		return loadContDataErr
	}
	if err := builder.DowloadComponents(installer.TempDir, containerumConfig); err != nil {
		return err
	}

	var rendered, renderErr = builder.RenderComponents(installer.TempDir, containerumConfig, builder.RenderWithValues(map[string]interface{}{}))
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
