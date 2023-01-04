package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConn interface {
	InTx(func(*sql.Tx) error) error
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

func (psql *postgresDbConn) Close() error {
	return psql.DB.Close()
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
