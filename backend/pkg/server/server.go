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
	app    app.SwimLogsApp
	logger *httplog.Logger
}

type ApiError struct {
	apidef.ErrorDetail
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("error %q, status %d", e.Title, e.Status)
}

func NewServer(app app.SwimLogsApp, logger *httplog.Logger, feOrigin string) http.Handler {
	r := chi.NewRouter()

	srv := SwimLogsServer{
		app:    app,
		logger: logger,
	}

	r.Use(topLevelMiddleware(feOrigin)...)
	r.Options("/*", nil)

	r.Group(func(r chi.Router) {
		r.Use(publicMiddleware(logger)...)

		r.Post("/trainings", handleInOut(srv.CreateTraining, srv.logger))
		r.Get(
			"/trainings/summaries",
			handleQueryOut(pageParamsExtractor, srv.SummariesPage, srv.logger),
		)
		r.Get("/trainings/summaries/current-week", handleOut(srv.SummariesCurrentWeek, srv.logger))
		r.Delete(
			"/trainings/{id}",
			handlePathStatus(pathIdExtractor, srv.DeleteTraining, srv.logger),
		)
		r.Get("/trainings/{id}", handlePathOut(pathIdExtractor, srv.TrainingById, srv.logger))
		r.Patch(
			"/trainings/{id}",
			handlePathInOut(pathIdExtractor, srv.EditTrainingSession, srv.logger),
		)
		r.Delete(
			"/trainings/{id}/sets/{setId}",
			handlePathStatus(pathIdAndSetIdExtractor, srv.DeleteSet, srv.logger),
		)
		r.Patch(
			"/trainings/{id}/sets/{setId}",
			handlePathInOut(pathIdAndSetIdExtractor, srv.EditSet, srv.logger),
		)
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

func handleOut[Out any](f OutFunc[Out], logger *httplog.Logger) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp EmptyType, in EmptyType) (Out, int, error) {
		return f(ctx)
	}
	return handle(tf, EmptyFunc, EmptyFunc, logger)
}

func handleInOut[In, Out any](f InOutFunc[In, Out], logger *httplog.Logger) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp EmptyType, in In) (Out, int, error) {
		return f(ctx, in)
	}
	return handle(tf, EmptyFunc, EmptyFunc, logger)
}

func handleQueryOut[QP, Out any](
	qpConv QueryParamsConv[QP],
	f InOutFunc[QP, Out],
	logger *httplog.Logger,
) http.HandlerFunc {
	tf := func(ctx context.Context, pp EmptyType, qp QP, in EmptyType) (Out, int, error) {
		return f(ctx, qp)
	}
	return handle(tf, EmptyFunc, qpConv, logger)
}

func handlePathStatus[PP any](
	ppConv PathParamsConv[PP],
	f PathStatusFunc[PP],
	logger *httplog.Logger,
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in EmptyType) (EmptyType, int, error) {
		status, err := f(ctx, pp)
		return EmptyType{}, status, err
	}
	return handle(tf, ppConv, EmptyFunc, logger)
}

func handlePathOut[PP, Out any](
	ppConv PathParamsConv[PP],
	f PathOutFunc[PP, Out],
	logger *httplog.Logger,
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in EmptyType) (Out, int, error) {
		return f(ctx, pp)
	}
	return handle(tf, ppConv, EmptyFunc, logger)
}

func handlePathInOut[PP, In, Out any](
	ppConv PathParamsConv[PP],
	f PathInOutFunc[PP, In, Out],
	logger *httplog.Logger,
) http.HandlerFunc {
	tf := func(ctx context.Context, pp PP, qp EmptyType, in In) (Out, int, error) {
		return f(ctx, pp, in)
	}
	return handle(tf, ppConv, EmptyFunc, logger)
}

func handle[PP, QP, In, Out any](
	f TargetFunc[PP, QP, In, Out],
	ppFunc PathParamsConv[PP],
	qpFunc QueryParamsConv[QP],
	logger *httplog.Logger,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := decode[In](w, r)
		if err != nil {
			apiErr := badRequest(err)
			encodeError(w, apiErr, logger)
			return
		}

		var apiErr *ApiError
		ppParams, err := ppFunc(r)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr, logger)
				return
			}
			logger.Error(
				UnexpectedError,
				slog.String("where", "path-params"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		qpParams, err := qpFunc(r)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr, logger)
				return
			}
			logger.Error(
				UnexpectedError,
				slog.String("where", "query-params"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		out, status, err := f(r.Context(), ppParams, qpParams, in)
		if err != nil {
			if errors.As(err, &apiErr) {
				encodeError(w, apiErr, logger)
				return
			}
			logger.Error(
				UnexpectedError,
				slog.String("where", "handler"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		encode(w, status, out, logger)
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

func pathIdAndSetIdExtractor(r *http.Request) (IdSetId, error) {
	id := chi.URLParam(r, "id")
	setId := chi.URLParam(r, "setId")
	uid, err := uuid.Parse(id)
	ids := IdSetId{}

	if err != nil {
		return IdSetId{}, invalidPathParam("id", id)
	}
	ids.id = uid

	setUid, err := uuid.Parse(setId)
	if err != nil {
		return IdSetId{}, invalidPathParam("setId", setId)
	}
	ids.setId = setUid

	return ids, nil
}

func invalidQueryParam(param, value string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:   "invalid.query.param",
			Detail: fmt.Sprintf("Invalid %q: %q", param, value),
			Status: http.StatusBadRequest,
			Title:  fmt.Sprintf("Invalid query param %q: %q", param, value),
		},
	}
}

func invalidPathParam(param, value string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Code:   "invalid.path.param",
			Detail: fmt.Sprintf("Invalid %q: %q", param, value),
			Status: http.StatusBadRequest,
			Title:  fmt.Sprintf("Invalid path param %q: %q", param, value),
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
