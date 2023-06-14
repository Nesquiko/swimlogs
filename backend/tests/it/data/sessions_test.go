package data

import (
	"testing"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveSessionSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}

	sess, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(newSess.Day, sess.Day)
	assert.Equal(newSess.StartTime, sess.StartTime)
	assert.Equal(newSess.DurationMin, sess.DurationMin)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(1, sessCount)
}

func TestSaveSessionInvalidDay(t *testing.T) {
	it.TestFilter(t)
	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Day("invalid"),
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.NotNil(err)
	assert.ErrorIs(err, data.ErrInvalidEnumType)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(0, sessCount)
}

func TestSaveSessionDurationLessThanZero(t *testing.T) {
	it.TestFilter(t)
	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: -1,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.NotNil(err)
	assert.ErrorIs(err, data.ErrCheckViolation)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(0, sessCount)
}

func TestSaveSessionDuplicateSession(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}

	_, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.Nil(err)

	_, err = PostgresDbConn.SaveSession(newSess)
	assert.NotNil(err)
	assert.ErrorIs(err, data.ErrUniqueViolation)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(1, sessCount)
}

func TestGetAllSessions(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	for _, d := range []openapi.Day{openapi.Monday, openapi.Tuesday, openapi.Wednesday, openapi.Thursday, openapi.Friday, openapi.Saturday, openapi.Sunday} {
		newSess := openapi.CreateSessionJSONBody{
			Day:         d,
			StartTime:   openapi.StartTime("10:00"),
			DurationMin: 60,
		}
		_, err := PostgresDbConn.SaveSession(newSess)
		assert.Nil(t, err)
	}

	sessions, err := PostgresDbConn.GetAllSessions()
	assert.Nil(t, err)
	assert.Len(t, sessions, 7)
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
	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}

	sess, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	assert.Nil(err)

	err = PostgresDbConn.DeleteSessionById(sess.Id)
	if !assert.Nil(err) {
		t.FailNow()
	}

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(0, sessCount)
}

func TestUpdateSessionSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}

	sess, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	if !assert.Nil(err) {
		t.FailNow()
	}

	newDay := openapi.Tuesday
	newStartTime := openapi.StartTime("11:00")
	updated := openapi.UpdateSessionJSONBody{
		Day:         &newDay,
		StartTime:   &newStartTime,
		DurationMin: nil,
	}

	sess, err = PostgresDbConn.UpdateSession(sess.Id, updated)
	if !assert.Nil(err) {
		t.FailNow()
	}

	assert.Equal(newDay, sess.Day)
	assert.Equal(newStartTime, sess.StartTime)
	assert.Equal(newSess.DurationMin, sess.DurationMin)

	sessCount, err := data.SqlWithResult[int](PostgresDbConn.DB, "select count(*) from sessions")
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(1, sessCount)
}

func TestUpdateSessionDbErrorCases(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })
	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}
	sess1, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	if !assert.Nil(err) {
		t.FailNow()
	}

	invalidDay := openapi.Day("invalid")
	invalidDuration := -1
	testCases := []struct {
		desc    string
		updated openapi.UpdateSessionJSONBody
		err     error
	}{
		{
			desc: "invalid day",
			updated: openapi.UpdateSessionJSONBody{
				Day:         &invalidDay,
				StartTime:   nil,
				DurationMin: nil,
			},
			err: data.ErrInvalidEnumType,
		},
		{
			desc: "invalid duration",
			updated: openapi.UpdateSessionJSONBody{
				Day:         nil,
				StartTime:   nil,
				DurationMin: &invalidDuration,
			},
			err: data.ErrCheckViolation,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(*testing.T) {
			_, err := PostgresDbConn.UpdateSession(sess1.Id, tC.updated)
			assert.NotNil(err)
			assert.ErrorIs(err, tC.err)
		})
	}
}

func TestUpdateSessionDuplicateSession(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateSessions(PostgresDbConn.DB) })

	newSess := openapi.CreateSessionJSONBody{
		Day:         openapi.Monday,
		StartTime:   openapi.StartTime("10:00"),
		DurationMin: 60,
	}
	sess1, err := PostgresDbConn.SaveSession(newSess)
	assert := assert.New(t)
	if !assert.Nil(err) {
		t.FailNow()
	}

	newSess.Day = openapi.Tuesday
	sess2, err := PostgresDbConn.SaveSession(newSess)
	if !assert.Nil(err) {
		t.FailNow()
	}

	// change sess2 to be a duplicate of sess1
	updated := openapi.UpdateSessionJSONBody{
		Day:         &sess2.Day,
		StartTime:   nil,
		DurationMin: nil,
	}

	_, err = PostgresDbConn.UpdateSession(sess1.Id, updated)
	assert.NotNil(err)
	assert.ErrorIs(err, data.ErrUniqueViolation)

	sess2Day, err := data.SqlWithResult[string](
		PostgresDbConn.DB,
		"select day from sessions where id = $1",
		sess2.Id,
	)
	if !assert.Nil(err) {
		t.FailNow()
	}
	assert.Equal(string(sess2.Day), sess2Day)
}
