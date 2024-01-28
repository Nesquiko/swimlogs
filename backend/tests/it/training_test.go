package it

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
)

func TestTraining_NotFound(t *testing.T) {
	url := TH.ts.URL + "/trainings/" + uuid.NewString()
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestTraining(t *testing.T) {
	newTraining := apidef.CreateTrainingRequest{
		DurationMin: 120,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				StartType:      apidef.None,
				TotalDistance:  400,
			},
			{
				Description:    asPtr("some Description"),
				DistanceMeters: 50,
				Equipment:      &[]apidef.EquipmentEnum{apidef.Fins},
				Repeat:         8,
				SetOrder:       1,
				StartSeconds:   asPtr(90),
				StartType:      apidef.Interval,
				TotalDistance:  400,
			},
		},
		Start:         time.Now(),
		TotalDistance: 800,
	}
	id := createTraining(t, &newTraining).Id

	url := TH.ts.URL + "/trainings/" + id.String()
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var training apidef.Training
	err = json.NewDecoder(res.Body).Decode(&training)
	res.Body.Close()
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(id, training.Id)
	assert.Equal(newTraining.DurationMin, training.DurationMin)
	assert.Equal(newTraining.TotalDistance, training.TotalDistance)
	assert.Equal(newTraining.Start.Year(), training.Start.Year())
	assert.Equal(newTraining.Start.Month(), training.Start.Month())
	assert.Equal(newTraining.Start.Day(), training.Start.Day())
	assert.Equal(newTraining.Start.Hour(), training.Start.Hour())
	assert.Equal(newTraining.Start.Minute(), training.Start.Minute())
	assert.Len(training.Sets, 2)

	for i := 0; i < len(newTraining.Sets); i++ {
		newSet := newTraining.Sets[i]
		set := training.Sets[i]

		assert.Equal(newSet.DistanceMeters, set.DistanceMeters)
		assert.Equal(newSet.TotalDistance, set.TotalDistance)
		assert.Equal(newSet.Repeat, set.Repeat)
		assert.Equal(newSet.SetOrder, set.SetOrder)
		assert.Equal(newSet.Description, set.Description)
		assert.Equal(newSet.Equipment, set.Equipment)
		assert.Equal(newSet.StartSeconds, set.StartSeconds)
		assert.Equal(newSet.StartType, set.StartType)
	}
}

func trainingById(t *testing.T, id uuid.UUID) apidef.Training {
	url := TH.ts.URL + "/trainings/" + id.String()
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var training apidef.Training
	err = json.NewDecoder(res.Body).Decode(&training)
	res.Body.Close()
	require.NoError(t, err)

	return training
}
