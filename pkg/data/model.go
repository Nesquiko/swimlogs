package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	Friday    string = "friday"
	Monday    string = "monday"
	Saturday  string = "saturday"
	Sunday    string = "sunday"
	Thursday  string = "thursday"
	Tuesday   string = "tuesday"
	Wednesday string = "wednesday"
)

const TimeLayout = "2006-01-02T15:04:05Z"

const (
	Interval string = "interval"
	None     string = "none"
	Pause    string = "pause"
)

// Base contains essential fields used by database tables
type Base struct {
	Id         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Session struct {
	Base
	Day         string
	DurationMin int
	StartTime   string
}

type Training struct {
	Base

	Date          time.Time
	Day           *string
	DurationMin   *int
	StartTime     *string
	TotalDistance int

	Blocks []Block
}

type Block struct {
	Id uuid.UUID

	Repeat        int
	Name          string
	TotalDistance int

	Sets []Set
}

type Set struct {
	Id uuid.UUID

	Repeat        int
	Distance      int
	What          string
	StartingRule  string
	RuleSeconds   sql.NullInt16
	TotalDistance int
}

func createBase() Base {
	return Base{
		Id:         uuid.New(),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}
