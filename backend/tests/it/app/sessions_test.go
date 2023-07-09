package app

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveSessionSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	req := openapi.CreateSessionJSONBody{
		Day:         openapi.Friday,
		DurationMin: 60,
		StartTime:   "10:00",
	}

	res := SwimLogsApp.SaveSession(req)
	assert := assert.New(t)
	assert.Equal(http.StatusCreated, res.Code())

	require.IsType(t, openapi.Session{}, res.Body())
	session := res.Body().(openapi.Session)
	assert.Equal(req.Day, session.Day)
	assert.Equal(req.DurationMin, session.DurationMin)
	assert.Equal(req.StartTime, session.StartTime)
}

func TestSaveSessionInvalidDay(t *testing.T) {
	it.TestFilter(t)
	req := openapi.CreateSessionJSONBody{
		Day:         openapi.Day("invalid"),
		DurationMin: 60,
		StartTime:   "10:00",
	}

	res := SwimLogsApp.SaveSession(req)
	assert := assert.New(t)
	assert.Equal(http.StatusBadRequest, res.Code())

	require.IsType(t, openapi.InvalidSession{}, res.Body())
	invalid := res.Body().(openapi.InvalidSession)
	assert.Equal(asPtr(fmt.Sprintf(app.DayErrFormat, req.Day)), invalid.Day)
}

func TestSaveSessionDuplicate(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	req := openapi.CreateSessionJSONBody{
		Day:         openapi.Friday,
		DurationMin: 60,
		StartTime:   "10:00",
	}

	res := SwimLogsApp.SaveSession(req)
	assert := assert.New(t)
	assert.Equal(http.StatusCreated, res.Code())

	res = SwimLogsApp.SaveSession(req)
	assert.Equal(http.StatusConflict, res.Code())
	assert.IsType(openapi.ErrorDetail{}, res.Body())
}

func TestGetSessions(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	for _, d := range []openapi.Day{openapi.Monday, openapi.Tuesday, openapi.Wednesday, openapi.Thursday, openapi.Friday, openapi.Saturday, openapi.Sunday} {
		newSess := openapi.CreateSessionJSONBody{
			Day:         d,
			StartTime:   "10:00",
			DurationMin: 60,
		}
		res := SwimLogsApp.SaveSession(newSess)
		require.IsType(t, openapi.Session{}, res.Body())
	}

	res := SwimLogsApp.GetSessions(openapi.GetSessionsParams{Page: 0, PageSize: 5})
	assert := assert.New(t)
	assert.Equal(http.StatusOK, res.Code())

	require.IsType(t, openapi.SessionsResponse{}, res.Body())
	sessions := res.Body().(openapi.SessionsResponse)
	assert.Len(sessions.Sessions, 5)
	assert.Equal(7, sessions.Pagination.Total)
	assert.Equal(0, sessions.Pagination.Page)
	assert.Equal(5, sessions.Pagination.PageSize)

	res = SwimLogsApp.GetSessions(openapi.GetSessionsParams{Page: 1, PageSize: 5})
	assert.Equal(http.StatusOK, res.Code())

	require.IsType(t, openapi.SessionsResponse{}, res.Body())
	sessions = res.Body().(openapi.SessionsResponse)
	assert.Len(sessions.Sessions, 2)
	assert.Equal(7, sessions.Pagination.Total)
	assert.Equal(1, sessions.Pagination.Page)
	assert.Equal(2, sessions.Pagination.PageSize)
}

func TestDeleteSessionByIdNotFound(t *testing.T) {
	it.TestFilter(t)
	res := SwimLogsApp.DeleteSessionById(uuid.New())
	assert.Equal(t, http.StatusNotFound, res.Code())
	assert.IsType(t, openapi.ErrorDetail{}, res.Body())
}

func TestDeleteSessionById(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Friday,
		StartTime:   "10:00",
		DurationMin: 60,
	}
	s := SwimLogsApp.SaveSession(newSess)
	require.IsType(t, openapi.Session{}, s.Body())
	session := s.Body().(openapi.Session)

	res := SwimLogsApp.DeleteSessionById(session.Id)
	assert.Equal(t, http.StatusNoContent, res.Code())
}

// func TestUpdateSessionNotFound(t *testing.T) {
// 	it.TestFilter(t)
//
// 	newDay := openapi.Monday
// 	res := SwimLogsApp.UpdateSession(uuid.New(), openapi.UpdateSessionJSONBody{Day: &newDay})
// 	assert.Equal(t, http.StatusNotFound, res.Code())
// 	assert.IsType(t, openapi.ErrorDetail{}, res.Body())
// }
//
// func TestUpdateSessionSuccessfully(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
//
// 	newSess := openapi.CreateSessionJSONBody{
// 		Day:         openapi.Friday,
// 		StartTime:   openapi.StartTime("10:00"),
// 		DurationMin: 60,
// 	}
// 	sess, err := PostgresDbConn.SaveSession(newSess)
// 	require.Nil(t, err)
//
// 	newDay := openapi.Monday
// 	res := SwimLogsApp.UpdateSession(sess.Id, openapi.UpdateSessionJSONBody{Day: &newDay})
// 	assert.Equal(t, http.StatusOK, res.Code())
//
// 	require.IsType(t, openapi.Session{}, res.Body())
// 	session := res.Body().(openapi.Session)
// 	assert.Equal(t, newDay, session.Day)
// 	assert.Equal(t, newSess.StartTime, session.StartTime)
// 	assert.Equal(t, newSess.DurationMin, session.DurationMin)
// }
//
// func TestUpdateSessionInvalidDay(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
//
// 	newSess := openapi.CreateSessionJSONBody{
// 		Day:         openapi.Friday,
// 		StartTime:   openapi.StartTime("10:00"),
// 		DurationMin: 60,
// 	}
// 	sess, err := PostgresDbConn.SaveSession(newSess)
// 	require.Nil(t, err)
//
// 	newDay := openapi.Day("invalid")
// 	res := SwimLogsApp.UpdateSession(sess.Id, openapi.UpdateSessionJSONBody{Day: &newDay})
// 	assert.Equal(t, http.StatusBadRequest, res.Code())
// 	assert.IsType(t, openapi.ErrorDetail{}, res.Body())
// }
