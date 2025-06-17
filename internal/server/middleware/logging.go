package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture the status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		entry := log.WithFields(log.Fields{
			"path":        r.URL.Path,
			"method":      r.Method,
			"status_code": lrw.statusCode,
			"duration":    duration.Seconds(),
			"from":        r.RemoteAddr,
		})

		if rec := r.Context().Value(panicKey); rec != nil {
			entry = entry.WithField("panic", rec)
			entry.Error("Recovered from panic")

			return
		}

		if lrw.statusCode >= http.StatusInternalServerError {
			entry.Error("Server error")
			return
		}

		entry.Info("Processed a request")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
