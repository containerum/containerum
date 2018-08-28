package builder

import (
	"fmt"
	"net/http"

	"github.com/containerum/containerum/embark/pkg/models/requirements"
	"github.com/mholt/archiver"
)

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
