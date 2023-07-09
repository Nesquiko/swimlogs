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

func TestGetTrainingDetailsForWeekSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	weekTime := time.Date(2023, 8, 24, 12, 0, 0, 0, time.UTC) // Thursday
	for i := 0; i < 7; i++ {
		newTraining := newTraining()
		newTraining.Start = weekTime.AddDate(0, 0, -i)
		_, err := PostgresDbConn.SaveTraining(newTraining)
		require.Nil(t, err)
	}

	expectedCount := 4
	trainings, err := PostgresDbConn.GetTrainingDetailsInWeek(weekTime)
	require.Nil(t, err)
	assert.Equal(t, expectedCount, len(trainings))
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
	assert.Equal(saved.Start, training.Start)
	assert.Equal(saved.DurationMin, training.DurationMin)
	assert.Equal(saved.TotalDistance, training.TotalDistance)
	assert.Equal(len(saved.Sets), len(training.Sets))
}

func TestSaveTrainingSuccessfully(t *testing.T) {
	it.TestFilter(t)
	t.Cleanup(func() { it.TruncateTrainings(PostgresDbConn.DB) })

	newTraining := newTraining()
	expectedSetCount := len(newTraining.Sets)

	_, err := PostgresDbConn.SaveTraining(newTraining)
	require.Nil(t, err)
	assert := assert.New(t)

	traininCount, err := data.SqlWithResult[int](
		PostgresDbConn.DB,
		"select count(*) from trainings",
	)
	require.Nil(t, err)
	assert.Equal(1, traininCount)

	setCount, err := data.SqlWithResult[int](
		PostgresDbConn.DB,
		"select count(*) from sets",
	)
	require.Nil(t, err)
	assert.Equal(expectedSetCount, setCount)
}

func newTraining() data.Training {
	trainingId := uuid.New()
	superSet1Id := uuid.New()
	superSet2Id := uuid.New()

	return data.Training{
		Id:            trainingId,
		Start:         time.Now(),
		DurationMin:   90,
		TotalDistance: 1900,
		Sets: []data.TrainingSet{
			{
				Id:             uuid.New(),
				TrainingId:     trainingId,
				SetOrder:       0,
				TotalDistance:  400,
				Repeat:         1,
				DistanceMeters: asPtr(400),
				Description:    asPtr("warm up freestyle"),
				StartType:      data.NoneStartType,
			},
			{
				Id:             uuid.New(),
				TrainingId:     trainingId,
				SetOrder:       1,
				TotalDistance:  600,
				Repeat:         3,
				DistanceMeters: asPtr(200),
				Description:    asPtr("drills"),
				StartType:      data.PauseStartType,
				StartSeconds:   asPtr(20),
			},
			{
				Id:            superSet1Id,
				TrainingId:    trainingId,
				SetOrder:      2,
				TotalDistance: 300,
				Repeat:        4,
				StartType:     data.NoneStartType,
				SubSets: &[]data.TrainingSet{
					{
						Id:             uuid.New(),
						TrainingId:     trainingId,
						SubSetOrder:    asPtr(0),
						TotalDistance:  50,
						Repeat:         1,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartType:      data.NoneStartType,
						StartSeconds:   asPtr(60),
					},
					{
						Id:             uuid.New(),
						TrainingId:     trainingId,
						SubSetOrder:    asPtr(1),
						TotalDistance:  25,
						Repeat:         1,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartType:      data.PauseStartType,
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				Id:             uuid.New(),
				TrainingId:     trainingId,
				SetOrder:       3,
				TotalDistance:  100,
				Repeat:         1,
				DistanceMeters: asPtr(100),
				Description:    asPtr("cool down"),
				StartType:      data.NoneStartType,
			},
			{
				Id:            superSet2Id,
				TrainingId:    trainingId,
				SetOrder:      4,
				TotalDistance: 300,
				Repeat:        4,
				StartType:     data.NoneStartType,
				SubSets: &[]data.TrainingSet{
					{
						Id:             uuid.New(),
						TrainingId:     trainingId,
						SubSetOrder:    asPtr(0),
						TotalDistance:  50,
						Repeat:         1,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartType:      data.NoneStartType,
						StartSeconds:   asPtr(60),
					},
					{
						Id:             uuid.New(),
						TrainingId:     trainingId,
						SubSetOrder:    asPtr(1),
						TotalDistance:  25,
						Repeat:         1,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartType:      data.PauseStartType,
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				Id:             uuid.New(),
				TrainingId:     trainingId,
				SetOrder:       5,
				TotalDistance:  200,
				Repeat:         1,
				DistanceMeters: asPtr(200),
				Description:    asPtr("breastroke cool down"),
				StartType:      data.NoneStartType,
			},
		},
	}
}

func asPtr[T any](v T) *T {
	return &v
}
