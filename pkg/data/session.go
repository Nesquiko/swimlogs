package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	Friday    string = "friday"
	Monday    string = "monday"
	Saturday  string = "saturday"
	Sunday    string = "sunday"
	Thursday  string = "thursday"
	Tuesday   string = "tuesday"
	Wednesday string = "wednesday"
)

const TIME_LAYOUT = "2006-01-02T15:04:05Z"

const (
	INSERT_SESSION = `insert into session (id, created_at, modified_at, day, starttime, duration) values ($1, $2, $3, $4, $5, $6)`
	SELECT_SESSION = `select * from session`
)

type Session struct {
	Base
	Day         string
	DurationMin int
	StartTime   string
}

func (db *postgresDbConn) SaveSession(session Session, tx *sql.Tx) (*uuid.UUID, error) {
	base := createBase()
	session.Base = base

	_, err := tx.Exec(
		INSERT_SESSION,
		session.Id,
		session.CreatedAt,
		session.ModifiedAt,
		session.Day,
		session.StartTime,
		session.DurationMin,
	)
	if err != nil {
		return nil, fmt.Errorf("SaveSession: %w", err)
	}

	return &session.Id, nil
}

func (db *postgresDbConn) GetAllSessions() ([]Session, error) {
	sessions := make([]Session, 0)

	rows, err := db.Query(SELECT_SESSION)
	if err != nil {
		return nil, fmt.Errorf("GetAllSessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		var startTime string
		err := rows.Scan(
			&session.Id,
			&session.CreatedAt,
			&session.ModifiedAt,
			&session.Day,
			&startTime,
			&session.DurationMin,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllSessions: %w", err)
		}

		st, err := time.Parse(TIME_LAYOUT, startTime)
		if err != nil {
			return nil, fmt.Errorf("GetAllSessions: %w", err)
		}
		session.StartTime = st.Format("15:04")

		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllSessions: %w", err)
	}

	return sessions, nil
}
