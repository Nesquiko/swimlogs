package integration

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
)

var (
	CleanSession  = "truncate table session"
	CleanTraining = "delete from training"
)

func cleanDB(db data.DBConn) {
	field := reflect.ValueOf(db).Elem().FieldByName("DB")

	if !field.IsValid() {
		log.Fatal("couldn't get the field")
	}

	sqlDB, ok := field.Interface().(*sql.DB)
	if !ok {
		log.Fatal(ok)
	}

	_, err := sqlDB.Exec(CleanSession)
	if err != nil {
		log.Fatal(err)
	}
	_, err = sqlDB.Exec(CleanTraining)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	validSession   = oapiGen.Session{Day: oapiGen.Saturday, DurationMin: 60, StartTime: "17:00"}
	invalidSession = oapiGen.Session{
		Day:         oapiGen.Day("INVALID-DAY"),
		DurationMin: 60,
		StartTime:   "17:00",
	}
)

func saveSession(s oapiGen.Session, t *testing.T) *uuid.UUID {
	req := oapiGen.CreateSessionRequestObject{
		Body: &s,
	}
	res, err := SwimLogsApp.CreateSession(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	session, ok := res.(oapiGen.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", session)
	}

	return &session.Id
}

func saveTraining(training oapiGen.Training, t *testing.T) *uuid.UUID {
	req := oapiGen.CreateTrainingRequestObject{Body: &training}
	res, err := SwimLogsApp.CreateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	trainingDetail, ok := res.(oapiGen.CreateTraining201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull response, but response was %+v", res)
	}
	return &trainingDetail.Id
}

func createValidTraining(sessionId *uuid.UUID) oapiGen.Training {
	totDist := 400
	day := oapiGen.Monday
	date, err := time.Parse("02.01.2006", "16.01.2023")
	if err != nil {
		panic(err)
	}
	dur := 60
	startTime := "16:00"
	return oapiGen.Training{
		Blocks: []oapiGen.Block{
			{
				Name:   "Warmp up",
				Num:    0,
				Repeat: 1,
				Sets: []oapiGen.Set{
					{
						Distance:     400,
						Num:          0,
						Repeat:       1,
						StartingRule: oapiGen.StartingRule{Rule: oapiGen.None},
						TotalDist:    &totDist,
					},
				},
				TotalDist: &totDist,
			},
		},
		Date:        types.Date{Time: date},
		Day:         &day,
		DurationMin: &dur,
		StartTime:   &startTime,
		SessionId:   sessionId,
	}
}

func createTrainingWithNoBlocks() oapiGen.Training {
	day := oapiGen.Monday
	date, err := time.Parse("02.01.2006", "16.01.2023")
	if err != nil {
		panic(err)
	}
	dur := 60
	startTime := "16:00"
	return oapiGen.Training{
		Blocks:      []oapiGen.Block{},
		Date:        types.Date{Time: date},
		Day:         &day,
		DurationMin: &dur,
		StartTime:   &startTime,
		SessionId:   nil,
	}
}

func getNameOfDay(year, month, day int) oapiGen.Day {
	str := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	dayTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		panic(err)
	}
	switch dayTime.Weekday() {
	case time.Monday:
		return oapiGen.Monday
	case time.Tuesday:
		return oapiGen.Tuesday
	case time.Wednesday:
		return oapiGen.Wednesday
	case time.Thursday:
		return oapiGen.Thursday
	case time.Friday:
		return oapiGen.Friday
	case time.Saturday:
		return oapiGen.Saturday
	case time.Sunday:
		return oapiGen.Sunday
	}

	panic("why is ther week day: " + dayTime.Weekday().String())
}
