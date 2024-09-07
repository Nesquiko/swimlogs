package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

func (s SwimLogsServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	request, decodeErr := Decode[api.NewTraining](w, r)
	if decodeErr != nil {
		encodeError(w, decodeErr)
		return
	}

	t, err := s.app.CreateTraining(r.Context(), request)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid training", slog.String("error", validationErr.Error()))
			encodeError(w, apiErr)
			return
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "CreateTraining"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusCreated, t)
}

func (s SwimLogsServer) TrainingById(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	t, err := s.app.TrainingById(r.Context(), id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn("training not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()))
		encodeError(w, notFound("training", id.String()))
		return
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "TrainingById"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusOK, t)
}

func (s SwimLogsServer) DeleteSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) DeleteTraining(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) EditSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) EditTrainingSession(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) MoveSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) ReplaceSetComponents(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}

func (s SwimLogsServer) SummariesPage(
	w http.ResponseWriter,
	r *http.Request,
	params api.SummariesPageParams,
) {
	panic("unimplemented")
}

// // (GET /trainings/summaries)
// func (s *SwimLogsServer) SummariesPage(
// 	ctx context.Context,
// 	params api.SummariesPageParams,
// ) (api.TrainingSummariesResponse, int, error) {
// 	if params.Page < 0 {
// 		slog.Warn(
// 			"invalid page query params",
// 			slog.Int("page", params.Page),
// 		)
// 		apiErr := invalidQueryParam("page", strconv.Itoa(params.Page))
// 		return api.TrainingSummariesResponse{}, apiErr.Status, apiErr
// 	} else if params.PageSize < 1 {
// 		slog.Warn(
// 			"invalid pageSize query params",
// 			slog.Int("pageSize", params.PageSize),
// 		)
// 		apiErr := invalidQueryParam("pageSize", strconv.Itoa(params.Page))
// 		return api.TrainingSummariesResponse{}, apiErr.Status, apiErr
// 	}
//
// 	summaries, total, err := s.app.TrainingSummariesPage(ctx, params)
// 	if err != nil {
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "SummariesPage"),
// 		)
// 		apiErr := internalServerError()
// 		return api.TrainingSummariesResponse{}, apiErr.Status, apiErr
// 	}
//
// 	pagination := api.Pagination{Page: params.Page, PageSize: len(summaries), Total: total}
// 	return api.TrainingSummariesResponse{
// 		Summaries:  summaries,
// 		Pagination: pagination,
// 	}, http.StatusOK, nil
// }
//
// // (DELETE /trainings/{id})
// func (s *SwimLogsServer) DeleteTraining(
// 	ctx context.Context,
// 	id uuid.UUID,
// ) (int, error) {
// 	err := s.app.DeleteTraining(ctx, id)
//
// 	if errors.Is(err, app.ErrNotFound) {
// 		slog.Warn(
// 			"training not found",
// 			slog.String("error", err.Error()),
// 			slog.String("id", id.String()),
// 		)
// 		return http.StatusNotFound, notFoundId("training", id)
// 	} else if err != nil {
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "DeleteTraining"),
// 		)
// 		apiErr := internalServerError()
// 		return apiErr.Status, apiErr
// 	}
//
// 	return http.StatusNoContent, nil
// }
//
// // (PATCH /trainings/{id})
// func (s *SwimLogsServer) EditTrainingSession(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	r api.EditSessionRequest,
// ) (api.EditSessionResponse, int, error) {
// 	td, err := s.app.EditTrainingSession(ctx, id, r)
// 	if err != nil {
// 		if validationErr, ok := err.(*app.ValidationError); ok {
// 			apiErr := fromValidationError(validationErr)
// 			slog.Warn("invalid training session", slog.String("error", validationErr.Error()))
// 			return api.TrainingSummary{}, apiErr.Status, apiErr
// 		}
//
// 		if errors.Is(err, app.ErrNotFound) {
// 			slog.Warn("training not found",
// 				slog.String("error", NotFoundCode),
// 				slog.String("id", id.String()))
// 			return api.TrainingSummary{}, http.StatusNotFound, notFoundId("training", id)
// 		}
//
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "EditTrainingSession"),
// 		)
// 		apiErr := internalServerError()
// 		return api.TrainingSummary{}, apiErr.Status, apiErr
// 	}
// 	return td, http.StatusOK, nil
// }
//
// // (DELETE /sets/{id})
// func (s *SwimLogsServer) DeleteSet(ctx context.Context, id uuid.UUID) (int, error) {
// 	err := s.app.DeleteSet(ctx, id)
//
// 	if errors.Is(err, app.ErrNotFound) {
// 		slog.Warn(
// 			"set not found",
// 			slog.String("error", err.Error()),
// 			slog.String("id", id.String()),
// 		)
// 		return http.StatusNotFound, notFoundId("set", id)
// 	} else if err != nil {
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "DeleteSet"),
// 		)
// 		apiErr := internalServerError()
// 		return apiErr.Status, apiErr
// 	}
//
// 	return http.StatusNoContent, nil
// }
//
// // (PATCH /sets/{id})
// func (s *SwimLogsServer) EditSet(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	r api.EditSetRequest,
// ) (api.EditSetResponse, int, error) {
// 	set, totalDistance, err := s.app.EditSet(ctx, id, r)
// 	if err != nil {
// 		if validationErr, ok := err.(*app.ValidationError); ok {
// 			apiErr := fromValidationError(validationErr)
// 			slog.Warn("invalid set", slog.String("error", validationErr.Error()))
// 			return api.EditSetResponse{}, apiErr.Status, apiErr
// 		}
//
// 		if errors.Is(err, app.ErrNotFound) {
// 			slog.Warn(
// 				"set not found",
// 				slog.String("error", err.Error()),
// 				slog.String("id", id.String()),
// 			)
// 			return api.EditSetResponse{}, http.StatusNotFound, notFoundId("set", id)
// 		}
//
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "EditSet"),
// 		)
// 		apiErr := internalServerError()
// 		return api.EditSetResponse{}, apiErr.Status, apiErr
// 	}
// 	return api.EditSetResponse{
// 		Set:           set,
// 		TotalDistance: totalDistance,
// 	}, http.StatusOK, nil
// }
//
// func (s *SwimLogsServer) MoveSet(
// 	ctx context.Context,
// 	id uuid.UUID,
// 	r api.MoveSetRequest,
// ) (api.Training, int, error) {
// 	t, err := s.app.MoveSet(ctx, id, r.NewSetOrder)
//
// 	if validationErr, ok := err.(*app.ValidationError); ok {
// 		apiErr := fromValidationError(validationErr)
// 		slog.Warn("invalid new set order", slog.String("error", validationErr.Error()))
// 		return api.Training{}, apiErr.Status, apiErr
// 	} else if errors.Is(err, app.ErrNotFound) {
// 		slog.Warn(
// 			"set not found",
// 			slog.String("error", err.Error()),
// 			slog.String("id", id.String()),
// 		)
// 		return api.Training{}, http.StatusNotFound, notFoundId("set", id)
// 	} else if err != nil {
// 		slog.Error(
// 			UnexpectedError,
// 			slog.String("error", err.Error()),
// 			slog.String("where", "MoveSet"),
// 		)
// 		apiErr := internalServerError()
// 		return api.Training{}, apiErr.Status, apiErr
// 	}
//
// 	return t, http.StatusOK, nil
// }
