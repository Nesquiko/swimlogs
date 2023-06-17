package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
	MaxBytes        = 1_048_576
)

// SwimLogsServer implements interface openapi.ServerInterface
// which is generated from OpenApi spec located in documentation.
type SwimLogsServer struct {
	app app.SwimLogsApp
}

func NewServerHandler(swimLogs app.SwimLogsApp, feOrigin string) http.Handler {
	s := SwimLogsServer{swimLogs}

	r := chi.NewRouter()
	serverOpts := openapi.ChiServerOptions{
		BaseRouter:  r,
		Middlewares: PublicMiddleware(feOrigin),
		BaseURL:     "",
	}

	// group for handling OPTIONS requests
	r.Group(func(r chi.Router) {
		r.Use(cors(feOrigin))
		r.Options("/*", nil)
	})

	r.Group(func(r chi.Router) {
		r.Get(serverOpts.BaseURL+"/monitoring/heartbeat", s.Heartbeat)
	})

	return openapi.HandlerWithOptions(&s, serverOpts)
}

func respondWithCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, code int, response any) {
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
	}
}

func readJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var dst T
	r.Body = http.MaxBytesReader(w, r.Body, int64(MaxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError

		invalidFieldPrefix := "json: unknown field "
		largeBodyErrorStr := "http: request body too large"

		switch {
		case errors.As(err, &syntaxErr):
			log.Debug().
				Err(err).
				Msgf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)
			return dst, fmt.Errorf(
				"body contains badly-formed JSON (at character %d)",
				syntaxErr.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			log.Debug().Err(err).Msg("body contains badly-formed JSON")
			return dst, errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				log.Debug().
					Err(err).
					Msgf("body contains incorrect JSON type for field %q", unmarshalTypeErr.Field)
				return dst, fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeErr.Field,
				)
			}
			log.Debug().
				Err(err).
				Msgf("body contains incorrect JSON type (at character %d)", unmarshalTypeErr.Offset)
			return dst, fmt.Errorf(
				"body contains incorrect JSON type (at character %d)",
				unmarshalTypeErr.Offset,
			)

		case errors.Is(err, io.EOF):
			log.Debug().Err(err).Msg("body must not be empty")
			return dst, errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), invalidFieldPrefix):
			fieldName := strings.TrimPrefix(err.Error(), invalidFieldPrefix)
			log.Debug().Msgf("body contains unknown key %s", fieldName)
			return dst, fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == largeBodyErrorStr:
			log.Debug().Err(err).Msg("body is too large")
			return dst, fmt.Errorf("body must not be larger than %d bytes", MaxBytes)

		case errors.As(err, &invalidUnmarshalErr):
			log.Error().Err(err).Msg("invalid unmarshal target")

		default:
			log.Error().Err(err).Msg("failed to decode request body")
			return dst, err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		log.Debug().Err(err).Msg("body must only contain a single JSON value")
		return dst, errors.New("body must only contain a single JSON value")
	}

	return dst, nil
}
