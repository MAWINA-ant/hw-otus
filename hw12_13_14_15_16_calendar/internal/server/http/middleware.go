package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		latency := time.Since(start)
		ip := getClientIP(r)

		logLine := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			ip,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.String(),
			r.Proto,
			wrapped.statusCode,
			latency.Milliseconds(),
			r.UserAgent(),
		)
		logger.Info(logLine)
	})
}
