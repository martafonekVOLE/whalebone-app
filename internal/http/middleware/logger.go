package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Logger returns a middleware that logs HTTP requests using zap.
func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			requestID := middleware.GetReqID(r.Context())

			logger.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("request_id", requestID),
			)

			defer func() {
				logger.Info("request completed",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("remote_ip", r.RemoteAddr),
					zap.String("request_id", requestID),
					zap.Int("status", ww.Status()),
					zap.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
