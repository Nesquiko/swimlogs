package data

import (
	"database/sql"
	"fmt"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var insertSession = `
insert into sessions (day, start_time, duration, created_at, modified_at)
values ($1::day, $2, $3, now(), now())
returning id, day, start_time, duration
`

func (db *PostgresDbConn) SaveSession(
	session openapi.CreateSessionJSONBody,
) (openapi.Session, error) {
	var s openapi.Session

	err := db.QueryRow(
		insertSession,
		session.Day,
		session.StartTime,
		session.DurationMin,
	).Scan(
		&s.Id,
		&s.Day,
		&s.StartTime,
		&s.DurationMin,
	)
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case CheckViolationCode:
			return openapi.Session{}, ErrCheckViolation
		case UniqueViolationCode:
			return openapi.Session{}, ErrUniqueViolation
		case InvalidEnumTypeCode:
			return openapi.Session{}, ErrInvalidEnumType
		default:
			return openapi.Session{}, fmt.Errorf("SaveSession: %w", err)
		}
	} else if err != nil {
		return openapi.Session{}, fmt.Errorf("SaveSession: %w", err)
	}

	return s, nil
}

var selectSessions = `
select
    id,
    day,
    start_time,
    duration
from sessions
order by day::day, start_time, duration
`

func (db *PostgresDbConn) GetAllSessions() ([]openapi.Session, error) {
	sessions := make([]openapi.Session, 0)

	rows, err := db.Query(selectSessions)
	if err != nil {
		return nil, fmt.Errorf("GetAllSessions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var session openapi.Session
		err := rows.Scan(
			&session.Id,
			&session.Day,
			&session.StartTime,
			&session.DurationMin,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllSessions: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllSessions: %w", err)
	}

	return sessions, nil
}

var DeleteSession = "delete from sessions where id = $1"

func (db *PostgresDbConn) DeleteSessionById(id uuid.UUID) error {
	res, err := db.Exec(DeleteSession, id)
	if err != nil {
		return fmt.Errorf("DeleteSession: %w", err)
	}
	// Err can be ignored, because Postgres supports rows affected
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrRowsNotFound
	}
	return nil
}

var SessionExists = "select count(1) from session where id = $1"

func (psql *PostgresDbConn) SessionExists(id uuid.UUID) (bool, error) {
	var exists int
	err := psql.QueryRow(SessionExists, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

var updateSession = `
with saved as (select
                   id,
                   day,
                   start_time,
                   duration
               from sessions
               where id = $1)
update sessions s
set day        = coalesce($2, saved.day),
    start_time = coalesce($3, saved.start_time),
    duration   = coalesce($4, saved.duration)
from saved
where s.id = saved.id
returning s.id, s.day, s.start_time, s.duration
`

var TxRepeatableRead = sql.TxOptions{
	Isolation: sql.LevelRepeatableRead,
	ReadOnly:  false,
}

func (db *PostgresDbConn) UpdateSession(
	id uuid.UUID,
	updated openapi.UpdateSessionJSONBody,
) (openapi.Session, error) {
	res, err := txWithResultAndOpts(
		db.DB,
		&TxRepeatableRead,
		func(tx *sql.Tx) (openapi.Session, error) {
			return updateSessionInTx(id, updated, tx)
		},
	)
	if err != nil {
		return openapi.Session{}, err
	}
	return res, nil
}

func updateSessionInTx(
	id uuid.UUID,
	updated openapi.UpdateSessionJSONBody,
	tx *sql.Tx,
) (openapi.Session, error) {
	var s openapi.Session

	err := tx.QueryRow(
		updateSession,
		id,
		updated.Day,
		updated.StartTime,
		updated.DurationMin,
	).Scan(
		&s.Id,
		&s.Day,
		&s.StartTime,
		&s.DurationMin,
	)

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case CheckViolationCode:
			return openapi.Session{}, ErrCheckViolation
		case UniqueViolationCode:
			return openapi.Session{}, ErrUniqueViolation
		case InvalidEnumTypeCode:
			return openapi.Session{}, ErrInvalidEnumType
		case SerializationFailureCode:
			return openapi.Session{}, ErrSerializationFailure
		default:
			return openapi.Session{}, fmt.Errorf("updateSessionInTx: %w", err)
		}
	} else if err == sql.ErrNoRows {
		return openapi.Session{}, ErrRowsNotFound
	} else if err != nil {
		return openapi.Session{}, fmt.Errorf("updateSessionInTx: %w", err)
	}

	return s, nil
}
