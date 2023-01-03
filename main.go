package main

import (
	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
	"github.com/labstack/echo/v4"
)

func main() {
	// dsn := "postgres://swimlogs:swimlogs@localhost:2345/swimlogs?s slmode=disable"
	dsn := "host=localhost port=2345 user=swimlogs password=swimlogs dbname=swimlogs sslmode=disable"
	db, err := data.NewDBConn(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := server.New(db)
	handler := oapiGen.NewStrictHandler(server, nil)
	e := echo.New()

	oapiGen.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":42069"))
}
