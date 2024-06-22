package data

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const MigrationsLocationFormat = "file://%s"

func (pool *PostgresDbPool) MigrateUp() error {
	migrationsFileUrl := fmt.Sprintf(MigrationsLocationFormat, pool.migrationsDir)
	m, err := migrate.New(migrationsFileUrl, pool.conStr)
	if err != nil {
		return fmt.Errorf("MigrateUp init: %w", err)
	}
	m.Log = slogLogger{verbose: true}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("MigrateUp: %w", err)
	}
	err, dbErr := m.Close()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("MigrateUp: %w", err)
	} else if dbErr != nil {
		return fmt.Errorf("MigrateUp: %w", dbErr)
	}

	return nil
}

// simple wrapper around slog which adheres to migrate.Logger interface
type slogLogger struct {
	verbose bool
}

func (l slogLogger) Printf(format string, v ...interface{}) {
	format = strings.TrimRight(format, "\n")
	msg := fmt.Sprintf(format, v...)
	slog.Info(msg)
}

func (l slogLogger) Verbose() bool {
	return l.verbose
}
