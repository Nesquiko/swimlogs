package server

import (
	"net/http"
	"os"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func Middleware(feOrigin string) []openapi.MiddlewareFunc {
	l := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(
		zerolog.ConsoleWriter{
			Out:             os.Stderr,
			FormatTimestamp: func(i interface{}) string { return time.Now().Format("2006-01-02 15:04:05.000") },
		},
	)

	oas, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load OpenAPI spec")
	}

	// dont move things around, order matters, executes last to first
	return []openapi.MiddlewareFunc{
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
		hlog.RemoteAddrHandler("ip"),
		hlog.NewHandler(l),
		middleware.OapiRequestValidator(oas),
		cors(feOrigin),
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
