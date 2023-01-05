package data

import (
	"database/sql"

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

const INSERT_SESSION = `
insert into session
(id, created_at, modified_at, day, starttime, duration)
values ($1, $2, $3, $4, $5, $6)
`

type Session struct {
	*Base
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
		return nil, err
	}

	return &session.Id, nil
}
