package readers

import (
	"compress/gzip"
	"io"
)

type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (g GzipReader) Read(p []byte) (n int, err error) {
	return g.zr.Read(p)
}

func (g GzipReader) Close() error {
	return g.zr.Close()
}

func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}
