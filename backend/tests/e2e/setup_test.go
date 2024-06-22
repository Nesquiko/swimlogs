//go:build integration

package e2e

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

const MigrationsTestDirPath = "../../migrations"

// URL of the test server, initializez in the TestMain function
var ServerUrl string

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logLevel := slog.LevelDebug
	server.SetupLogger(logLevel, false)

	pgContainer := preparePostgres(ctx)
	portBindings, err := pgContainer.Ports(ctx)
	if err != nil {
		slog.Error("couldn't retrive test containers ports", slog.String("error", err.Error()))
		os.Exit(1)
	}
	port := portBindings["5432/tcp"][0]

	appHost := "127.0.0.1"
	appPort := "42069"
	testFlags := []string{
		"-host", appHost,
		"-port", appPort,
		"-log", strconv.Itoa(int(logLevel)),
		"-db-host", port.HostIP,
		"-db-port", port.HostPort,
		"-db-user", "swimlogs",
		"-db-pass", "swimlogs",
		"-db-name", "swimlogs",
		"-migrations", MigrationsTestDirPath,
	}
	ServerUrl = fmt.Sprintf("http://%s:%s", appHost, appPort)

	go server.Run(ctx, testFlags)
	err = waitForReady(
		ctx,
		5*time.Second,
		200*time.Millisecond,
		ServerUrl+"/monitoring/heartbeat",
	)
	if err != nil {
		slog.Error("ready endpoint not answering", slog.String("error", err.Error()))
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		slog.Error("failed to terminate container", slog.String("error", err.Error()))
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func preparePostgres(ctx context.Context) *postgres.PostgresContainer {
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
		slog.Error("failed to initialize container", slog.String("error", err.Error()))
		os.Exit(1)
	}
	conStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		slog.Error("failed to get connection string", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pool, err := data.NewPostgresPool(conStr, MigrationsTestDirPath)
	if err != nil {
		slog.Error("failed to create new pool", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = pool.MigrateUp()
	if err != nil {
		slog.Error("failed to migrate db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pool.Close()

	err = pgContainer.Snapshot(ctx, postgres.WithSnapshotName("clean"))
	if err != nil {
		slog.Error("failed to make snapshot", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return pgContainer
}
