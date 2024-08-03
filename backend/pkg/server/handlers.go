package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

// (POST /trainings)
func (s *SwimLogsServer) CreateTraining(
	ctx context.Context,
	r apidef.CreateTrainingRequest,
) (apidef.CreateTrainingReponse, int, error) {
	td, err := s.app.CreateTraining(ctx, r)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid training", slog.String("error", validationErr.Error()))
			return apidef.TrainingSummary{}, apiErr.Status, apiErr
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "CreateTraining"),
		)
		apiErr := internalServerError()
		return apidef.TrainingSummary{}, apiErr.Status, apiErr
	}
	return td, http.StatusCreated, nil
}

// (GET /trainings/summaries)
func (s *SwimLogsServer) SummariesPage(
	ctx context.Context,
	params apidef.SummariesPageParams,
) (apidef.TrainingSummariesResponse, int, error) {
	if params.Page < 0 {
		slog.Warn(
			"invalid page query params",
			slog.Int("page", params.Page),
		)
		apiErr := invalidQueryParam("page", strconv.Itoa(params.Page))
		return apidef.TrainingSummariesResponse{}, apiErr.Status, apiErr
	} else if params.PageSize < 1 {
		slog.Warn(
			"invalid pageSize query params",
			slog.Int("pageSize", params.PageSize),
		)
		apiErr := invalidQueryParam("pageSize", strconv.Itoa(params.Page))
		return apidef.TrainingSummariesResponse{}, apiErr.Status, apiErr
	}

	summaries, total, err := s.app.TrainingSummariesPage(ctx, params.Page, params.PageSize)
	if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "SummariesPage"),
		)
		apiErr := internalServerError()
		return apidef.TrainingSummariesResponse{}, apiErr.Status, apiErr
	}

	pagination := apidef.Pagination{Page: params.Page, PageSize: len(summaries), Total: total}
	return apidef.TrainingSummariesResponse{
		Summaries:  summaries,
		Pagination: pagination,
	}, http.StatusOK, nil
}

// (GET /trainings/summaries/current-week)
func (s *SwimLogsServer) SummariesCurrentWeek(
	ctx context.Context,
) (apidef.TrainingSummariesCurrentWeekResponse, int, error) {
	summaries, err := s.app.TrainingSummariesCurrentWeek(ctx)
	if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "SummariesCurrentWeek"),
		)
		apiErr := internalServerError()
		return apidef.TrainingSummariesCurrentWeekResponse{}, apiErr.Status, apiErr
	}

	return apidef.TrainingSummariesCurrentWeekResponse{Summaries: summaries}, http.StatusOK, nil
}

// (DELETE /trainings/{id})
func (s *SwimLogsServer) DeleteTraining(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	err := s.app.DeleteTraining(ctx, id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"training not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		return http.StatusNotFound, notFound("training", id.String())
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "DeleteTraining"),
		)
		apiErr := internalServerError()
		return apiErr.Status, apiErr
	}

	return http.StatusNoContent, nil
}

// (GET /trainings/{id})
func (s *SwimLogsServer) TrainingById(
	ctx context.Context,
	id uuid.UUID,
) (apidef.Training, int, error) {
	t, err := s.app.Training(ctx, id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn("training not found",
			slog.String("error", NotFoundCode),
			slog.String("id", id.String()))
		return apidef.Training{}, http.StatusNotFound, notFound("training", id.String())
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "TrainingById"),
		)
		apiErr := internalServerError()
		return apidef.Training{}, apiErr.Status, apiErr
	}

	return t, http.StatusOK, nil
}

// (PATCH /trainings/{id})
func (s *SwimLogsServer) EditTrainingSession(
	ctx context.Context,
	id uuid.UUID,
	r apidef.EditSessionRequest,
) (apidef.EditSessionResponse, int, error) {
	td, err := s.app.EditTrainingSession(ctx, id, r)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid training session", slog.String("error", validationErr.Error()))
			return apidef.TrainingSummary{}, apiErr.Status, apiErr
		}

		if errors.Is(err, app.ErrNotFound) {
			slog.Warn("training not found",
				slog.String("error", NotFoundCode),
				slog.String("id", id.String()))
			return apidef.TrainingSummary{}, http.StatusNotFound, notFound("training", id.String())
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "EditTrainingSession"),
		)
		apiErr := internalServerError()
		return apidef.TrainingSummary{}, apiErr.Status, apiErr
	}
	return td, http.StatusOK, nil
}

// (DELETE /trainings/{id}/sets/{setId})
func (s *SwimLogsServer) DeleteSet(
	ctx context.Context,
	params IdSetId,
) (int, error) {
	err := s.app.DeleteSet(ctx, params.id, params.setId)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"set not found",
			slog.String("error", err.Error()),
			slog.String("trainingId", params.id.String()),
			slog.String("setId", params.setId.String()),
		)
		return http.StatusNotFound, notFound("set", params.setId.String())
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "DeleteSet"),
		)
		apiErr := internalServerError()
		return apiErr.Status, apiErr
	}

	return http.StatusNoContent, nil
}

// (PATCH /trainings/{id}/sets/{setId})
func (s *SwimLogsServer) EditSet(
	ctx context.Context,
	params IdSetId,
	r apidef.EditSetRequest,
) (apidef.EditSetResponse, int, error) {
	set, totalDistance, err := s.app.EditSet(ctx, params.id, params.setId, r)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid set", slog.String("error", validationErr.Error()))
			return apidef.EditSetResponse{}, apiErr.Status, apiErr
		}

		if errors.Is(err, app.ErrNotFound) {
			slog.Warn(
				"set not found",
				slog.String("error", err.Error()),
				slog.String("trainingId", params.id.String()),
				slog.String("setId", params.setId.String()),
			)
			return apidef.EditSetResponse{}, http.StatusNotFound, notFound(
				"set",
				params.setId.String(),
			)
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "EditSet"),
		)
		apiErr := internalServerError()
		return apidef.EditSetResponse{}, apiErr.Status, apiErr
	}
	return apidef.EditSetResponse{
		Set:           set,
		TotalDistance: totalDistance,
	}, http.StatusOK, nil
}

func (s *SwimLogsServer) MoveSet(
	ctx context.Context,
	params IdSetId,
	r apidef.MoveSetRequest,
) (apidef.Training, int, error) {
	t, err := s.app.MoveSet(ctx, params.id, params.setId, r.NewSetOrder)

	if validationErr, ok := err.(*app.ValidationError); ok {
		apiErr := fromValidationError(validationErr)
		slog.Warn("invalid new set order", slog.String("error", validationErr.Error()))
		return apidef.Training{}, apiErr.Status, apiErr
	} else if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"set not found",
			slog.String("error", err.Error()),
			slog.String("trainingId", params.id.String()),
			slog.String("setId", params.setId.String()),
		)
		return apidef.Training{}, http.StatusNotFound, notFound("set", params.setId.String())
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "MoveSet"),
		)
		apiErr := internalServerError()
		return apidef.Training{}, apiErr.Status, apiErr
	}

	return t, http.StatusOK, nil
}
