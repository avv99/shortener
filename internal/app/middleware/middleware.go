package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	*gzip.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(p []byte) (int, error) {
	return grw.Writer.Write(p)
}

func (grw gzipResponseWriter) Header() http.Header {
	return grw.ResponseWriter.Header()
}

func (grw gzipResponseWriter) WriteHeader(statusCode int) {
	grw.ResponseWriter.WriteHeader(statusCode)
}

func GZipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			r.Body, _ = gzip.NewReader(r.Body)
		}

		next.ServeHTTP(w, r)
	})
}
