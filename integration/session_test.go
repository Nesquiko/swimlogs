package integration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/google/uuid"
)

var (
	validSession   = oapiGen.Session{Day: oapiGen.Saturday, DurationMin: 60, StartTime: "17:00"}
	invalidSession = oapiGen.Session{
		Day:         oapiGen.Day("INVALID-DAY"),
		DurationMin: 60,
		StartTime:   "17:00",
	}
)

// TODO update session
func TestDeleteSession(t *testing.T) {
	sessionReq := oapiGen.CreateSessionRequestObject{Body: &validSession}

	response, err := SwimLogsApp.CreateSession(sessionReq)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	session, ok := response.(oapiGen.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", session)
	}

	id := session.Id

	res, err := SwimLogsApp.DeleteSession(id)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	if _, ok := res.(oapiGen.DeleteSession200Response); !ok {
		t.Fatalf("expected successfull res, but was (of type %t), %v", response, response)
	}
}

func TestDeleteSessionNotFound(t *testing.T) {
	id := uuid.New()
	response, err := SwimLogsApp.DeleteSession(id)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := response.(oapiGen.DeleteSession404JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Session wasn't found"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Session with Id '%s' wasn't found", id.String())
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}
}

func TestGetAllSessions(t *testing.T) {
	sessionReq := oapiGen.CreateSessionRequestObject{Body: &validSession}

	n := 5
	for i := 0; i < n; i++ {
		response, err := SwimLogsApp.CreateSession(sessionReq)
		if err != nil {
			t.Fatalf("expected no error, but was %v", err)
		}
		session, ok := response.(oapiGen.CreateSession201JSONResponse)
		if !ok {
			t.Fatalf("expected successfull reponse, but response was %+v", session)
		}
	}

	response, err := SwimLogsApp.GetAllSessions()
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	sessions, ok := response.(oapiGen.GetAllSessions200JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", sessions)
	}

	if len(*sessions.Sessions) != n {
		t.Fatalf("count of sessions, expected %d, but there were %d", n, len(*sessions.Sessions))
	}
}

func TestCreateSessionValid(t *testing.T) {
	req := oapiGen.CreateSessionRequestObject{Body: &validSession}
	response, err := SwimLogsApp.CreateSession(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	session, ok := response.(oapiGen.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", session)
	}

	if session.Id == uuid.Nil {
		t.Error("id wasn't returned in reponse")
	}
	if session.Day != validSession.Day {
		t.Errorf("day, expected %q, but was %q", validSession.Day, session.Day)
	}
	if session.StartTime != validSession.StartTime {
		t.Errorf("start time, expected %q, but was %q", validSession.StartTime, session.StartTime)
	}
	if session.DurationMin != validSession.DurationMin {
		t.Errorf("duration, expected %q, but was %q", validSession.DurationMin, session.DurationMin)
	}
	cleanDB(DB)
}

func TestCreateSessionInvalid(t *testing.T) {
	req := oapiGen.CreateSessionRequestObject{
		Body: &invalidSession,
	}
	response, err := SwimLogsApp.CreateSession(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errorDetail, ok := response.(oapiGen.CreateSession400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errorDetail)
	}

	expectedTitle := "Invalid request"
	if errorDetail.Title != expectedTitle {
		t.Fatalf("error deatils title, expected %q, but was %q", expectedTitle, errorDetail.Title)
	}

	expectedDetail := "There were invalid session attributes"
	if errorDetail.Detail != expectedDetail {
		t.Fatalf(
			"error deatils detail, expected %q, but was %q",
			expectedDetail,
			errorDetail.Detail,
		)
	}

	expectedAddProps := map[string]string{
		"day": fmt.Sprintf("Unknown day name '%s'", invalidSession.Day),
	}
	if !reflect.DeepEqual(expectedAddProps, errorDetail.AdditionalProperties) {
		t.Fatalf(
			"error details additional properties, expected %v, but was %v",
			expectedAddProps,
			errorDetail.AdditionalProperties,
		)
	}
	cleanDB(DB)
}
