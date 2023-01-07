package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var InsertSession = "insert into session (id, created_at, modified_at, day, starttime, duration) values ($1, $2, $3, $4, $5, $6)"

func (db *postgresDbConn) SaveSession(session Session, tx *sql.Tx) (*uuid.UUID, error) {
	base := createBase()
	session.Base = base

	_, err := tx.Exec(
		InsertSession,
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

var SelectSession = "select * from session"

func (db *postgresDbConn) GetAllSessions() ([]Session, error) {
	sessions := make([]Session, 0)

	rows, err := db.Query(SelectSession)
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

		st, err := time.Parse(TimeLayout, startTime)
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

var DeleteSession = "delete from session where id = $1"

func (psql *postgresDbConn) DeleteSession(id uuid.UUID, tx *sql.Tx) error {
	res, err := tx.Exec(DeleteSession, id)
	if err != nil {
		return fmt.Errorf("DeleteSession: %w", err)
	}

	// Err can be ignored, because Postgres supports rows affected
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrRowNotFound
	}

	return nil
}

var UpdateSession = "update session set modified_at = now(), day = $2, starttime = $3, duration = $4 where id = $1 returning id, day, starttime, duration"

func (psql *postgresDbConn) UpdateSession(
	id uuid.UUID,
	updated Session,
	tx *sql.Tx,
) (Session, error) {
	var session Session
	var startTime string
	err := tx.QueryRow(
		UpdateSession,
		id,
		updated.Day,
		updated.StartTime,
		updated.DurationMin,
	).Scan(
		&session.Id,
		&session.Day,
		&startTime,
		&session.DurationMin,
	)
	if err == sql.ErrNoRows {
		return session, ErrRowNotFound
	} else if err != nil {
		return session, fmt.Errorf("UpdateSession: %w", err)
	}

	st, err := time.Parse(TimeLayout, startTime)
	if err != nil {
		return session, fmt.Errorf("UpdateSession: %w", err)
	}
	session.StartTime = st.Format("15:04")
	return session, nil
}
