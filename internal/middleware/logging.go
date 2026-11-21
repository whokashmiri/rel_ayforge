package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	if rec.statusCode != 0 {
		return
	}
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *responseRecorder) Unwrap() http.ResponseWriter {
	return rec.ResponseWriter
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := &responseRecorder{
				ResponseWriter: w,
			}
			next.ServeHTTP(recorder, r)
			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
			duration := time.Since(start)
			logger.InfoContext(
				r.Context(),
				"http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", statusCode,
				"duration_ms", float64(duration.Microseconds())/1000,
				"request_id", GetRequestID(r.Context()),
			)
		})

	}
}
