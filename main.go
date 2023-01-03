package main

import (
	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/labstack/echo/v4"
)

func main() {
	dsn := "host=localhost user=swimlogs password=swimlogs dbname=swimlogs port=2345"
	db, err := data.NewDBConn(dsn)
	if err != nil {
		panic(err)
	}

	server := app.NewApp(db)
	handler := oapiGen.NewStrictHandler(server, nil)
	e := echo.New()

	oapiGen.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":42069"))
}
