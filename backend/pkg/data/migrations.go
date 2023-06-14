package data

import (
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (db *PostgresDbConn) MigrateUp(showLogs bool) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Error().Err(err).Msg("failed to create database driver")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Error().Err(err).Msg("failed to create migration instance")
		return err
	}
	m.Log = zerologLogger{showLogs: showLogs, verbose: true}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to migrate up")
		return err
	}

	return nil
}

func (db *PostgresDbConn) CleanMigrateUp(showLogs bool) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Error().Err(err).Msg("failed to create database driver")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+db.conf.MigrationsDir, "postgres", driver)
	if err != nil {
		log.Error().Err(err).Msg("failed to create migration instance")
		return err
	}
	m.Log = zerologLogger{showLogs: showLogs, verbose: true}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to migrate down")
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Error().Err(err).Msg("failed to migrate up")
		return err
	}

	return nil
}

type zerologLogger struct {
	showLogs bool
	verbose  bool
}

func (l zerologLogger) Printf(format string, v ...interface{}) {
	if !l.showLogs {
		return
	}
	format = strings.TrimRight(format, "\n")
	log.Info().Msgf(format, v...)
}

func (l zerologLogger) Verbose() bool {
	return l.verbose
}
