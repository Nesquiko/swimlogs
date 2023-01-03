package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConn interface {
}

var models = []any{&Session{}}

func NewDBConn(dsn string) (DBConn, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(models...)

	return &postgresDbConn{db}, nil
}

type postgresDbConn struct {
	*gorm.DB
}
