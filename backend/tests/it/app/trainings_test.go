package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/Nesquiko/swimlogs/tests/it"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrainingsDetailsSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	weekTime := time.Date(2023, 8, 24, 12, 0, 0, 0, time.UTC) // Thursday
	for i := 0; i < 7; i++ {
		newTraining := newTraining()
		newTraining.Start = weekTime.AddDate(0, 0, -i)
		res := SwimLogsApp.SaveTraining(newTraining)
		require.Equal(t, http.StatusCreated, res.Code())
		require.IsType(t, openapi.TrainingDetail{}, res.Body())
	}

	expectedTotalCount := 7
	res := SwimLogsApp.GetTrainingDetails(openapi.GetTrainingsDetailsParams{Page: 0, PageSize: 5})
	require.Equal(t, http.StatusOK, res.Code())
	require.IsType(t, openapi.TrainingDetailsResponse{}, res.Body())
	response := res.Body().(openapi.TrainingDetailsResponse)

	assert := assert.New(t)
	assert.Equal(expectedTotalCount, response.Pagination.Total)
	assert.Equal(5, len(response.Details))

	res = SwimLogsApp.GetTrainingDetails(openapi.GetTrainingsDetailsParams{Page: 1, PageSize: 5})
	require.Equal(t, http.StatusOK, res.Code())
	require.IsType(t, openapi.TrainingDetailsResponse{}, res.Body())
	response = res.Body().(openapi.TrainingDetailsResponse)
	assert.Equal(expectedTotalCount, response.Pagination.Total)
	assert.Equal(2, len(response.Details))
}

func TestGetTrainingsDetailsForCurrentWeekSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	weekTime := time.Date(2023, 8, 24, 12, 0, 0, 0, time.UTC) // Thursday
	for i := 0; i < 7; i++ {
		newTraining := newTraining()
		newTraining.Start = weekTime.AddDate(0, 0, -i)
		res := SwimLogsApp.SaveTraining(newTraining)
		require.Equal(t, http.StatusCreated, res.Code())
		require.IsType(t, openapi.TrainingDetail{}, res.Body())
	}

	expectedCount := 4
	res := SwimLogsApp.GetTrainingDetailsInWeek(weekTime)
	require.Equal(t, http.StatusOK, res.Code())
	require.IsType(t, openapi.TrainingDetailsCurrentWeekResponse{}, res.Body())
	trainings := res.Body().(openapi.TrainingDetailsCurrentWeekResponse).Details
	assert.Equal(t, expectedCount, len(trainings))
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
	assert.True(saved.Start.Truncate(time.Second).Equal(training.Start.Truncate(time.Second)))
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
