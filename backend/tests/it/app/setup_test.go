package app

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/tests/it"
)

var (
	SwimLogsApp    app.SwimLogsApp
	PostgresDbConn *data.PostgresDbConn
)

func TestMain(m *testing.M) {
	it.MainFilter()
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:             os.Stderr,
			FormatTimestamp: func(i any) string { return time.Now().Format("2006-01-02 15:04:05.000") },
		}).
		With().
		Caller().
		Logger()

	testConf, ok := it.GetTestDbConf()
	if !ok {
		log.Fatal().Msg("failed to get test database configuration")
		os.Exit(1)
	}
	testConf.MigrationsDir = "../../../migrations"
	db, err := data.NewPostgresConn(testConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.CleanMigrateUp(false); err != nil {
		log.Fatal().Err(err).Msg("failed to clean and migrate up")
	}

	SwimLogsApp = app.New(db)
	PostgresDbConn = db

	os.Exit(m.Run())
}
