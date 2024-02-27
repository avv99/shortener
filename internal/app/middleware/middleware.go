package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Accept-Encoding клиента на поддержку gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Устанавливаем заголовок Content-Encoding в gzip
			w.Header().Set("Content-Encoding", "gzip")

			// Создаем новый gzip.Writer, который будет записывать данные в w
			gz := gzip.NewWriter(w)
			defer gz.Close()

			// Заменяем w на gz
			w = &gzipResponseWriter{Writer: gz, ResponseWriter: w}
		}

		// Продолжаем выполнение цепочки обработчиков с обновленным w
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	// Записываем сжатые данные в gzip.Writer
	return w.Writer.Write(b)
}
