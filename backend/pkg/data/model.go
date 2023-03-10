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
	Version    int
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
	Blocks        []Block
}

type Block struct {
	Id            uuid.UUID
	Num           int
	Repeat        int
	Name          string
	TotalDistance int

	Sets []Set
}

type Set struct {
	Id            uuid.UUID
	Num           int
	Repeat        int
	Distance      int
	What          string
	StartingRule  string
	RuleSeconds   sql.NullInt16
	TotalDistance int
}

type completeTraining struct {
	tId      uuid.UUID
	tVersion int
	tDate    time.Time
	tDay     string
	tStartT  string
	tDur     int
	tTotDist int

	bId      uuid.UUID
	bNum     int
	bRepeat  int
	bName    string
	bTotDist int

	sId        uuid.UUID
	sRepeat    int
	sNum       int
	sDist      int
	sWhat      string
	sStartRule string
	sRuleSecs  sql.NullInt16
	sTotDist   int
}

func createBase() Base {
	return Base{
		Id:         uuid.New(),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Version:    0,
	}
}
