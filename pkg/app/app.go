package app

import "github.com/Nesquiko/swimlogs/pkg/data"

type SwimLogs interface {
}

func New(db data.DBConn) SwimLogs {
	return &swimLogsApp{}
}

type swimLogsApp struct {
}
