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
	validationErr := validateNewTraining(r)
	if validationErr != nil {
		slog.Warn("invalid training", slog.String("error", validationErr.Error()))
		return apidef.TrainingSummary{}, validationErr.Status, validationErr
	}

	td, err := s.app.CreateTraining(r)
	if err != nil {
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

	summaries, total, err := s.app.TrainingSummariesPage(params.Page, params.PageSize)
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

// TODO context everywhere!

// (GET /trainings/summaries/current-week)
func (s *SwimLogsServer) SummariesCurrentWeek(
	ctx context.Context,
) (apidef.TrainingSummariesCurrentWeekResponse, int, error) {
	panic("not implemented")
}

// (DELETE /trainings/{id})
func (s *SwimLogsServer) DeleteTraining(
	ctx context.Context,
	id uuid.UUID,
) (int, error) {
	err := s.app.DeleteTraining(id)

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
	t, err := s.app.Training(id)

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
	panic("not implemented")
}

// (DELETE /trainings/{id}/sets/{setId})
func (s *SwimLogsServer) DeleteSet(
	ctx context.Context,
	params IdSetId,
) (int, error) {
	panic("not implemented")
}

// (PATCH /trainings/{id}/sets/{setId})
func (s *SwimLogsServer) EditSet(
	ctx context.Context,
	params IdSetId,
	r apidef.EditSetRequest,
) (apidef.EditSetResponse, int, error) {
	panic("not implemented")
}

// // (GET /trainings/details)
// func (s *SwimLogsServer) TrainingDetails(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	params apidef.TrainingDetailsParams,
// ) {
// 	if params.Page < 0 || params.PageSize < 1 {
// 		log.Warn().
// 			Int("page", params.Page).
// 			Int("pageSize", params.PageSize).
// 			Msg("invalid query params")
// 		respondWithCode(w, http.StatusBadRequest)
// 		return
// 	}
//
// 	details, total, err := s.app.TrainingDetailsPage(params.Page, params.PageSize)
// 	if err != nil {
// 		log.Warn().Err(err).Int("page", params.Page).Int("pageSize", params.PageSize)
// 		respondWithCode(w, http.StatusInternalServerError)
// 		return
// 	}
//
// 	pagination := apidef.Pagination{Page: params.Page, PageSize: len(details), Total: total}
// 	response := apidef.TrainingDetailsResponse{Details: details, Pagination: pagination}
// 	respondWithJSON(w, http.StatusOK, response)
// }
//
// // (GET /trainings/details/current-week)
// func (s *SwimLogsServer) TrainingDetailsCurrentWeek(w http.ResponseWriter, r *http.Request) {
// 	details, err := s.app.TrainingDetailsCurrentWeek()
// 	if err != nil {
// 		log.Error().Err(err).Msg("internal server error")
// 		respondWithCode(w, http.StatusInternalServerError)
// 		return
// 	}
//
// 	response := apidef.TrainingDetailsCurrentWeekResponse{
// 		Details: details,
// 	}
// 	respondWithJSON(w, http.StatusOK, response)
// }
//
// // (DELETE /trainings/{id})
// func (s *SwimLogsServer) DeleteTraining(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	id types.UUID,
// ) {
// 	err := s.app.DeleteTraining(id)
// 	if errors.Is(err, app.ErrNotFound) {
// 		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
// 		respondWithCode(w, http.StatusNotFound)
// 		return
// 	} else if err != nil {
// 		log.Error().Err(err).Msg("internal server error")
// 		respondWithCode(w, http.StatusInternalServerError)
// 		return
// 	}
//
// 	respondWithCode(w, http.StatusNoContent)
// }
//
// // (GET /trainings/{id})
// func (s *SwimLogsServer) Training(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	id types.UUID,
// ) {
// 	t, err := s.app.Training(id)
// 	if errors.Is(err, app.ErrNotFound) {
// 		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
// 		respondWithCode(w, http.StatusNotFound)
// 		return
// 	} else if err != nil {
// 		log.Error().Err(err).Msg("internal server error")
// 		respondWithCode(w, http.StatusInternalServerError)
// 		return
// 	}
//
// 	response := apidef.TrainingResponse(t)
// 	respondWithJSON(w, http.StatusOK, response)
// }
//
// // (PUT /trainings/{id})
// func (s *SwimLogsServer) EditTraining(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	id types.UUID,
// ) {
// 	req, err := readJSON[apidef.EditTrainingRequest](w, r)
// 	if err != nil {
// 		log.Warn().Err(err).Msg("failed to read request body")
// 		respondWithCode(w, http.StatusBadRequest)
// 		return
// 	}
//
// 	td, err := s.app.EditTraining(id, req)
// 	if errors.Is(err, app.ErrNotFound) {
// 		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
// 		respondWithCode(w, http.StatusNotFound)
// 		return
// 	} else if err != nil {
// 		log.Warn().Err(err).Msg("invalid training")
// 		respondWithCode(w, http.StatusBadRequest)
// 		return
// 	}
//
// 	response := apidef.EditTrainingResponse(td)
// 	respondWithJSON(w, http.StatusOK, response)
// }
