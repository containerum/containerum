package chart

import "k8s.io/helm/pkg/proto/hapi/chart"

type Chart chart.Metadata

/*
		Name        string       `json:"name"`
		Version     string       `json:"version"`
		AppVersion  string       `json:"appVersion"`
		Description string       `json:"description"`
		Engine      string       `json:"engine"`
		Home        string       `json:"home"`
		Icon        string       `json:"icon"`
		Keywords    []string     `json:"keywords"`
		Maintainers []Maintainer `json:"maintainers"`
		Sources     []string     `json:"sources"`

}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

*/
