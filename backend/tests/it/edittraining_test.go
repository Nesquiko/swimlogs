package it

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestEditTraining_NotFound(t *testing.T) {
	notSavedTraining := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				StartType:      apidef.None,
			},
		},
		Start: time.Now(),
	}
	id := createTraining(t, &notSavedTraining).Id
	training := trainingById(t, id)

	client := http.Client{}
	url := TH.ts.URL + "/trainings/" + uuid.NewString()
	req, err := json.Marshal(training)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(req))
	require.NoError(t, err)
	request.Header.Add("Content-Type", server.ApplicationJSON)
	res, err := client.Do(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestEditTraining(t *testing.T) {
	newTraining := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				StartType:      apidef.None,
			},
			{
				Description:    asPtr("some Description"),
				DistanceMeters: 50,
				Equipment:      &[]apidef.EquipmentEnum{apidef.Fins},
				Repeat:         8,
				SetOrder:       1,
				StartSeconds:   asPtr(90),
				StartType:      apidef.Interval,
			},
		},
		Start: time.Now(),
	}
	id := createTraining(t, &newTraining).Id
	training := trainingById(t, id)

	training.DurationMin = 120

	training.Sets[0].Repeat = 2
	training.Sets[0].DistanceMeters = 800
	training.Sets[0].Equipment = &[]apidef.EquipmentEnum{apidef.Board}

	training.Sets[1].Repeat = 4
	training.Sets[1].DistanceMeters = 200
	training.Sets[1].Equipment = nil

	training.Sets = append(training.Sets, apidef.TrainingSet{
		DistanceMeters: 50,
		Equipment:      &[]apidef.EquipmentEnum{apidef.Monofin, apidef.Snorkel},
		Id:             uuid.New(),
		Repeat:         5,
		SetOrder:       2,
		StartType:      apidef.None,
		TotalDistance:  250,
	})

	client := http.Client{}
	url := TH.ts.URL + "/trainings/" + training.Id.String()
	req, err := json.Marshal(training)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(req))
	require.NoError(t, err)
	request.Header.Add("Content-Type", server.ApplicationJSON)
	res, err := client.Do(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var response apidef.EditTrainingResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	res.Body.Close()
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(training.DurationMin, response.DurationMin)
	assert.Equal(2650, response.TotalDistance)
	assert.Equal(training.Start.Year(), response.Start.Year())
	assert.Equal(training.Start.Month(), response.Start.Month())
	assert.Equal(training.Start.Day(), response.Start.Day())
	assert.Equal(training.Start.Hour(), response.Start.Hour())
	assert.Equal(training.Start.Minute(), response.Start.Minute())

	var count int
	err = data.SqlWithResult(
		TH.pool,
		"select count(*) from sets where training_id = $1",
		[]any{training.Id},
		[]any{&count},
	)
	require.NoError(t, err)

	assert.Equal(3, count)

	result := struct {
		repeat         int
		distanceMeters int
		totalDistance  int
		equipment      []string
	}{}
	err = data.SqlWithResult(
		TH.pool,
		"select repeat, distance_meters, equipment, total_distance from sets where id = $1",
		[]any{training.Sets[0].Id},
		[]any{&result.repeat, &result.distanceMeters, &result.equipment, &result.totalDistance},
	)
	require.NoError(t, err)

	assert.Equal(2, result.repeat)
	assert.Equal(800, result.distanceMeters)
	assert.Equal(1600, result.totalDistance)
	assert.Equal(string(apidef.Board), result.equipment[0])

	err = data.SqlWithResult(
		TH.pool,
		"select repeat, distance_meters, equipment, total_distance from sets where id = $1",
		[]any{training.Sets[1].Id},
		[]any{&result.repeat, &result.distanceMeters, &result.equipment, &result.totalDistance},
	)
	require.NoError(t, err)

	assert.Equal(4, result.repeat)
	assert.Equal(200, result.distanceMeters)
	assert.Equal(800, result.totalDistance)
	assert.Nil(result.equipment)

	err = data.SqlWithResult(
		TH.pool,
		"select repeat, distance_meters, equipment, total_distance from sets where id != $1 and id != $2 and training_id = $3",
		[]any{
			training.Sets[0].Id,
			training.Sets[1].Id,
			training.Id,
		}, // I don't have the Id of the newly created set
		[]any{&result.repeat, &result.distanceMeters, &result.equipment, &result.totalDistance},
	)
	require.NoError(t, err)

	assert.Equal(5, result.repeat)
	assert.Equal(50, result.distanceMeters)
	assert.Equal(250, result.totalDistance)
	assert.Equal(string(apidef.Monofin), result.equipment[0])
	assert.Equal(string(apidef.Snorkel), result.equipment[1])
}
