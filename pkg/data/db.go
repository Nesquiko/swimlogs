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

	// GetSessionById returns a session with matching id, if it doesn't exist
	// returns ErrRowNotFound
	GetSessionById(id uuid.UUID) (Session, error)

	// SaveTraining saves passed training into storage, MUST be executed in
	// transaction (see DBConn.InTx)
	SaveTraining(t Training, tx *sql.Tx) (*uuid.UUID, error)

	// SaveTrainingWithSesssionData saves passed training with data from session
	// with matching sId into storage. Returns Training with only values copied
	// from session with matching sId (Day, StartTime, DurationMin).
	// MUST be executed in transaction (see DBConn.InTx)
	SaveTrainingWithSesssionData(t Training, sId uuid.UUID, tx *sql.Tx) (*Training, error)

	// GetTraining returns paginated list with length of pageSize and with
	// offset of page
	GetTrainings(page, pageSize int) ([]Training, error)

	// DeleteTraining deletes a training with matching id, if no training with
	// matching id exists, returns ErrRowNotFound. Also all the blocks and
	// sets are deleted. MUST be executed in transaction (see DBConn.InTx)
	DeleteTraining(id uuid.UUID, tx *sql.Tx) error

	// GetTrainingById returns a training with matching id. If it doesn't exist
	// returns ErrRowNotFound.
	GetTrainingById(id uuid.UUID) (Training, error)

	// UpdateTrainingById updates training with matching id accordingly to t.
	// MUST be executed in transaction (see DBConn.InTx)
	UpdateTrainingById(id uuid.UUID, t Training, tx *sql.Tx) error

	// GetDetailsOfTrainings returns paginated list of Training with only values
	// needed in detail.
	GetDetailsOfTrainings(page, pageSize int) ([]Training, error)

	// GetDetailsOfTrainingsCurrentWeek returns list of Trainings, in current
	// week, with only values needed in detail.
	GetDetailsOfTrainingsCurrentWeek() ([]Training, error)

	// GetTrainingCount return how many trainings are save in storage
	GetTrainingCount() (int, error)

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
