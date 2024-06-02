package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	_ "time/tzdata"

	"github.com/go-chi/httplog/v2"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

const (
	AppHostDefault  = "localhost"
	AppPortDefault  = "42069"
	LogLevelDefault = slog.LevelInfo
	DbHostDefault   = "localhost"
	DbPortDefault   = "5432"
	DbUserDefault   = "swimlogs"
	DbPassDefault   = "swimlogs"
	DbNameDefault   = "swimlogs"
	FEOriginDefault = "http://localhost:3000"
	TzDefault       = "Europe/Bratislava"
)

func run(ctx context.Context, w io.Writer, args []string) error {
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

	feOrigin := flags.String("fe-origin", FEOriginDefault, "frontend origin")
	tz := flags.String("tz", TzDefault, "timezone in which the app is running")
	flags.Parse(args)

	logger := httplog.NewLogger("swimlogs-api", httplog.Options{
		LogLevel: slog.Level(*logLevel),
	})

	loc, err := time.LoadLocation(*tz)
	if err != nil {
		logger.Error("failed to load timezone", slog.String("err", err.Error()))
		return err
	}
	logger.Info("loaded timezone", slog.String("tz", loc.String()))
	time.Local = loc

	conStr := data.ConnectionString(*dbUser, *dbPass, *dbHost, *dbName, *dbPort)
	pool, err := data.NewPostgresPool(conStr, "migrations")
	if err != nil {
		logger.Error("failed to connect to database", slog.String("err", err.Error()))
	}
	defer pool.Close()

	if err := pool.MigrateUp(true); err != nil {
		logger.Error("failed to migrate up", slog.String("err", err.Error()))
	}

	app := app.New(pool)
	srv := server.NewServer(app, logger, *feOrigin)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(*host, *port),
		Handler: srv,
	}

	go func() {
		logger.Info("starting server", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.String("err", err.Error()))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", slog.String("err", err.Error()))
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Println(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
