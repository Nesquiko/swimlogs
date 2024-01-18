package server

import (
	"net/http"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/rs/zerolog/log"
)

// (POST /trainings)
func (s *SwimLogsServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[openapi.NewTraining](w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	res := s.app.SaveTraining(req)
	respondWithJSON(w, res.Code(), res.Body())
}

// (GET /trainings/details)
func (s *SwimLogsServer) GetTrainingsDetails(
	w http.ResponseWriter,
	r *http.Request,
	params openapi.GetTrainingsDetailsParams,
) {
	res := s.app.GetTrainingDetails(params)
	respondWithJSON(w, res.Code(), res.Body())
}

// (GET /trainings/details/current-week)
func (s *SwimLogsServer) GetTrainingsDetailsCurrentWeek(w http.ResponseWriter, r *http.Request) {
	res := s.app.GetTrainingDetailsForCurrentWeek()
	respondWithJSON(w, res.Code(), res.Body())
}

// (DELETE /trainings/{id})
func (s *SwimLogsServer) DeleteTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	res := s.app.DeleteTrainingById(id)
	if res.HasError() {
		respondWithJSON(w, res.Code(), res.Body())
		return
	}
	respondWithCode(w, http.StatusNoContent)
}

// (GET /trainings/{id})
func (s *SwimLogsServer) GetTrainingById(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	res := s.app.GetTrainingById(id)
	respondWithJSON(w, res.Code(), res.Body())
}

// (PUT /trainings/{id})
func (s *SwimLogsServer) UpdateTraining(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
}
