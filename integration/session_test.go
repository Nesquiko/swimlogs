package integration

import (
	"fmt"
	"testing"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

func TestUpdateSession(t *testing.T) {
	req := oapiGen.CreateSessionRequestObject{
		Body: &validSession,
	}
	res, err := SwimLogsApp.CreateSession(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	session, ok := res.(oapiGen.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", session)
	}

	expected := oapiGen.Session{
		Day:         oapiGen.Tuesday,
		StartTime:   "09:00",
		DurationMin: 120,
	}
	request := oapiGen.UpdateSessionRequestObject{Id: session.Id, Body: &expected}
	response, err := SwimLogsApp.UpdateSession(request)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	updated, ok := response.(oapiGen.UpdateSession200JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", updated)
	}

	if expected.Day != updated.Day {
		t.Errorf("day, expected %s, but was %s", expected.Day, updated.Day)
	}
	if expected.StartTime != updated.StartTime {
		t.Errorf("start time, expected %s, but was %s", expected.StartTime, updated.StartTime)
	}
	if expected.DurationMin != updated.DurationMin {
		t.Errorf("duration, expected %d, but was %d", expected.DurationMin, updated.DurationMin)
	}

	cleanDB(DB)
}

func TestUpdateSessionNotFound(t *testing.T) {
	request := oapiGen.UpdateSessionRequestObject{Id: validSession.Id, Body: &validSession}
	response, err := SwimLogsApp.UpdateSession(request)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	errDetail, ok := response.(oapiGen.UpdateSession404JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Session wasn't found"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Session with Id '%s' wasn't found", validSession.Id.String())
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	// if len(errDetail.AdditionalProperties) != 0 {
	// 	t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	// }
}

func TestUpdateSessionInvalid(t *testing.T) {
	req := oapiGen.CreateSessionRequestObject{
		Body: &validSession,
	}
	res, err := SwimLogsApp.CreateSession(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	session, ok := res.(oapiGen.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("expected successfull reponse, but response was %+v", session)
	}

	request := oapiGen.UpdateSessionRequestObject{Id: session.Id, Body: &invalidSession}
	response, err := SwimLogsApp.UpdateSession(request)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	errorDetail, ok := response.(oapiGen.UpdateSession400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errorDetail)
	}

	// expectedTitle := "Invalid request"
	// if errorDetail.Title != expectedTitle {
	// 	t.Fatalf("error deatils title, expected %q, but was %q", expectedTitle, errorDetail.Title)
	// }
	//
	// expectedDetail := "There were invalid session attributes"
	// if errorDetail.Detail != expectedDetail {
	// 	t.Fatalf(
	// 		"error deatils detail, expected %q, but was %q",
	// 		expectedDetail,
	// 		errorDetail.Detail,
	// 	)
	// }
	//
	// expectedAddProps := map[string]string{
	// 	"day": fmt.Sprintf("Unknown day name '%s'", invalidSession.Day),
	// }
	// if !reflect.DeepEqual(expectedAddProps, errorDetail.AdditionalProperties) {
	// 	t.Fatalf(
	// 		"error details additional properties, expected %v, but was %v",
	// 		expectedAddProps,
	// 		errorDetail.AdditionalProperties,
	// 	)
	// }
	cleanDB(DB)
}

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
	// if len(errDetail.AdditionalProperties) != 0 {
	// 	t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	// }
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

	if len(sessions.Sessions) != n {
		t.Fatalf("count of sessions, expected %d, but there were %d", n, len(sessions.Sessions))
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

	// expectedTitle := "Invalid request"
	// if errorDetail.Title != expectedTitle {
	// 	t.Fatalf("error deatils title, expected %q, but was %q", expectedTitle, errorDetail.Title)
	// }
	//
	// expectedDetail := "There were invalid session attributes"
	// if errorDetail.Detail != expectedDetail {
	// 	t.Fatalf(
	// 		"error deatils detail, expected %q, but was %q",
	// 		expectedDetail,
	// 		errorDetail.Detail,
	// 	)
	// }
	//
	// expectedAddProps := map[string]string{
	// 	"day": fmt.Sprintf("Unknown day name '%s'", invalidSession.Day),
	// }
	// if !reflect.DeepEqual(expectedAddProps, errorDetail.AdditionalProperties) {
	// 	t.Fatalf(
	// 		"error details additional properties, expected %v, but was %v",
	// 		expectedAddProps,
	// 		errorDetail.AdditionalProperties,
	// 	)
	// }
	cleanDB(DB)
}
