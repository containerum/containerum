package compression

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
)

// Gzip compression support
type Gzip struct {
	*Options
}

// NewGzip creates a Gzip + Options variable
func NewGzip() *Gzip {
	gz := new(Gzip)
	gz.Options = new(Options)
	return gz
}

// Compress to gzip
func (gz *Gzip) Compress(content []byte) ([]byte, error) {
	if !gz.Options.Compress {
		return content, nil
	}

	// method
	var m int
	switch gz.Options.Method {
	case "NoCompression":
		m = flate.NoCompression
		break
	case "BestSpeed":
		m = flate.BestSpeed
		break
	case "BestCompression":
		m = flate.BestCompression
		break
	default:
		m = flate.DefaultCompression
		break
	}

	// compress
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, m)
	if err != nil {
		return nil, err
	}

	// insert content
	_, err = w.Write(content)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	// compressed content
	return b.Bytes(), nil
}
