package builder

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"github.com/mholt/archiver"
	"golang.org/x/sync/errgroup"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	_ = helm.Client{}
)

type Client struct {
	*helm.Client
	downloadMu sync.Mutex
}

func NewCLient(host string) *Client {
	return &Client{
		Client: helm.NewClient(
			helm.Host(host),
			helm.ConnectTimeout(60),
		),
	}
}

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

func (client *Client) downloadDependency(getter *http.Client, dir string, chart requirements.Dependency) error {

	client.downloadMu.Lock() // Lock

	var resp, err = func() (*http.Response, error) {
		fmt.Printf("Downloading and unpacking %q...", chart.FileName())

		defer client.downloadMu.Unlock() // Unlock

		var addr, err = chart.URL()
		if err != nil {
			return nil, ErrUnableToFetchChart{Chart: chart.Name, Reason: err}
		}
		var resp, errGet = getter.Get(addr)
		if errGet != nil {
			return nil, ErrUnableToFetchChart{Chart: chart.Name, Reason: errGet}
		}
		var size = resp.ContentLength
		if size <= 0 {
			size = 2 * 1 << 10
		}
		fmt.Printf("%d bytes\n", size)
		return resp, nil
	}()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// var bar = pb.New(int(resp.ContentLength))
	// bar.Start()
	var errUnpack = archiver.TarGz.Read(resp.Body, dir)
	// bar.Finish()
	if errUnpack != nil {
		return ErrUnableToFetchChart{Chart: chart.Name, Reason: errUnpack}
	}
	return nil
}

func (client *Client) LoadChartFromDir(dir string) (*chart.Chart, error) {
	return chartutil.LoadDir(dir)
}

func (client *Client) getRequirements(dir string) (requirements.Requirements, error) {
	var ch, err = client.LoadChartFromDir(dir)
	if err != nil {
		return requirements.Requirements{}, err
	}
	return requirements.FromChart(ch)
}

func (client *Client) FetchAllDeps(rootRequirements requirements.Requirements, dir string) (cgraph.Graph, error) {
	//	if err := client.DownloadRequirements(dir, rootRequirements); err != nil {
	//		return err
	//	}
	var deps = requirements.NewQueue(len(rootRequirements.Dependencies))
	deps.Push(rootRequirements.Dependencies...)
	var getter = &http.Client{
		Timeout: 60 * time.Second,
	}
	var downloaded = map[string]bool{}

	var graph = make(cgraph.Graph)
	graph.AddNode("containerum", rootRequirements.Names(), func() error {
		fmt.Printf("installing %q\n", "containerum")
		return nil
	})
	for dep := range deps {
		dep := dep

		fmt.Printf("Resolving %q, %d deps left\n", dep, len(deps))
		var depDep []string
		if !downloaded[dep.Name] {
			var depDir = path.Join(dir, dep.Name)
			if err := client.downloadDependency(getter, dir, dep); err != nil {
				fmt.Println(err)
				continue
			}
			var depReq, errDepReq = client.getRequirements(depDir)
			if errDepReq != nil {
				if !strings.Contains(errDepReq.Error(), ".yaml not found") {
					return nil, errDepReq
				}
			}

			var depChart, errLoadChart = client.LoadChartFromDir(depDir)
			if errLoadChart != nil {
				fmt.Println(errLoadChart)
				continue
			}

			if len(depChart.GetDependencies()) == 0 {
				fmt.Printf("\t%q depends on %v\n", dep.Name, depReq.Dependencies)
				deps.Push(depReq.Dependencies...)
				depDep = depReq.Names()
			} else {
				fmt.Printf("\tDeps of %q are already vendored in 'charts' dir\n", dep.Name)
			}
			downloaded[dep.Name] = true
		} else {
			fmt.Printf("%\tq is already fetched", dep)
		}

		fmt.Printf("\tAdding %q to graph\n", dep)
		graph.AddNode(dep.Name, depDep, func() error {
			fmt.Printf("installing %q\n", dep)
			return nil
		})

		if len(deps) == 0 {
			close(deps)
		}
	}
	fmt.Printf("Charts to install:\n\t%s\n", strings.Join(graph.Nodes(), "\n\t"))
	return graph, nil
}

func (client *Client) install(namespace, dir string, valuesFile string) error {
	var ch, errLoadDir = chartutil.LoadDir(dir)
	if errLoadDir != nil {
		return errLoadDir
	}
	var result, errInstallChart = client.Client.InstallReleaseFromChart(ch, namespace,
		helm.InstallWait(true))
	if errInstallChart != nil {
		return errInstallChart
	}
	_ = result.Release
	return nil
}
