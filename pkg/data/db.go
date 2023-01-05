package data

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// ErrRowNotFound indicated no rows were found
var ErrRowNotFound = errors.New("didn't find row")

// DBConn is an interface for persisting and reading data from storage
type DBConn interface {
	// InTx executes given function in a transaction, reverting all changes
	// made in the transaction, if function returns an error
	InTx(f func(*sql.Tx) error) error

	// SaveSession saves passed session into storage, MUST be executed in
	// transaction (see DBConn.InTx)
	SaveSession(s Session, tx *sql.Tx) (*uuid.UUID, error)

	// GetAllSessions returns all sessions saved in storage
	GetAllSessions() ([]Session, error)

	// DeleteSession deletes a session with matching id, if no session with
	// matching uuid id, returns ErrRowNotFound. MUST be executed in
	// transaction (see DBConn.InTx)
	DeleteSession(id uuid.UUID, tx *sql.Tx) error

	// UpdateSession updates a session with matching id with values from updated.
	// Only Day, StartTime, DurationMin can are updated, other fields are ignored.
	// Also ModifiedAt is updated to current time of execution.
	UpdateSession(id uuid.UUID, updated Session, tx *sql.Tx) (Session, error)

	// Closes the db connection
	Close() error
}

// NewDBConn attempts to connect to a DB specified by the dsn, if it can't,
// it retunrs an error.
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
