package builder

import (
	chartutil "k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
)

var (
	_ = helm.Client{}
)

type Client struct {
	*helm.Client
}

func NewCLient(host string) *Client {
	return &Client{
		helm.NewClient(
			helm.Host(host),
			helm.ConnectTimeout(60),
		),
	}
}

func (client *Client) Install(toCollection, dir string, valuesFile string) error {
	var ch, errLoadDir = chartutil.LoadDir(dir)
	if errLoadDir != nil {
		return errLoadDir
	}
	var result, errInstallChart = client.Client.InstallReleaseFromChart(ch, toCollection,
		helm.InstallWait(false))
	if errInstallChart != nil {
		return errInstallChart
	}
	_ = result.Release
	return nil
}
