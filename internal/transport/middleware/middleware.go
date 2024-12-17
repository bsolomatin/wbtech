package middleware

import (
	"net/http"
	"time"
	"wbzerolevel/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info(r.Context(), "Request received", zap.String("Method", r.Method), zap.String("URL", r.URL.Path))
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info(r.Context(), "Request processed", zap.String("Method", r.Method), zap.String("URL", r.URL.Path), zap.String("Processed time", duration.String()))
		})
	}
}