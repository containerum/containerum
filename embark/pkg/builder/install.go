package builder

import (
	"fmt"
	"path"

	"io/ioutil"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"k8s.io/helm/pkg/helm"
)

func (client *Client) InstallChartWithDependencies(namespace, dir string, valuesFile string) error {
	var chartRequirements, getRequirementsErr = client.getRequirements(dir)
	if getRequirementsErr != nil {
		return getRequirementsErr
	}

	var dependencyGraph, fetchDepsErr = client.FetchAllDeps(chartRequirements, path.Join(dir, "charts"))
	if fetchDepsErr != nil {
		return fetchDepsErr
	}
	var installOptions = []helm.InstallOption{
		helm.InstallWait(true), /* blocks until chart is installed */
	}
	if valuesFile != "" {
		var valuesData, loadValuesErr = ioutil.ReadFile(valuesFile)
		if loadValuesErr != nil {
			return loadValuesErr
		}
		installOptions = append(installOptions,
			helm.ValueOverrides(valuesData))
	}

	var installationGraph = make(cgraph.Graph)
	dependencyGraph.Walk(Containerum, func(node string, _ []string, children []string) {
		installationGraph.AddNode(node, children, func() error {
			var chartDir string
			switch node {
			case Containerum:
				chartDir = dir
			default:
				chartDir = path.Join(dir, "charts", node)
			}
			fmt.Printf("Installing %q\n", node)
			var _, installErr = client.InstallRelease(
				namespace, /* kubernetes namespace */
				chartDir,  /*dir with chart */
				installOptions...)
			return installErr
		})
	})
	return installationGraph.Execute(Containerum)
}
