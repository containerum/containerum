package tiller

import (
	"fmt"
	"net/url"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/logger"
	kubeAPIv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeClientAPI "k8s.io/client-go/tools/clientcmd/api"
)

// InstallTiller installs Tiller on Kubernetes cluster, defined in kubeClientAPI.Config and returns Tiller API URL
func InstallTiller(log logger.Logger, config kubeClientAPI.Config) (a string, e error) {
	if log == nil {
		log = logger.StdLogger()
	}
	var kube, err = kube.NewKubeClient(config)
	if err != nil {
		return "", emberr.ErrUnableToInstallTiler{Prefix: "unable to init kube client", Reason: err}
	}

	log.Info("Searching for tiller deployment")
	var deployments, fetchDeploymentsListErr = kube.AppsV1().
		Deployments(defaultTiller.Namespace).List(kubeAPIv1.ListOptions{})
	if fetchDeploymentsListErr != nil {
		return "", fetchDeploymentsListErr
	}
	var depl, deploymentAlreadyExists = findDepl(deployments.Items, defaultTiller.Name)
	if !deploymentAlreadyExists {
		log.Info("Installing tiller deployment")
		var d, createTilerErr = kube.AppsV1().
			Deployments(defaultTiller.Namespace).
			Create(&defaultTiller)
		if createTilerErr != nil {
			return "", emberr.ErrUnableToInstallTiler{Prefix: "unable to create tiller deploy", Reason: createTilerErr}
		}
		depl = *d
	} else {
		log.Info("Tiller deployment already exists")
	}

	log.Info("Searching for tiller service")
	var services, fetchServicesErr = kube.CoreV1().
		Services(defaultTiller.Namespace).List(kubeAPIv1.ListOptions{})
	if fetchServicesErr != nil {
		return "", fetchServicesErr
	}
	var serv, serviceAlreadyExists = findServ(services.Items, defaultTillerService.Name)
	if !serviceAlreadyExists {
		if !deploymentAlreadyExists {
			var port, ok = getFirstPort(depl)
			if !ok {
				return "", emberr.ErrUnableToInstallTiler{Prefix: "invalid tiller deployment: no container ports!"}
			}
			defaultTillerService.Spec.Ports[0].Port = port
		}
		log.Info("Installing tiller service")
		var s, createTillerService = kube.CoreV1().
			Services(defaultTiller.Namespace).Create(&defaultTillerService)
		if createTillerService != nil {
			return "", emberr.ErrUnableToInstallTiler{Prefix: "unable to install tiller service", Reason: createTillerService}
		}
		serv = *s
	} else {
		log.Info("Tiller service already exists")
	}
	var servAddr = &url.URL{
		Scheme: "http",
		Host:   serv.Name + fmt.Sprintf(":%d", serv.Spec.Ports[0].Port),
	}
	return servAddr.String(), nil
}
