package it

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

type TestHarness struct {
	ts   *httptest.Server
	pool *data.PostgresDbPool
}

var TH TestHarness

func TestMain(m *testing.M) {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:             os.Stderr,
			FormatTimestamp: func(i any) string { return time.Now().Format("2006-01-02 15:04:05.000") },
		}).
		With().
		Caller().
		Logger()

	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:16.1"),
		postgres.WithDatabase("swimlogs"),
		postgres.WithUsername("swimlogs"),
		postgres.WithPassword("swimlogs"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize container, %s", err.Error())
		os.Exit(1)
	}

	conStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get connection string, %s", err.Error())
		os.Exit(1)
	}

	pool, err := data.NewPostgresPool(conStr, "../../migrations")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create new pool, %s", err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	err = pool.MigrateUp(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to migrate db, %s", err.Error())
		os.Exit(1)
	}

	swimlogs := app.New(pool)
	h := server.NewServerHandler(swimlogs, "")
	ts := httptest.NewServer(h)
	defer ts.Close()

	TH = TestHarness{
		ts:   ts,
		pool: pool,
	}

	exitCode := m.Run()
	if err := pgContainer.Terminate(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to terminate container, %s", err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func asPtr[T any](v T) *T {
	return &v
}

func (th TestHarness) CleanTrainings(t *testing.T) {
	err := data.Sql(th.pool, "truncate trainings cascade")
	require.NoError(t, err)
}
