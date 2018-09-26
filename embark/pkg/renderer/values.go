package renderer

import (
	"github.com/containerum/containerum/embark/pkg/models/chart"
	"github.com/containerum/containerum/embark/pkg/models/release"
	kubeRelease "k8s.io/helm/pkg/proto/hapi/release"
)

type Values struct {
	Values  map[string]interface{}
	Chart   chart.Chart
	Release release.Release
}

func DefaultValues() Values {
	return Values{
		Release: release.Release{
			Release: kubeRelease.Release{
				Name: "containerum",
			},
			Service: "containerum",
		},
	}
}
