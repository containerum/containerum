package builder

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"github.com/mholt/archiver"
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

func (client *Client) DownloadRequirements(dir string, reqs requirements.Requirements) error {
	var getter = &http.Client{
		Timeout: 60 * time.Second,
	}
	for _, req := range reqs.Dependencies {
		var addr, err = req.URL()
		if err != nil {
			return err
		}
		var resp, errGet = getter.Get(addr)
		if errGet != nil {
			return errGet
		}
		defer resp.Body.Close()
		if err := archiver.TarGz.Read(resp.Body, path.Join(os.TempDir(), req.FileName())); err != nil {
			return err
		}
	}
	return nil
}

func (client *Client) Install(namespace, dir string, valuesFile string) error {
	var ch, errLoadDir = chartutil.LoadDir(dir)
	if errLoadDir != nil {
		return errLoadDir
	}
	var result, errInstallChart = client.Client.InstallReleaseFromChart(ch, namespace,
		helm.InstallWait(false))
	if errInstallChart != nil {
		return errInstallChart
	}
	_ = result.Release
	return nil
}
