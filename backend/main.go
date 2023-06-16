package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	AppHostEnvVar = "APP_HOST"
	AppPortEnvVar = "APP_PORT"

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

	dbHost := flag.String("db-host", os.Getenv(DbHostEnvVar), "db host")
	dbPort := flag.String("db-port", os.Getenv(DbPortEnvVar), "db port")
	dbUser := flag.String("db-user", os.Getenv(DbUserEnvVar), "user connecting to db")
	dbPass := flag.String("db-pass", os.Getenv(DbPassEnvVar), "password for connecting to db")
	dbName := flag.String("db-name", os.Getenv(DbNameEnvVar), "to which db to connnect")

	feOrigin := flag.String("fe-origin", os.Getenv(FEOriginEnvVar), "frontend origin")

	jsonLogs := flag.Bool("json-logs", false, "whether to log in json format")
	flag.Parse()

	log.Logger = log.With().Caller().Logger()
	if !*jsonLogs {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:             os.Stderr,
				FormatTimestamp: func(i interface{}) string { return time.Now().Format("2006-01-02 15:04:05.000") },
			}).
			With().Caller().Logger()
	}

	conf := data.PostgresConnConf{
		Host: *dbHost,
		Port: *dbPort,
		User: *dbUser,
		Pass: *dbPass,
		Db:   *dbName,
	}
	db, err := data.NewPostgresConn(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.MigrateUp(true); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate up")
	}

	swimlogs := app.New(db)
	h := server.NewServerHandler(swimlogs, *feOrigin)

	addr := *appHost + ":" + *appPort
	log.Info().Str("addr", addr).Msg("starting server")

	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal().Err(err).Msg("handler failed")
	}
}
