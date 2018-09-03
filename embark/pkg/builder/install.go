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

	var log = client.log

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
		log.Infof("Using values from %q\n", valuesFile)
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
	log.Infof("Installing containerum through tiller %q\n", client.host)
	var installErr = installationGraph.Execute(Containerum)
	if installErr != nil {
		return emberr.ErrUnableToInstallChart{Chart: Containerum, Reason: installErr}
	}
	return nil
}

// NOTE:
func (client *Client) Install(namespace, dir string) error {
	var ch, loadChartErr = client.LoadChartFromDir(dir)
	if loadChartErr != nil {
		return loadChartErr
	}
	var rendered, err = RenderChart(ch)
	if err != nil {
		return err
	}

	for _, depl := range rendered.Deployments {
		var _, createDeplErr = client.
			kube.
			AppsV1().
			Deployments(namespace).
			Create(&depl)
		if createDeplErr != nil {
			return createDeplErr
		}
	}

	for _, serv := range rendered.Services {
		var _, createServErr = client.
			kube.
			CoreV1().
			Services(namespace).
			Create(&serv)
		if createServErr != nil {
			return err
		}
	}

	for _, ingr := range rendered.Ingresses {
		var _, createIngr = client.
			kube.ExtensionsV1beta1().
			Ingresses(namespace).
			Create(&ingr)
		if createIngr != nil {
			return createIngr
		}
	}

	return nil
}
