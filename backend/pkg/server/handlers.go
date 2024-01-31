package server

import (
	"errors"
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

// (POST /trainings)
func (s *SwimLogsServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[apidef.CreateTrainingRequest](w, r)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	td, err := s.app.CreateTraining(req)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request")
		respondWithCode(w, http.StatusBadRequest)
		return
	}
	respondWithJSON(w, http.StatusCreated, td)
}

// (GET /trainings/details)
func (s *SwimLogsServer) TrainingDetails(
	w http.ResponseWriter,
	r *http.Request,
	params apidef.TrainingDetailsParams,
) {
	if params.Page < 0 || params.PageSize < 1 {
		log.Warn().
			Int("page", params.Page).
			Int("pageSize", params.PageSize).
			Msg("invalid query params")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	details, total, err := s.app.TrainingDetailsPage(params.Page, params.PageSize)
	if err != nil {
		log.Warn().Err(err).Int("page", params.Page).Int("pageSize", params.PageSize)
		respondWithCode(w, http.StatusInternalServerError)
		return
	}

	pagination := apidef.Pagination{Page: params.Page, PageSize: len(details), Total: total}
	response := apidef.TrainingDetailsResponse{Details: details, Pagination: pagination}
	respondWithJSON(w, http.StatusOK, response)
}

// (GET /trainings/details/current-week)
func (s *SwimLogsServer) TrainingDetailsCurrentWeek(w http.ResponseWriter, r *http.Request) {
	details, err := s.app.TrainingDetailsCurrentWeek()
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		respondWithCode(w, http.StatusInternalServerError)
		return
	}

	response := apidef.TrainingDetailsCurrentWeekResponse{
		Details: details,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// (DELETE /trainings/{id})
func (s *SwimLogsServer) DeleteTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	err := s.app.DeleteTraining(id)
	if errors.Is(err, app.ErrNotFound) {
		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
		respondWithCode(w, http.StatusNotFound)
		return
	} else if err != nil {
		log.Error().Err(err).Msg("internal server error")
		respondWithCode(w, http.StatusInternalServerError)
		return
	}

	respondWithCode(w, http.StatusNoContent)
}

// (GET /trainings/{id})
func (s *SwimLogsServer) Training(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	t, err := s.app.Training(id)
	if errors.Is(err, app.ErrNotFound) {
		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
		respondWithCode(w, http.StatusNotFound)
		return
	} else if err != nil {
		log.Error().Err(err).Msg("internal server error")
		respondWithCode(w, http.StatusInternalServerError)
		return
	}

	response := apidef.TrainingResponse(t)
	respondWithJSON(w, http.StatusOK, response)
}

// (PUT /trainings/{id})
func (s *SwimLogsServer) EditTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	req, err := readJSON[apidef.EditTrainingRequest](w, r)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	td, err := s.app.EditTraining(id, req)
	if errors.Is(err, app.ErrNotFound) {
		log.Warn().Err(err).Str("id", id.String()).Msg("training not found")
		respondWithCode(w, http.StatusNotFound)
		return
	} else if err != nil {
		log.Warn().Err(err).Msg("invalid training")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	response := apidef.EditTrainingResponse(td)
	respondWithJSON(w, http.StatusOK, response)
}
