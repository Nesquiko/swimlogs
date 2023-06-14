package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrainingsDetailsForCurrentWeekSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	req := newTraining()
	req.Date = types.Date{Time: time.Now()}
	var saved openapi.TrainingDetail
	{
		res := SwimLogsApp.SaveTraining(req)
		require.Equal(t, http.StatusCreated, res.Code())
		require.IsType(t, openapi.TrainingDetail{}, res.Body())
		saved = res.Body().(openapi.TrainingDetail)
	}

	res := SwimLogsApp.GetTrainingDetailsForCurrentWeek()
	require.Equal(t, http.StatusOK, res.Code())
	require.IsType(t, openapi.TrainingDetailsCurrentWeekResponse{}, res.Body())
	trainings := res.Body().(openapi.TrainingDetailsCurrentWeekResponse).Details
	require.Equal(t, 1, len(trainings))

	assert := assert.New(t)
	assert.Equal(saved.Id, trainings[0].Id)
	assert.Equal(saved.StartTime, trainings[0].StartTime)
	assert.Equal(saved.DurationMin, trainings[0].DurationMin)
	assert.Equal(saved.TotalDistance, trainings[0].TotalDistance)
}

func TestGetTrainingByIdSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	req := newTraining()
	res := SwimLogsApp.SaveTraining(req)
	require.Equal(t, http.StatusCreated, res.Code())

	require.IsType(t, openapi.TrainingDetail{}, res.Body())
	saved := res.Body().(openapi.TrainingDetail)

	trainingById := SwimLogsApp.GetTrainingById(saved.Id)
	assert := assert.New(t)
	assert.Equal(http.StatusOK, trainingById.Code())

	require.IsType(t, openapi.Training{}, trainingById.Body())
	training := trainingById.Body().(openapi.Training)
	assert.Equal(saved.StartTime, training.StartTime)
	assert.Equal(saved.DurationMin, training.DurationMin)
	assert.Equal(saved.TotalDistance, training.TotalDistance)
}

func TestGetTrainingByIdNotFound(t *testing.T) {
	it.TestFilter(t)

	res := SwimLogsApp.GetTrainingById(uuid.New())
	assert := assert.New(t)
	assert.Equal(http.StatusNotFound, res.Code())
}

func TestSaveTrainingSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	req := newTraining()

	res := SwimLogsApp.SaveTraining(req)
	assert := assert.New(t)
	assert.Equal(res.Code(), http.StatusCreated)

	require.IsType(t, openapi.TrainingDetail{}, res.Body())
	training := res.Body().(openapi.TrainingDetail)
	assert.Equal(req.Date.Unix(), training.Date.Unix())
	assert.Equal(req.StartTime, training.StartTime)
	assert.Equal(req.DurationMin, training.DurationMin)
	assert.Equal(req.TotalDistance, training.TotalDistance)
}

func newTraining() openapi.NewTraining {
	seconds := 30
	return openapi.NewTraining{
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
}
