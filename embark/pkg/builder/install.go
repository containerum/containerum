package builder

import (
	"fmt"
	"path"

	"io/ioutil"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/emberr"
	"k8s.io/helm/pkg/helm"
)

func (client *Client) InstallChartWithDependencies(namespace, dir string, valuesFile string) error {
	var chartRequirements, getRequirementsErr = client.getRequirements(dir)
	if getRequirementsErr != nil {
		return emberr.ErrUnableToInstallChart{Prefix: "unable to load requirements", Chart: Containerum, Reason: getRequirementsErr}
	}

	var dependencyGraph, fetchDepsErr = client.FetchAllDeps(chartRequirements, path.Join(dir, "charts"))
	if fetchDepsErr != nil {
		return emberr.ErrUnableToInstallChart{Prefix: "unable to fetch all deps", Chart: Containerum, Reason: fetchDepsErr}
	}
	var installOptions = []helm.InstallOption{
		helm.InstallTimeout(60),
		helm.InstallWait(true), /* blocks until chart is installed */
		helm.InstallDryRun(true),
	}
	if valuesFile != "" {
		fmt.Printf("Using valuues from %q\n", valuesFile)
		var valuesData, loadValuesErr = ioutil.ReadFile(valuesFile)
		if loadValuesErr != nil {
			return emberr.ErrUnableToInstallChart{Prefix: "unable to load values file", Chart: Containerum, Reason: loadValuesErr}
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
			fmt.Printf("Installing %q from %q\n", node, chartDir)
			var ch, errLoadChart = client.LoadChartFromDir(chartDir)
			if errLoadChart != nil {
				return fmt.Errorf("unable to load chart: %v", errLoadChart)
			}
			var _, installErr = client.InstallReleaseFromChart(ch, namespace, installOptions...)
			return installErr
		})
	})
	fmt.Printf("Installing containerum through tiller %q\n", client.host)
	var installErr = installationGraph.Execute(Containerum)
	if installErr != nil {
		return emberr.ErrUnableToInstallChart{Chart: Containerum, Reason: installErr}
	}
	return nil
}
