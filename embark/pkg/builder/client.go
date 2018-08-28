package builder

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"golang.org/x/sync/errgroup"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	_ = helm.Client{}
)

const (
	Containerum = "containerum"
)

type Client struct {
	*helm.Client
	host       string
	downloadMu sync.Mutex
}

func NewCLient(host string) *Client {
	return &Client{
		host: host,
		Client: helm.NewClient(
			helm.Host(host),
			helm.ConnectTimeout(60),
		),
	}
}

// Downloads requirements to target dir.
// Doesn't resolve recursive dependencies
func (client *Client) DownloadRequirements(dir string, reqs requirements.Requirements) error {
	const timeout = 60 * time.Second
	var getter = &http.Client{
		Timeout: timeout,
	}
	if err := MkDirIfNotExists(dir); err != nil {
		return err
	}
	var groupTimeout = time.Duration(len(reqs.Dependencies)) * timeout
	var ctx, done = context.WithTimeout(context.Background(), groupTimeout)
	defer done()
	var group, _ = errgroup.WithContext(ctx)
	for _, req := range reqs.Dependencies {
		req := req
		group.Go(func() error {
			return client.downloadDependency(getter, dir, req)
		})
	}
	return group.Wait()
}

// Loads chart from dir
func (client *Client) LoadChartFromDir(dir string) (*chart.Chart, error) {
	return chartutil.LoadDir(dir)
}

// Extract requirements from chart dir if exist
func (client *Client) getRequirements(dir string) (requirements.Requirements, error) {
	var ch, err = client.LoadChartFromDir(dir)
	if err != nil {
		return requirements.Requirements{}, err
	}
	return requirements.FromChart(ch)
}
