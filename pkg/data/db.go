package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConn interface {
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
