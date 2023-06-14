package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	UniqueViolationCode      = "23505"
	ForeignKeyViolationCode  = "23503"
	CheckViolationCode       = "23514"
	SerializationFailureCode = "40001"
	InvalidEnumTypeCode      = "22P02"
)

var (
	ErrRowsNotFound         = errors.New("didn't find row")
	ErrCheckViolation       = errors.New("check constraint violated")
	ErrUniqueViolation      = errors.New("unique constraint violated")
	ErrSerializationFailure = errors.New("serialization failure")
	ErrInvalidEnumType      = errors.New("invalid enum type")
	ErrForeignKeyViolation  = errors.New("foreign key constraint violated")
)

// NewPostgresConn attempts to connect to a DB specified by the PostgresConnConf.
// If the config is invalid, the program exits.
// If the connection fails, an error is returned.
func NewPostgresConn(conf PostgresConnConf) (*PostgresDbConn, error) {
	checkDbConf(conf)

	db, err := sql.Open("postgres", conf.dsn())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDbConn{conf, db}, nil
}

type PostgresDbConn struct {
	conf PostgresConnConf
	*sql.DB
}

func (psql *PostgresDbConn) Close() error {
	return psql.DB.Close()
}

type PostgresConnConf struct {
	Host          string
	Port          string
	Db            string
	User          string
	Pass          string
	MigrationsDir string
}

func (conf PostgresConnConf) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Pass,
		conf.Db,
	)
}

// checkDbConf checks that the required fields are set in the PostgresConnConf,
// and exits the program if one is missing.
func checkDbConf(conf PostgresConnConf) {
	isValid := true
	if conf.Host == "" {
		log.Error().Msg("host is not set")
		isValid = false
	}
	if conf.Port == "" {
		log.Error().Msg("port is not set")
		isValid = false
	}
	if conf.Db == "" {
		log.Error().Msg("db is not set")
		isValid = false
	}
	if conf.User == "" {
		log.Error().Msg("user is not set")
		isValid = false
	}
	if conf.Pass == "" {
		log.Error().Msg("pass is not set")
		isValid = false
	}

	if !isValid {
		log.Fatal().Msg("invalid db config")
	}
}

func Sql(db *sql.DB, sql string, args ...any) error {
	_, err := db.Exec(sql, args...)
	if err != nil {
		log.Error().Err(err).Msg("sql")
		return fmt.Errorf("sql: %w", err)
	}
	return nil
}

func SqlWithResult[R any](db *sql.DB, sql string, args ...any) (R, error) {
	var res R
	err := db.QueryRow(sql, args...).Scan(&res)
	if err != nil {
		log.Error().Err(err).Msg("sql with result")
		return res, fmt.Errorf("sql with result: %w", err)
	}
	return res, nil
}

func tx(db *sql.DB, f func(*sql.Tx) error) error {
	return txWithOpts(db, nil, f)
}

func txWithOpts(db *sql.DB, opts *sql.TxOptions, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(context.Background(), opts)
	if err != nil {
		log.Error().Err(err).Msg("begin tx")
		return err
	}

	err = f(tx)
	if err != nil {
		_ = tx.Rollback()
		log.Warn().Err(err).Msg("rolling back tx")
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("commit tx")
		return err
	}

	return nil
}

func txWithResult[R any](db *sql.DB, f func(*sql.Tx) (R, error)) (R, error) {
	return txWithResultAndOpts(db, nil, f)
}

func txWithResultAndOpts[R any](
	db *sql.DB,
	opts *sql.TxOptions,
	f func(*sql.Tx) (R, error),
) (R, error) {
	var res R
	tx, err := db.BeginTx(context.Background(), opts)
	if err != nil {
		log.Error().Err(err).Msg("begin txWithResult")
		return res, err
	}

	res, err = f(tx)
	if err != nil {
		_ = tx.Rollback()
		log.Warn().Err(err).Msg("rolling back txWithResult")
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("commit txWithResult")
		return res, err
	}

	return res, nil
}
