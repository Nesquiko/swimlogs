package state

import (
	"errors"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

type SessionStateStorage struct {
	s *oapiGen.Session
}

func (sss *SessionStateStorage) IsInCache(id uuid.UUID) bool {
	if sss.s == nil {
		return false
	}
	return sss.s.Id == id
}

func (sss *SessionStateStorage) SaveSession(s oapiGen.Session) {
	sss.s = &s
}

func (sss *SessionStateStorage) GetSession() (oapiGen.Session, error) {
	if sss.s == nil {
		return oapiGen.Session{}, errors.New("no session stored")
	}
	return *sss.s, nil
}

func (sss *SessionStateStorage) InvalidateTrainingCache() {
	sss.s = nil
}

type SessionStateStorageSetter interface {
	SessionStateStorageSet(sss *SessionStateStorage)
}

// SessionStateStorage ref for wiring
type SessionStateStorageRef struct {
	*SessionStateStorage
}

func (sssr *SessionStateStorageRef) SessionStateStorageSet(sss *SessionStateStorage) {
	sssr.SessionStateStorage = sss
}
