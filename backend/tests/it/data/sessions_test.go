package data

import (
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveSessionSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := data.Session{
		Day:         "Monday",
		StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		DurationMin: 60,
	}

	sess, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(newSess.Day, sess.Day)
	assert.Equal(newSess.StartTime, sess.StartTime)
	assert.Equal(newSess.DurationMin, sess.DurationMin)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	require.Nil(t, err)
	assert.Equal(1, sessCount)
}

func TestSaveSessionInvalidDay(t *testing.T) {
	it.TestFilter(t)

	newSess := data.Session{
		Day:         "Invalid",
		StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		DurationMin: 60,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, data.ErrInvalidEnumType)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	require.Nil(t, err)
	assert.Equal(t, 0, sessCount)
}

func TestSaveSessionDurationLessThanZero(t *testing.T) {
	it.TestFilter(t)

	newSess := data.Session{
		Day:         "Monday",
		StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		DurationMin: -1,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, data.ErrCheckViolation)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	require.Nil(t, err)
	assert.Equal(t, 0, sessCount)
}

func TestSaveSessionDuplicateSession(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := data.Session{
		Day:         "Monday",
		StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		DurationMin: 60,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	require.Nil(t, err)

	_, err = PostgresDbConn.SaveSession(newSess)
	require.NotNil(t, err)
	assert.ErrorIs(t, err, data.ErrUniqueViolation)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	require.Nil(t, err)
	assert.Equal(t, 1, sessCount)
}

func TestGetAllSessions(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	for _, d := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		newSess := data.Session{
			Day:         d,
			StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
			DurationMin: 60,
		}
		_, err := PostgresDbConn.SaveSession(newSess)
		require.Nil(t, err)
	}

	assert := assert.New(t)
	sessions, total, err := PostgresDbConn.GetSessions(0, 5)
	assert.Nil(err)
	assert.Len(sessions, 5)
	assert.Equal(7, total)

	sessions, total, err = PostgresDbConn.GetSessions(1, 5)
	assert.Nil(err)
	assert.Len(sessions, 2)
	assert.Equal(7, total)
}

func TestDeleteSessionByIdUnknownId(t *testing.T) {
	it.TestFilter(t)
	err := PostgresDbConn.DeleteSessionById(uuid.New())
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, data.ErrRowsNotFound)
}

func TestDeleteSessionById(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	newSess := data.Session{
		Day:         "Tuesday",
		StartTime:   time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
		DurationMin: 60,
	}

	sess, err := PostgresDbConn.SaveSession(newSess)
	require.Nil(t, err)

	err = PostgresDbConn.DeleteSessionById(sess.Id)
	require.Nil(t, err)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	require.Nil(t, err)
	assert.Equal(t, 0, sessCount)
}

// func TestUpdateSessionSuccessfully(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
//
// 	newSess := data.Session{
// 		Day:         "Tuesday",
// 		StartTime:   time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
// 		DurationMin: 60,
// 	}
//
// 	sess, err := PostgresDbConn.SaveSession(newSess)
// 	require.Nil(t, err)
//
// 	newDay := openapi.Tuesday
// 	newStartTime := openapi.StartTime("11:00")
// 	updated := openapi.UpdateSessionJSONBody{
// 		Day:         &newDay,
// 		StartTime:   &newStartTime,
// 		DurationMin: nil,
// 	}
//
// 	sess, err = PostgresDbConn.UpdateSession(sess.Id, updated)
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
//
// 	assert.Equal(newDay, sess.Day)
// 	assert.Equal(newStartTime, sess.StartTime)
// 	assert.Equal(newSess.DurationMin, sess.DurationMin)
//
// 	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
// 	assert.Equal(1, sessCount)
// }
//
// func TestUpdateSessionDbErrorCases(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
// 	newSess := openapi.CreateSessionJSONBody{
// 		Day:         openapi.Monday,
// 		StartTime:   openapi.StartTime("10:00"),
// 		DurationMin: 60,
// 	}
// 	sess1, err := PostgresDbConn.SaveSession(newSess)
// 	assert := assert.New(t)
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
//
// 	invalidDay := openapi.Day("invalid")
// 	invalidDuration := -1
// 	testCases := []struct {
// 		desc    string
// 		updated openapi.UpdateSessionJSONBody
// 		err     error
// 	}{
// 		{
// 			desc: "invalid day",
// 			updated: openapi.UpdateSessionJSONBody{
// 				Day:         &invalidDay,
// 				StartTime:   nil,
// 				DurationMin: nil,
// 			},
// 			err: data.ErrInvalidEnumType,
// 		},
// 		{
// 			desc: "invalid duration",
// 			updated: openapi.UpdateSessionJSONBody{
// 				Day:         nil,
// 				StartTime:   nil,
// 				DurationMin: &invalidDuration,
// 			},
// 			err: data.ErrCheckViolation,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(*testing.T) {
// 			_, err := PostgresDbConn.UpdateSession(sess1.Id, tC.updated)
// 			assert.NotNil(err)
// 			assert.ErrorIs(err, tC.err)
// 		})
// 	}
// }
//
// func TestUpdateSessionDuplicateSession(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
//
// 	newSess := openapi.CreateSessionJSONBody{
// 		Day:         openapi.Monday,
// 		StartTime:   openapi.StartTime("10:00"),
// 		DurationMin: 60,
// 	}
// 	sess1, err := PostgresDbConn.SaveSession(newSess)
// 	assert := assert.New(t)
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
//
// 	newSess.Day = openapi.Tuesday
// 	sess2, err := PostgresDbConn.SaveSession(newSess)
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
//
// 	// change sess2 to be a duplicate of sess1
// 	updated := openapi.UpdateSessionJSONBody{
// 		Day:         &sess2.Day,
// 		StartTime:   nil,
// 		DurationMin: nil,
// 	}
//
// 	_, err = PostgresDbConn.UpdateSession(sess1.Id, updated)
// 	assert.NotNil(err)
// 	assert.ErrorIs(err, data.ErrUniqueViolation)
//
// 	sess2Day, err := data.SqlWithResult[string](
// 		PostgresDbConn.DB,
// 		"select day from sessions where id = $1",
// 		sess2.Id,
// 	)
// 	if !assert.Nil(err) {
// 		t.FailNow()
// 	}
// 	assert.Equal(string(sess2.Day), sess2Day)
// }
