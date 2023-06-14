package data

import (
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrainingDetailsForThisWeekSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	newTraining := newTraining()
	newTraining.Date = types.Date{Time: time.Now()}
	saved, err := PostgresDbConn.SaveTraining(newTraining)
	require.Nil(t, err)

	trainings, err := PostgresDbConn.GetTrainingDetailsForThisWeek()
	require.Nil(t, err)
	require.Equal(t, 1, len(trainings))

	assert := assert.New(t)
	assert.Equal(saved.Id, trainings[0].Id)
	assert.Equal(saved.StartTime, trainings[0].StartTime)
	assert.Equal(saved.DurationMin, trainings[0].DurationMin)
	assert.Equal(saved.TotalDistance, trainings[0].TotalDistance)
}

func TestGetTrainingByIdNotFound(t *testing.T) {
	it.TestFilter(t)

	_, err := PostgresDbConn.GetTrainingById(uuid.New())
	require.NotNil(t, err)
	assert.ErrorIs(t, err, data.ErrRowsNotFound)
}

func TestGetTrainingByIdSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	newTraining := newTraining()
	saved, err := PostgresDbConn.SaveTraining(newTraining)
	require.Nil(t, err)

	training, err := PostgresDbConn.GetTrainingById(saved.Id)
	require.Nil(t, err)

	assert := assert.New(t)
	assert.Equal(saved.Id, training.Id)
	assert.Equal(saved.StartTime, training.StartTime)
	assert.Equal(saved.DurationMin, training.DurationMin)
	assert.Equal(saved.TotalDistance, training.TotalDistance)
}

func TestSaveTrainingSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	newTraining := newTraining()

	_, err := PostgresDbConn.SaveTraining(newTraining)
	require.Nil(t, err)
	assert := assert.New(t)

	traininCount, err := data.SqlWithResult[int](
		PostgresDbConn.DB,
		"select count(*) from trainings",
	)
	require.Nil(t, err)
	assert.Equal(1, traininCount)

	blockCount, err := data.SqlWithResult[int](
		PostgresDbConn.DB,
		"select count(*) from blocks",
	)
	require.Nil(t, err)
	assert.Equal(2, blockCount)

	setCount, err := data.SqlWithResult[int](
		PostgresDbConn.DB,
		"select count(*) from sets",
	)
	require.Nil(t, err)
	assert.Equal(4, setCount)
}

func newTraining() openapi.NewTraining {
	seconds := 30
	newTraining := openapi.NewTraining{
		Date:          types.Date{Time: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)},
		StartTime:     "10:00",
		DurationMin:   60,
		TotalDistance: 2200,
		Blocks: []openapi.NewBlock{
			{
				Name:          "Warmup",
				Num:           0,
				Repeat:        1,
				TotalDistance: 1000,
				Sets: []openapi.NewTrainingSet{
					{
						Num:          0,
						Distance:     400,
						Repeat:       1,
						What:         "Freestyle",
						StartingRule: openapi.StartingRule{Type: openapi.None},
					},
					{
						Num:      1,
						Distance: 200,
						Repeat:   3,
						What:     "Breaststroke",
						StartingRule: openapi.StartingRule{
							Type:    openapi.Interval,
							Seconds: &seconds,
						},
					},
				},
			},
			{
				Name:          "Main",
				Num:           1,
				Repeat:        2,
				TotalDistance: 1200,
				Sets: []openapi.NewTrainingSet{
					{
						Num:      0,
						Distance: 50,
						Repeat:   10,
						What:     "50 fast, 50 slow",
						StartingRule: openapi.StartingRule{
							Type:    openapi.Interval,
							Seconds: &seconds,
						},
					},
					{
						Num:      1,
						Distance: 100,
						Repeat:   1,
						What:     "Cool down",
						StartingRule: openapi.StartingRule{
							Type: openapi.None,
						},
					},
				},
			},
		},
	}
	return newTraining
}
