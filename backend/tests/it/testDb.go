package it

import (
	"database/sql"
	"os"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/rs/zerolog/log"
)

func TruncateSessions(db *sql.DB) {
	err := data.Sql(db, "truncate sessions cascade")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to truncate sessions")
	}
}

func TruncateTrainings(db *sql.DB) {
	err := data.Sql(db, "truncate trainings cascade")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to truncate trainings")
	}
}

func GetTestDbConf() (data.PostgresConnConf, bool) {
	host, ok := os.LookupEnv("TEST_DB_HOST")
	if !ok {
		log.Error().Msg("TEST_DB_HOST not set")
		return data.PostgresConnConf{}, false
	}

	port, ok := os.LookupEnv("TEST_DB_PORT")
	if !ok {
		log.Error().Msg("TEST_DB_PORT not set")
		return data.PostgresConnConf{}, false
	}

	user, ok := os.LookupEnv("TEST_DB_USER")
	if !ok {
		log.Error().Msg("TEST_DB_USER not set")
		return data.PostgresConnConf{}, false
	}

	pass, ok := os.LookupEnv("TEST_DB_PASSWORD")
	if !ok {
		log.Error().Msg("TEST_DB_PASSWORD not set")
		return data.PostgresConnConf{}, false
	}

	db, ok := os.LookupEnv("TEST_DB_NAME")
	if !ok {
		log.Error().Msg("TEST_DB_NAME not set")
		return data.PostgresConnConf{}, false
	}

	return data.PostgresConnConf{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		Db:   db,
	}, true
}
