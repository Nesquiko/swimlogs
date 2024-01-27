package server

import (
	"net/http"
	"os"
	"time"

	chiMidleware "github.com/go-chi/chi/v5/middleware"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/Nesquiko/swimlogs/apidef"
)

// type StrictHTTPHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error)
// type StrictHTTPMiddlewareFunc func(f StrictHTTPHandlerFunc, operationID string) StrictHTTPHandlerFunc

func publicMiddleware(feOrigin string) []apidef.MiddlewareFunc {
	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout, FormatTimestamp: func(i interface{}) string { return time.Now().Format("2006-01-02 15:04:05.000") }})

	oas, err := apidef.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load OpenAPI spec")
	}
	oas.Servers = nil // removes validation of server, since we are using proxy

	// dont move things around, order matters, executes last to first
	return []apidef.MiddlewareFunc{
		nethttpmiddleware.OapiRequestValidator(oas),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Int("status", status).
				Int("size", size).
				// defaults is miliseconds, to override use zerolog.DurationFieldUnit
				Dur("duration", duration).
				Msg("")
		}),
		hlog.RequestHandler("req"),
		hlog.UserAgentHandler("user_agent"),
		hlog.CustomHeaderHandler("ip", "X-Real-Ip"),
		hlog.NewHandler(l),
		cors(feOrigin),
		chiMidleware.Recoverer,
	}
}

func cors(feOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// w.Header().Add("Vary", "Origin") not needed, because there is only one allowed origin
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Set("Access-Control-Allow-Origin", feOrigin)

			if r.Method == http.MethodOptions &&
				r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
