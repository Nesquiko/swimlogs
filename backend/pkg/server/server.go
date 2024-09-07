package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

const (
	ContentType            = "Content-Type"
	ApplicationJSON        = "application/json"
	ApplicationProblemJSON = "application/problem+json"
	MaxBytes               = 1_048_576

	EncodingError   = "unexptected encoding error"
	UnexpectedError = "unexptected error"
)

type SwimLogsServer struct {
	app app.SwimLogsApp
}

type ApiError struct {
	api.ErrorDetail
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("error %q, status %d", e.Title, e.Status)
}

func NewServer(
	app app.SwimLogsApp,
	spec *openapi3.T,
	middlewareLogger *httplog.Logger,
	feOrigin string,
) http.Handler {
	r := chi.NewMux()
	r.Use(heartbeat())
	srv := SwimLogsServer{app: app}

	// When this PR https://github.com/oapi-codegen/oapi-codegen/pull/1608 is
	// merged, I will reconsider strict
	// strictHandler := api.NewStrictHandlerWithOptions(srv, nil, api.StrictHTTPServerOptions{
	// 	RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
	// 		// RequestErrorHandlerFunc should not needed, there is validation done with OpenApi spec
	// 		slog.Error(
	// 			"unexpected error handling in RequestErrorHandlerFunc",
	// 			slog.String("error", err.Error()),
	// 		)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	},
	// 	ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
	// 		var apiErr *ApiError
	// 		if errors.As(err, &apiErr) {
	// 			encodeError(w, apiErr)
	// 			return
	// 		}
	// 		slog.Error(
	// 			UnexpectedError,
	// 			slog.String("error", err.Error()),
	// 			slog.String("where", "handler"),
	// 		)
	// 		http.Error(w, "internal server error", http.StatusInternalServerError)
	// 	},
	// })

	validationOpts := OapiValidationOptions{
		spec:         spec,
		errorHandler: validationErrorHandler,
	}

	return api.HandlerWithOptions(srv, api.ChiServerOptions{
		BaseURL:     "",
		BaseRouter:  r,
		Middlewares: middleware(middlewareLogger, feOrigin, validationOpts),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			var invalidParamErr *api.InvalidParamFormatError
			var requiredParamError *api.RequiredParamError

			switch {
			case errors.As(err, &invalidParamErr):
				slog.Warn(
					"invalid path param",
					slog.String("error", err.Error()),
					slog.String("where", "ErrorHandlerFunc"),
				)
				encodeError(w, fromInvalidParamErr(invalidParamErr, chi.URLParam(r, "id")))
			case errors.As(err, &requiredParamError):
				slog.Warn(
					"missing required path param",
					slog.String("error", err.Error()),
					slog.String("where", "ErrorHandlerFunc"),
				)
				encodeError(w, fromRequiredParamErr(requiredParamError))
			default:
				slog.Error(
					"unexpected error handling in ErrorHandlerFunc",
					slog.String("error", err.Error()),
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		},
	})
}

func validationErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	slog.Warn("validationErrorHandler", "message", message)
	split := strings.Split(message, ": ")
	hashIndex := strings.Index(split[1], "#")

	var schema *string = nil
	if hashIndex != -1 {
		s := split[1][hashIndex:]
		schema = &s
	}

	path := extractSchemaPath(split[2 : len(split)-1])
	reason := split[len(split)-1]

	encodeError(w, validationError(path, reason, statusCode, schema))
}

func extractSchemaPath(subPaths []string) string {
	path := ""
	for _, subPath := range subPaths {
		subPath = strings.TrimPrefix(subPath, "Error at")
		subPath = strings.TrimSpace(subPath)
		subPath = strings.Trim(subPath, "\"")
		path += subPath
	}
	return path
}

const (
	SchemaValidationErrorCode = "invalid.request.schema"
	ValidationErrorCode       = "invalid.request"
	ValidationErrorTitle      = "Request doesn't comply with schema"
	ValidationErrorDetail     = "Validation failed, %s"
)

func validationError(path string, reason string, statusCode int, schema *string) *ApiError {
	additionalProperties := make(map[string]any)
	if path != "" {
		additionalProperties["path"] = path
	}
	if reason != "" {
		additionalProperties["reason"] = reason
	}
	if schema != nil {
		additionalProperties["schema"] = schema
	}
	if len(additionalProperties) == 0 {
		additionalProperties = nil
	}

	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:                 SchemaValidationErrorCode,
			Detail:               fmt.Sprintf(ValidationErrorDetail, reason),
			Status:               statusCode,
			Title:                ValidationErrorTitle,
			AdditionalProperties: additionalProperties,
		},
	}
}

func internalServerError() *ApiError {
	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:   "internal.server.error",
			Title:  "Internal Server Error",
			Detail: "Unexpected error on server",
			Status: http.StatusInternalServerError,
		},
	}
}

const (
	NotFoundCode         = "not.found"
	NotFoundTitleFormat  = "%s was not found"
	NotFoundDetailFormat = "%s with id '%s' was not found"
)

func notFoundId(resoure string, id uuid.UUID) *ApiError {
	return notFound(resoure, id.String())
}

func notFound(resoure, id string) *ApiError {
	if resoure == "" {
		resoure = "Resource"
	}

	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:   NotFoundCode,
			Title:  fmt.Sprintf(NotFoundTitleFormat, resoure),
			Detail: fmt.Sprintf(NotFoundDetailFormat, resoure, id),
			Status: http.StatusNotFound,
		},
	}
}

func fromValidationError(e *app.ValidationError) *ApiError {
	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:                 e.Code,
			Title:                e.Title,
			Detail:               e.Detail,
			Status:               e.Status,
			AdditionalProperties: e.AdditionalProperties,
		},
	}
}

const (
	InvalidParamErrorCode   = "invalid.path.param"
	InvalidParamErrorTitle  = "Invalid path param"
	InvalidParamErrorDetail = "Invalid path param %q: %q"
)

func fromInvalidParamErr(e *api.InvalidParamFormatError, invalidParam string) *ApiError {
	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:   InvalidParamErrorCode,
			Title:  InvalidParamErrorTitle,
			Detail: fmt.Sprintf(InvalidParamErrorDetail, e.ParamName, invalidParam),
			Status: http.StatusBadRequest,
		},
	}
}

const (
	RequiredParamErrorCode   = "required.path.param"
	RequiredParamErrorTitle  = "Required path param"
	RequiredParamErrorDetail = "Required path param %q is missing"
)

func fromRequiredParamErr(e *api.RequiredParamError) *ApiError {
	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:   RequiredParamErrorCode,
			Title:  RequiredParamErrorTitle,
			Detail: fmt.Sprintf(RequiredParamErrorDetail, e.ParamName),
			Status: http.StatusBadRequest,
		},
	}
}
