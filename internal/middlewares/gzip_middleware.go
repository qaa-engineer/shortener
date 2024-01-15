package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/qaa-engineer/shortener/internal/readers"
	"github.com/qaa-engineer/shortener/internal/writers"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriter := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			gzipWriter := writers.NewGzipWriter(w)
			responseWriter = gzipWriter
			defer func(gzipWriter *writers.GzipWriter) {
				err := gzipWriter.Close()
				if err != nil {
					log.Println(err)
				}
			}(gzipWriter)
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			compressReader, err := readers.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = compressReader
			defer func(compressReader *readers.GzipReader) {
				err := compressReader.Close()
				if err != nil {
					log.Println(err)
				}
			}(compressReader)
		}

		next.ServeHTTP(responseWriter, r)
	})
}
