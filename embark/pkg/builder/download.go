package builder

import (
	"net/http"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/mholt/archiver"
)

func downloadDependency(getter *http.Client, dir, componentURI string) error {
	var resp, err = getter.Get(componentURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var errUnpack = archiver.TarGz.Read(resp.Body, dir)
	// bar.Finish()
	if errUnpack != nil {
		return emberr.ErrUnableToDownloadDependencies{Reason: errUnpack}
	}
	return nil
}
