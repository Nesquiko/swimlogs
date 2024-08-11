package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/apidef"
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
	apidef.ErrorDetail
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("error %q, status %d", e.Title, e.Status)
}

func fromValidationError(e *app.ValidationError) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:                 e.Code,
			Title:                e.Title,
			Detail:               e.Detail,
			Status:               e.Status,
			AdditionalProperties: e.AdditionalProperties,
		},
	}
}

func NewServer(
	app app.SwimLogsApp,
	middlewareLogger *httplog.Logger,
	feOrigin string,
) http.Handler {
	r := chi.NewRouter()

	srv := SwimLogsServer{
		app: app,
	}

	r.Use(topLevelMiddleware(feOrigin)...)
	r.Options("/*", nil)

	r.Group(func(r chi.Router) {
		r.Use(publicMiddleware(middlewareLogger)...)

		r.Post("/trainings", handleInOut(srv.CreateTraining))
		r.Get("/trainings/summaries", handleQueryOut(pageParamsExtractor, srv.SummariesPage))
		r.Delete("/trainings/{id}", handlePathStatus(pathIdExtractor, srv.DeleteTraining))
		r.Get("/trainings/{id}", handlePathOut(pathIdExtractor, srv.TrainingById))
		r.Patch("/trainings/{id}", handlePathInOut(pathIdExtractor, srv.EditTrainingSession))
		r.Delete("/sets/{id}", handlePathStatus(pathIdExtractor, srv.DeleteSet))
		r.Patch("/sets/{id}", handlePathInOut(pathIdExtractor, srv.EditSet))
		r.Patch("/sets/{id}/move", handlePathInOut(pathIdExtractor, srv.MoveSet))
	})

	return r
}

type (
	EmptyType struct{}

	PathParamsConv[PP any]  func(r *http.Request) (PP, error)
	QueryParamsConv[QP any] func(r *http.Request) (QP, error)

	OutFunc[Out any]          func(context.Context) (Out, int, error)
	InOutFunc[In, Out any]    func(context.Context, In) (Out, int, error)
	QueryOutFunc[QP, Out any] func(context.Context, QP) (Out, int, error)

	PathStatusFunc[PP any]         func(context.Context, PP) (int, error)
	PathOutFunc[PP, Out any]       func(context.Context, PP) (Out, int, error)
	PathInOutFunc[PP, In, Out any] func(context.Context, PP, In) (Out, int, error)

	TargetFunc[PP, QP, In, Out any] func(context.Context, PP, QP, In) (Out, int, error)
)

var EmptyFunc = func(r *http.Request) (EmptyType, error) { return EmptyType{}, nil }

func handleOut[Out any](f OutFunc[Out]) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp EmptyType, in EmptyType) (Out, int, error) {
		return f(ctx)
	}
	return handle(tf, EmptyFunc, EmptyFunc)
}

func handleInOut[In, Out any](f InOutFunc[In, Out]) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp EmptyType, in In) (Out, int, error) {
		return f(ctx, in)
	}
	return handle(tf, EmptyFunc, EmptyFunc)
}

func handleQueryOut[QP, Out any](
	qpConv QueryParamsConv[QP],
	f InOutFunc[QP, Out],
) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp QP, in EmptyType) (Out, int, error) {
		return f(ctx, qp)
	}
	return handle(tf, EmptyFunc, qpConv)
}

func handlePathStatus[PP any](
	ppConv PathParamsConv[PP],
	f PathStatusFunc[PP],
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in EmptyType) (EmptyType, int, error) {
		status, err := f(ctx, pp)
		return EmptyType{}, status, err
	}
	return handle(tf, ppConv, EmptyFunc)
}

func handlePathOut[PP, Out any](
	ppConv PathParamsConv[PP],
	f PathOutFunc[PP, Out],
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in EmptyType) (Out, int, error) {
		return f(ctx, pp)
	}
	return handle(tf, ppConv, EmptyFunc)
}

func handlePathInOut[PP, In, Out any](
	ppConv PathParamsConv[PP],
	f PathInOutFunc[PP, In, Out],
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in In) (Out, int, error) {
		return f(ctx, pp, in)
	}
	return handle(tf, ppConv, EmptyFunc)
}

func handle[PP, QP, In, Out any](
	f TargetFunc[PP, QP, In, Out],
	ppFunc PathParamsConv[PP],
	qpFunc QueryParamsConv[QP],
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in In
		if _, ok := any(in).(EmptyType); !ok {
			decodedIn, err := decode[In](w, r)
			if err != nil {
				apiErr := badRequest(err)
				encodeError(w, apiErr)
				return
			}
			in = decodedIn
		}

		var apiErr *ApiError
		ppParams, err := ppFunc(r)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr)
				return
			}
			slog.Error(
				UnexpectedError,
				slog.String("error", err.Error()),
				slog.String("where", "path-params"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		qpParams, err := qpFunc(r)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr)
				return
			}
			slog.Error(
				UnexpectedError,
				slog.String("error", err.Error()),
				slog.String("where", "query-params"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		out, status, err := f(r.Context(), ppParams, qpParams, in)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr)
				return
			}
			slog.Error(
				UnexpectedError,
				slog.String("error", err.Error()),
				slog.String("where", "handler"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		if _, ok := any(out).(EmptyType); ok {
			w.WriteHeader(status)
			return
		}
		encode(w, status, out)
	})
}

func pageParamsExtractor(r *http.Request) (apidef.SummariesPageParams, error) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return apidef.SummariesPageParams{}, invalidQueryParam(
			"page",
			r.URL.Query().Get("page"),
		)
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		return apidef.SummariesPageParams{}, invalidQueryParam(
			"pageSize",
			r.URL.Query().Get("pageSize"),
		)
	}

	return apidef.SummariesPageParams{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func pathIdExtractor(r *http.Request) (uuid.UUID, error) {
	id := chi.URLParam(r, "id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, invalidPathParam("id", id)
	}
	return uid, nil
}

type IdSetId struct {
	id    uuid.UUID
	setId uuid.UUID
}

const (
	InvalidQueryParamCode         = "invalid.query.param"
	InvalidQueryParamTitleFormat  = "Invalid %q: %q"
	InvalidQueryParamDetailFormat = "Invalid query param %q: %q"
)

func invalidQueryParam(param, value string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:   InvalidQueryParamCode,
			Title:  fmt.Sprintf(InvalidQueryParamTitleFormat, param, value),
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf(InvalidQueryParamDetailFormat, param, value),
		},
	}
}

const (
	InvalidPathParamCode         = "invalid.path.param"
	InvalidPathParamTitleFormat  = "Invalid %q: %q"
	InvalidPathParamDetailFormat = "Invalid path param %q: %q"
)

func invalidPathParam(param, value string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:   InvalidPathParamCode,
			Title:  fmt.Sprintf(InvalidPathParamTitleFormat, param, value),
			Status: http.StatusBadRequest,
			Detail: fmt.Sprintf(InvalidPathParamDetailFormat, param, value),
		},
	}
}

func badRequest(err error) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:   "invalid.request",
			Title:  "Bad request",
			Detail: fmt.Sprintf("Request was invalid due to %q", err.Error()),
			Status: http.StatusBadRequest,
		},
	}
}

func internalServerError() *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
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
		ErrorDetail: apidef.ErrorDetail{
			Code:   NotFoundCode,
			Title:  fmt.Sprintf(NotFoundTitleFormat, resoure),
			Detail: fmt.Sprintf(NotFoundDetailFormat, resoure, id),
			Status: http.StatusNotFound,
		},
	}
}
