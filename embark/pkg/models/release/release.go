package release

import "k8s.io/helm/pkg/proto/hapi/release"

type Release struct {
	release.Release
	Service string
}
