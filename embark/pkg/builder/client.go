package builder

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/logger"
	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"golang.org/x/sync/errgroup"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

const (
	Containerum = "containerum"
)

type Client struct {
	downloadMu sync.Mutex
	kube       *kube.Kube
	log        logger.Logger
	timeout    time.Duration
}

func NewCLient(options ...clientOptions) (*Client, error) {
	var clientConfig = DefaultClientOptionsPtr().Merge(options...)

	var kubeConfig, getKubeConfigErr = clientConfig.kubeConfig()
	if getKubeConfigErr != nil {
		return nil, getKubeConfigErr
	}

	var kubeClient, newKubeClientErr = kube.NewKubeClient(kubeConfig)
	if newKubeClientErr != nil {
		return nil, newKubeClientErr
	}

	return &Client{
		kube:    kubeClient,
		log:     clientConfig.log,
		timeout: clientConfig.timeout,
	}, nil
}

// Downloads requirements to target dir.
// Doesn't resolve recursive dependencies
func (client *Client) DownloadRequirements(dir string, reqs requirements.Requirements) error {
	var getter = &http.Client{
		Timeout: client.timeout,
	}
	if err := MkDirIfNotExists(dir); err != nil {
		return err
	}
	var groupTimeout = time.Duration(len(reqs.Dependencies)) * client.timeout
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
