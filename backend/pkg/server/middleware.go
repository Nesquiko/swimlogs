package server

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	validation_middleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/Nesquiko/swimlogs/pkg/api"
)

type OapiValidationOptions struct {
	spec         *openapi3.T
	errorHandler func(w http.ResponseWriter, message string, statusCode int)
}

func middleware(
	logger *httplog.Logger,
	feOrigin string,
	opts OapiValidationOptions,
) []api.MiddlewareFunc {
	return []api.MiddlewareFunc{
		chi_middleware.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{feOrigin},
			AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			MaxAge:         300,
		}),
		chi_middleware.RealIP,
		validation_middleware.OapiRequestValidatorWithOptions(
			opts.spec,
			&validation_middleware.Options{ErrorHandler: opts.errorHandler},
		),
		httplog.RequestLogger(logger),
		chi_middleware.AllowContentType(ApplicationJSON),
	}
}

func heartbeat() func(http.Handler) http.Handler {
	return chi_middleware.Heartbeat("/monitoring/heartbeat")
}
