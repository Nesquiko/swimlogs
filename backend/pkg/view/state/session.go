package state

import (
	"errors"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/view/util"
)

type SessionStateStorage struct {
	editSession *oapiGen.Session
	sessions    []oapiGen.Session

	reloadSessions bool
}

func (sss *SessionStateStorage) InsertNewSession(s oapiGen.Session) {
	sss.sessions = append(sss.sessions, s)
	util.OrderByDays(&sss.sessions)
	sss.reloadSessions = true
}

func (sss *SessionStateStorage) SaveEditSession(s oapiGen.Session) {
	sss.editSession = &s
}

func (sss *SessionStateStorage) GetEditSession() (oapiGen.Session, error) {
	if sss.editSession == nil {
		return oapiGen.Session{}, errors.New("no session stored")
	}
	return *sss.editSession, nil
}

func (sss *SessionStateStorage) SaveSessions(sessions []oapiGen.Session) {
	sss.sessions = sessions
}

func (sss *SessionStateStorage) UpdateSessionInStorage(s oapiGen.Session) {
	for i := range sss.sessions {
		if sss.sessions[i].Id == s.Id {
			sss.sessions[i].Day = s.Day
			sss.sessions[i].StartTime = s.StartTime
			sss.sessions[i].DurationMin = s.DurationMin
			sss.sessions[i].Version = s.Version
			break
		}
	}
}

func (sss *SessionStateStorage) GetSessions() ([]oapiGen.Session, error) {
	if sss.sessions == nil || len(sss.sessions) == 0 || sss.reloadSessions {
		return nil, errors.New("no sessions stored")
	}
	return sss.sessions, nil
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
