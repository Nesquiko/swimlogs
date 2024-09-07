package server

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-chi/httplog/v2"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

const (
	AppHostDefault        = "localhost"
	AppPortDefault        = "42069"
	LogLevelDefault       = slog.LevelInfo
	DbHostDefault         = "localhost"
	DbPortDefault         = "5432"
	DbUserDefault         = "swimlogs"
	DbPassDefault         = "swimlogs"
	DbNameDefault         = "swimlogs"
	FEOriginDefault       = "http://localhost:3000"
	TzDefault             = "Europe/Bratislava"
	MigrationsPathDefault = "migrations"
	ProdModeDefault       = false
)

func Run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	flags := flag.NewFlagSet("flags", flag.ExitOnError)

	host := flags.String("host", AppHostDefault, "application host")
	port := flags.String("port", AppPortDefault, "application port")
	logLevel := flags.Int(
		"log",
		int(LogLevelDefault),
		"application log level (-4, 0, 4, 8)",
	)

	dbHost := flags.String("db-host", DbHostDefault, "db host")
	dbPort := flags.String("db-port", DbPortDefault, "db port")
	dbUser := flags.String("db-user", DbUserDefault, "user connecting to db")
	dbPass := flags.String("db-pass", DbPassDefault, "password for connecting to db")
	dbName := flags.String("db-name", DbNameDefault, "to which db to connnect")
	migratinsDirPath := flags.String(
		"migrations",
		MigrationsPathDefault,
		"path to dir with migration scripts",
	)
	prodMode := flags.Bool("prod", ProdModeDefault, "start server in production mode")

	feOrigin := flags.String("fe-origin", FEOriginDefault, "frontend origin")
	tz := flags.String("tz", TzDefault, "timezone in which the app is running")
	flags.Parse(args)

	httpLogger := SetupLogger(slog.Level(*logLevel), *prodMode)

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		slog.Error("failed to load timezone", slog.String("error", err.Error()))
		return err
	}
	httpLogger.Info("loaded timezone", slog.String("tz", loc.String()))
	time.Local = loc

	conStr := data.ConnectionString(*dbUser, *dbPass, *dbHost, *dbName, *dbPort)
	pool, err := data.NewPostgresPool(conStr, *migratinsDirPath)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
	}
	defer pool.Close()

	if err := pool.MigrateUp(); err != nil {
		slog.Error("failed to migrate up", slog.String("error", err.Error()))
		os.Exit(1)
	}

	spec, err := api.GetSwagger()
	if err != nil {
		slog.Error("failed to load OpenApi spec", slog.String("error", err.Error()))
		os.Exit(1)
	}
	spec.Servers = nil // running behind proxy

	app := app.New(pool)
	srv := NewServer(app, spec, httpLogger, *feOrigin)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(*host, *port),
		Handler: srv,
	}

	go func() {
		slog.Info("starting server", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error listening and serving", slog.String("error", err.Error()))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		slog.Info("interrupt received, shutting down server")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error shutting down http server", slog.String("error", err.Error()))
		}
	}()
	wg.Wait()
	return nil
}

func SetupLogger(logLevel slog.Level, prodMode bool) *httplog.Logger {
	logger := httplog.NewLogger("swimlogs-api", httplog.Options{
		LogLevel: slog.Level(logLevel),
		Concise:  !prodMode,
	})
	slog.SetDefault(logger.Logger)
	return logger
}
