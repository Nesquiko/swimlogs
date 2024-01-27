package server

import (
	"net/http"

	"github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog/log"

	"github.com/Nesquiko/swimlogs/apidef"
)

// (POST /trainings)
func (s *SwimLogsServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[apidef.CreateTrainingRequest](w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	td, err := s.app.CreateTraining(req)
	if err != nil {
		log.Error().Err(err).Msg("invalid request")
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
}

// (GET /trainings/details/current-week)
func (s *SwimLogsServer) TrainingDetailsCurrentWeek(w http.ResponseWriter, r *http.Request) {
}

// (DELETE /trainings/{id})
func (s *SwimLogsServer) DeleteTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
}

// (GET /trainings/{id})
func (s *SwimLogsServer) TrainingById(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
}

// (PUT /trainings/{id})
func (s *SwimLogsServer) EditTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
}
