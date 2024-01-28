package data

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

const MigrationsLocationFormat = "file://%s"

func (pool *PostgresDbPool) MigrateUp(showLogs bool) error {
	migrationsFileUrl := fmt.Sprintf(MigrationsLocationFormat, pool.migrationsDir)
	m, err := migrate.New(migrationsFileUrl, pool.conStr)
	if err != nil {
		return fmt.Errorf("MigrateUp init: %w", err)
	}
	m.Log = zerologLogger{showLogs: showLogs, verbose: true}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("MigrateUp: %w", err)
	}

	return nil
}

func (pool *PostgresDbPool) CleanMigrateUp(showLogs bool) error {
	migrationsFileUrl := fmt.Sprintf(MigrationsLocationFormat, pool.migrationsDir)
	m, err := migrate.New(migrationsFileUrl, pool.conStr)
	if err != nil {
		return fmt.Errorf("CleanMigrateUp init: %w", err)
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
