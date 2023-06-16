package server

import "net/http"

// (GET /monitoring/heartbeat)
func (s *SwimLogsServer) Heartbeat(w http.ResponseWriter, r *http.Request) {
	respondWithCode(w, http.StatusOK)
}
