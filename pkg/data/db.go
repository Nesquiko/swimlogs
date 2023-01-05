package data

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var ErrNotFound = errors.New("didn't find row")

type DBConn interface {
	InTx(func(*sql.Tx) error) error

	SaveSession(Session, *sql.Tx) (*uuid.UUID, error)
	GetAllSessions() ([]Session, error)
	DeleteSession(uuid.UUID, *sql.Tx) error
	UpdateSession(id uuid.UUID, updated Session, tx *sql.Tx) (Session, error)

	Close() error
}

func NewDBConn(dsn string) (DBConn, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresDbConn{db}, nil
}

type postgresDbConn struct {
	*sql.DB
}

func (psql *postgresDbConn) InTx(f func(*sql.Tx) error) error {
	tx, err := psql.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (psql *postgresDbConn) Close() error {
	return psql.DB.Close()
}
