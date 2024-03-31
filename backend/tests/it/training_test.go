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
	exptectedTraining := apidef.CreateTrainingRequest{
		DurationMin: 120,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				StartType:      apidef.None,
				TotalDistance:  400,
				Group:          asPtr(apidef.Bifi),
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
	id := createTraining(t, &exptectedTraining).Id

	url := TH.ts.URL + "/trainings/" + id.String()
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var actualTraining apidef.Training
	err = json.NewDecoder(res.Body).Decode(&actualTraining)
	res.Body.Close()
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(id, actualTraining.Id)
	assert.Equal(exptectedTraining.DurationMin, actualTraining.DurationMin)
	assert.Equal(exptectedTraining.TotalDistance, actualTraining.TotalDistance)
	assert.Equal(exptectedTraining.Start.Year(), actualTraining.Start.Year())
	assert.Equal(exptectedTraining.Start.Month(), actualTraining.Start.Month())
	assert.Equal(exptectedTraining.Start.Day(), actualTraining.Start.Day())
	assert.Equal(exptectedTraining.Start.Hour(), actualTraining.Start.Hour())
	assert.Equal(exptectedTraining.Start.Minute(), actualTraining.Start.Minute())
	assert.Len(actualTraining.Sets, 2)

	for i := 0; i < len(exptectedTraining.Sets); i++ {
		exptectedSet := exptectedTraining.Sets[i]
		actualSet := actualTraining.Sets[i]

		assert.Equal(exptectedSet.DistanceMeters, actualSet.DistanceMeters)
		assert.Equal(exptectedSet.TotalDistance, actualSet.TotalDistance)
		assert.Equal(exptectedSet.Repeat, actualSet.Repeat)
		assert.Equal(exptectedSet.SetOrder, actualSet.SetOrder)
		assert.Equal(exptectedSet.Description, actualSet.Description)
		assert.Equal(exptectedSet.Equipment, actualSet.Equipment)
		assert.Equal(exptectedSet.StartSeconds, actualSet.StartSeconds)
		assert.Equal(exptectedSet.StartType, actualSet.StartType)
		assert.Equal(exptectedSet.Group, actualSet.Group)
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
