package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Session struct {
	Id          uuid.UUID
	Day         string
	StartTime   time.Time
	DurationMin int

	CreatedAt  time.Time
	ModifiedAt time.Time
}

var insertSession = `
insert into sessions (day, start_time, duration_min, created_at, modified_at)
values ($1::day, $2, $3, $4, $4)
returning id, day, start_time, duration_min, created_at, modified_at
`

func (db *PostgresDbConn) SaveSession(session Session) (Session, error) {
	var s Session

	err := db.QueryRow(insertSession, session.Day, session.StartTime, session.DurationMin, time.Now()).
		Scan(&s.Id, &s.Day, &s.StartTime, &s.DurationMin, &s.CreatedAt, &s.ModifiedAt)

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case CheckViolationCode:
			return Session{}, ErrCheckViolation
		case UniqueViolationCode:
			return Session{}, ErrUniqueViolation
		case InvalidEnumTypeCode:
			return Session{}, ErrInvalidEnumType
		default:
			return Session{}, fmt.Errorf("SaveSession: %w", err)
		}
	} else if err != nil {
		return Session{}, fmt.Errorf("SaveSession: %w", err)
	}

	return s, nil
}

var selectSessions = `
select id, day, start_time, duration_min, count(*) over ()
from sessions
order by day::day, start_time, duration_min
limit $1 offset $2
`

// GetSessions returns a paginated list of sessions and total count of sessions.
// Page starts at 0
func (db *PostgresDbConn) GetSessions(page, pageSize int) ([]Session, int, error) {
	sessions := make([]Session, 0)

	rows, err := db.Query(selectSessions, pageSize, page*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAllSessions: %w", err)
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var session Session
		err := rows.Scan(
			&session.Id,
			&session.Day,
			&session.StartTime,
			&session.DurationMin,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("GetAllSessions: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GetAllSessions: %w", err)
	}

	return sessions, total, nil
}

var deleteSession = "delete from sessions where id = $1"

func (db *PostgresDbConn) DeleteSessionById(id uuid.UUID) error {
	res, err := db.Exec(deleteSession, id)
	if err != nil {
		return fmt.Errorf("DeleteSession: %w", err)
	}
	// Err can be ignored, because Postgres supports rows affected
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrRowsNotFound
	}
	return nil
}

var sessionExists = "select count(1) from session where id = $1"

func (psql *PostgresDbConn) SessionExists(id uuid.UUID) (bool, error) {
	var exists int
	err := psql.QueryRow(sessionExists, id).Scan(&exists)
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
