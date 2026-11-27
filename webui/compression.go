package webui

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer      *gzip.Writer
	compress    bool
	wroteHeader bool
	closed      bool
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	contentType := strings.ToLower(w.Header().Get("Content-Type"))
	if (strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/")) &&
		w.Header().Get("Content-Encoding") == "" {
		w.compress = true
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Del("Content-Length")
		w.writer = gzip.NewWriter(w.ResponseWriter)
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.compress {
		return w.writer.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

func (w *gzipResponseWriter) Close() error {
	if w.writer != nil && !w.closed {
		w.closed = true
		return w.writer.Close()
	}
	return nil
}

// withAPICompression compresses JSON/text API responses when the client
// advertises gzip. Static files are handled by staticHandler.
func withAPICompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead || !strings.HasPrefix(r.URL.Path, "/api/") ||
			!strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") ||
			strings.HasPrefix(r.URL.Path, "/api/feeds/poster") ||
			strings.HasPrefix(r.URL.Path, "/api/feeds/cover") ||
			strings.HasPrefix(r.URL.Path, "/api/theme/background") {
			next.ServeHTTP(w, r)
			return
		}
		gw := &gzipResponseWriter{ResponseWriter: w}
		defer gw.Close()
		next.ServeHTTP(gw, r)
	})
}
