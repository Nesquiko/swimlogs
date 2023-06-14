package server

import (
	"net/http"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/rs/zerolog/log"
)

// (GET /sessions)
func (s *SwimLogsServer) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	res := s.app.GetAllSessions()
	respondWithJSON(w, res.Code(), res.Body())
}

// (POST /sessions)
func (s *SwimLogsServer) CreateSession(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[openapi.CreateSessionJSONBody](w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	res := s.app.SaveSession(req)
	respondWithJSON(w, res.Code(), res.Body())
}

// (DELETE /sessions/{id})
func (s *SwimLogsServer) DeleteSession(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	res := s.app.DeleteSessionById(id)
	if res.HasError() {
		respondWithJSON(w, res.Code(), res.Body())
	} else {
		respondWithCode(w, res.Code())
	}
}

// (PUT /sessions/{id})
func (s *SwimLogsServer) UpdateSession(
	w http.ResponseWriter,
	r *http.Request,
	id types.UUID,
) {
	req, err := readJSON[openapi.UpdateSessionJSONBody](w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		respondWithCode(w, http.StatusBadRequest)
		return
	}

	res := s.app.UpdateSession(id, req)
	respondWithJSON(w, res.Code(), res.Body())
}
