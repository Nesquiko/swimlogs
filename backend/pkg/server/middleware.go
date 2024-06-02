package server

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func topLevelMiddleware(feOrigin string) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Heartbeat("/monitoring/heartbeat"),
		cors.Handler(cors.Options{
			AllowedOrigins: []string{feOrigin},
			AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			MaxAge:         300,
		}),
	}
}

func publicMiddleware(logger *httplog.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.RealIP,
		httplog.RequestLogger(logger),
		middleware.AllowContentType(ApplicationJSON),
		middleware.Recoverer,
	}
}
