package main

import (
	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	logger := log.Sugar()

	// dsn := "postgres://swimlogs:swimlogs@localhost:2345/swimlogs?s slmode=disable"
	dsn := "host=localhost port=2345 user=swimlogs password=swimlogs dbname=swimlogs sslmode=disable"
	db, err := data.NewDBConn(dsn)
	if err != nil {
		logger.Fatalf("couldn't connect to db, '%v'", err)
	}
	defer db.Close()

	swimlogs := app.New(db, logger)
	server := server.New(swimlogs)
	handler := oapiGen.NewStrictHandler(server, nil)
	e := echo.New()
	e.Use(middleware.Logger())

	oapiGen.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":42069"))
}
