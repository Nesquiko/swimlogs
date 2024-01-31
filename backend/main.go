package main

import (
	"flag"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

const (
	AppHostEnvVar       = "APP_HOST"
	AppPortEnvVar       = "APP_PORT"
	AppDebugLevelEnvVar = "APP_DEBUG_LEVEL"

	DbHostEnvVar = "DATABASE_HOST"
	DbPortEnvVar = "DATABASE_PORT"
	DbNameEnvVar = "DATABASE_NAME"
	DbUserEnvVar = "DATABASE_USER"
	DbPassEnvVar = "DATABASE_PASSWORD"

	FEOriginEnvVar = "FE_ORIGIN"
)

func main() {
	appHost := flag.String("host", os.Getenv(AppHostEnvVar), "application host")
	appPort := flag.String("port", os.Getenv(AppPortEnvVar), "application port")
	appDebugLevel := flag.Int(
		"debug-level",
		1,
		"application debug level (trace = -1, debug = 0, info = 1, warn = 2, error = 3, fatal = 4, panic = 5)",
	)

	dbHost := flag.String("db-host", os.Getenv(DbHostEnvVar), "db host")
	dbPort := flag.String("db-port", os.Getenv(DbPortEnvVar), "db port")
	dbUser := flag.String("db-user", os.Getenv(DbUserEnvVar), "user connecting to db")
	dbPass := flag.String("db-pass", os.Getenv(DbPassEnvVar), "password for connecting to db")
	dbName := flag.String("db-name", os.Getenv(DbNameEnvVar), "to which db to connnect")

	feOrigin := flag.String("fe-origin", os.Getenv(FEOriginEnvVar), "frontend origin")
	_ = feOrigin

	jsonLogs := flag.Bool("json-logs", false, "whether to log in json format")

	tz := flag.String("tz", os.Getenv("TZ"), "timezone in which the app is running")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.Level(*appDebugLevel))
	log.Logger = log.With().Caller().Logger()
	if !*jsonLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, FormatTimestamp: func(i interface{}) string { return time.Now().Format("2006-01-02 15:04:05.000") }}).
			With().
			Caller().
			Logger()
	}

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load timezone")
	}
	log.Info().Str("tz", loc.String()).Msg("loaded timezone")
	time.Local = loc

	conStr := data.ConnectionString(*dbUser, *dbPass, *dbHost, *dbName, *dbPort)
	pool, err := data.NewPostgresPool(conStr, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := pool.MigrateUp(true); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate up")
	}

	swimlogs := app.New(pool)
	h := server.NewServerHandler(swimlogs, *feOrigin)

	addr := *appHost + ":" + *appPort
	log.Info().Str("addr", addr).Msg("starting server")

	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal().Err(err).Msg("handler failed")
	}
}
