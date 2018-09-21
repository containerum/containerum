package builder

import (
	"io"
	"net/http"
	"runtime"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/mholt/archiver"
)

func DownloadComponent(getter *http.Client, dir, componentURI string) error {
	var resp, err = getter.Get(componentURI)
	if err != nil {
		return emberr.ErrUnableToDownloadDependencies{Reason: err}
	}
	defer resp.Body.Close()
	var errUnpack = archiver.TarGz.Read(GoshedReader{resp.Body}, dir)
	if errUnpack != nil {
		return emberr.ErrUnableToDownloadDependencies{Reason: errUnpack}
	}
	return nil
}

type GoshedReader struct {
	io.Reader
}

func (re GoshedReader) Read(data []byte) (int, error) {
	runtime.Gosched()
	return re.Reader.Read(data)
}
