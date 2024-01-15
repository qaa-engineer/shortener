package writers

import (
	"compress/gzip"
	"net/http"
)

type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *GzipWriter) Write(bytes []byte) (int, error) {
	return c.w.Write(bytes)
}

func (c *GzipWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *GzipWriter) Close() error {
	return c.zw.Close()
}
