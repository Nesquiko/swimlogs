package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestGetTrainingsDetailsForCurrentWeekSuccessfully(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })
//
// 	req := newTraining()
// 	req.Date = types.Date{Time: time.Now()}
// 	var saved openapi.TrainingDetail
// 	{
// 		res := SwimLogsApp.SaveTraining(req)
// 		require.Equal(t, http.StatusCreated, res.Code())
// 		require.IsType(t, openapi.TrainingDetail{}, res.Body())
// 		saved = res.Body().(openapi.TrainingDetail)
// 	}
//
// 	res := SwimLogsApp.GetTrainingDetailsForCurrentWeek()
// 	require.Equal(t, http.StatusOK, res.Code())
// 	require.IsType(t, openapi.TrainingDetailsCurrentWeekResponse{}, res.Body())
// 	trainings := res.Body().(openapi.TrainingDetailsCurrentWeekResponse).Details
// 	require.Equal(t, 1, len(trainings))
//
// 	assert := assert.New(t)
// 	assert.Equal(saved.Id, trainings[0].Id)
// 	assert.Equal(saved.StartTime, trainings[0].StartTime)
// 	assert.Equal(saved.DurationMin, trainings[0].DurationMin)
// 	assert.Equal(saved.TotalDistance, trainings[0].TotalDistance)
// }

// func TestGetTrainingByIdSuccessfully(t *testing.T) {
// 	it.TestFilter(t)
// 	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })
//
// 	req := newTraining()
// 	res := SwimLogsApp.SaveTraining(req)
// 	require.Equal(t, http.StatusCreated, res.Code())
//
// 	require.IsType(t, openapi.TrainingDetail{}, res.Body())
// 	saved := res.Body().(openapi.TrainingDetail)
//
// 	trainingById := SwimLogsApp.GetTrainingById(saved.Id)
// 	assert := assert.New(t)
// 	assert.Equal(http.StatusOK, trainingById.Code())
//
// 	require.IsType(t, openapi.Training{}, trainingById.Body())
// 	training := trainingById.Body().(openapi.Training)
// 	assert.Equal(saved.StartTime, training.StartTime)
// 	assert.Equal(saved.DurationMin, training.DurationMin)
// 	assert.Equal(saved.TotalDistance, training.TotalDistance)
// }

// func TestGetTrainingByIdNotFound(t *testing.T) {
// 	it.TestFilter(t)
//
// 	res := SwimLogsApp.GetTrainingById(uuid.New())
// 	assert := assert.New(t)
// 	assert.Equal(http.StatusNotFound, res.Code())
// }

func TestSaveTrainingSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	req := newTraining()

	res := SwimLogsApp.SaveTraining(req)
	assert := assert.New(t)
	assert.Equal(res.Code(), http.StatusCreated)

	require.IsType(t, openapi.TrainingDetail{}, res.Body())
	training := res.Body().(openapi.TrainingDetail)
	assert.True(req.Start.Truncate(time.Second).Equal(training.Start.Truncate(time.Second)))
	assert.Equal(req.DurationMin, training.DurationMin)
	assert.Equal(req.TotalDistance, training.TotalDistance)
}

func newTraining() openapi.NewTraining {
	return openapi.NewTraining{
		Start:         time.Now(),
		DurationMin:   90,
		TotalDistance: 1900,
		Sets: []openapi.NewTrainingSet{
			{
				SetOrder:       asPtr(0),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  400,
				DistanceMeters: asPtr(400),
				Description:    asPtr("warm up freestyle"),
			},
			{
				SetOrder:       asPtr(1),
				Repeat:         3,
				StartType:      openapi.Pause,
				TotalDistance:  600,
				DistanceMeters: asPtr(200),
				Description:    asPtr("drills"),
				StartSeconds:   asPtr(20),
			},
			{
				SetOrder:      asPtr(2),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 300,
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  50,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  25,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(3),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  100,
				DistanceMeters: asPtr(100),
				Description:    asPtr("cool down"),
			},
			{
				SetOrder:      asPtr(4),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 300,
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  50,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  25,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(5),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  200,
				DistanceMeters: asPtr(200),
				Description:    asPtr("breastroke cool down"),
			},
		},
	}
}

func asPtr[T any](v T) *T {
	return &v
}
