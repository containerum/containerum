package builder

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/emberr"
	kubeCoreV1 "k8s.io/api/core/v1"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	_ = kubeCoreV1.PersistentVolumeClaim{}
)

const (
	Deployment = "deployment"
	Ingress    = "ingress"
	Service    = "service"
	Volume     = "volume"
	Job        = "job"
	Configmap  = "configmap"
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

	if valuesFile != "" {
		log.Infof("Using values from %q\n", valuesFile)
		var valuesData, loadValuesErr = ioutil.ReadFile(valuesFile)
		if loadValuesErr != nil {
			return emberr.ErrUnableToInstallChart{Prefix: "unable to load values file", Chart: Containerum, Reason: loadValuesErr}
		}
		// TODO: add values to config
		_ = valuesData
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
			log.Infof("Installing %q from %q\n", node, chartDir)
			var ch, errLoadChart = client.LoadChartFromDir(chartDir)
			if errLoadChart != nil {
				return fmt.Errorf("unable to load chart: %v", errLoadChart)
			}
			for _, dep := range ch.GetDependencies() {
				var installChartDepErr = client.install(namespace, dep, nil)
				if installChartDepErr != nil {
					return installChartDepErr
				}
			}
			return client.install(namespace, ch, nil)
		})
	})
	log.Infof("Installing containerum\n")
	var installErr = installationGraph.Execute(Containerum)
	if installErr != nil {
		return emberr.ErrUnableToInstallChart{Chart: Containerum, Reason: installErr}
	}
	return nil
}

func (client *Client) install(namespace string, ch *chart.Chart, values chartutil.Values) error {
	var rendered, err = RenderChart(ch, renderOptions{
		Values: values,
	})
	if err != nil {
		return err
	}

	var installationOrder = []string{Volume, Configmap, Job, Deployment, Service, Ingress}

	for _, resourceType := range installationOrder {
		switch resourceType {
		case Volume:
			// FIXME: fix volume creation
			/*
				for _, vol := range rendered.Volumes {
					_ = vol
					var _, createVolumeErr = client.
						kube.CoreV1().
						PersistentVolumeClaims("namespace").
						Create(vol.PersistentVolumeClaim)
					if createVolumeErr != nil {
						return createVolumeErr
					}
				}
			*/
		case Configmap:
			for _, configmap := range rendered.Configs {
				var _, createConfigmapErr = client.
					kube.
					CoreV1().
					ConfigMaps(namespace).
					Create(&configmap)
				if createConfigmapErr != nil {
					return createConfigmapErr
				}
			}
		case Job:
			for _, job := range rendered.Jobs {
				var _, createJobErr = client.
					kube.
					BatchV1().
					Jobs(namespace).
					Create(&job)
				if createJobErr != nil {
					return createJobErr
				}
			}
		case Deployment:
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
		case Service:
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
		case Ingress:
			for _, ingr := range rendered.Ingresses {
				var _, createIngrErr = client.
					kube.ExtensionsV1beta1().
					Ingresses(namespace).
					Create(&ingr)
				if createIngrErr != nil {
					return createIngrErr
				}
			}
		default:
			panic(fmt.Sprintf("[embark.pkg.builder.Client.Install] unexected installation order state: %q", resourceType))
		}
	}

	return nil
}
