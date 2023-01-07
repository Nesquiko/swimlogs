package integration

import (
	"fmt"
	"os"
	"testing"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"go.uber.org/zap"
)

var (
	DB          data.DBConn
	SwimLogsApp app.SwimLogs
)

func TestMain(m *testing.M) {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	logger := log.Sugar()

	testDbPort := os.Getenv("TESTDB_PORT")
	testDsn := fmt.Sprintf(
		"host=localhost port=%s user=swimlogs_test password=swimlogs_test dbname=swimlogs_test sslmode=disable",
		testDbPort,
	)
	DB, err = data.NewDBConn(testDsn)
	if err != nil {
		logger.Fatalf("Couldn't connect to testing database!, error: %v", err)
		return
	}
	defer DB.Close()

	SwimLogsApp = app.New(DB, logger)

	os.Exit(m.Run())
}
