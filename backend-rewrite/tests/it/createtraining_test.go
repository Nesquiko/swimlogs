package it

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestCreateTraining_ResponseMatches(t *testing.T) {
	request := apidef.CreateTrainingRequest{
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
	req, err := json.Marshal(request)
	require.NoError(t, err)

	url := TH.ts.URL + "/trainings"
	res, err := http.Post(url, server.ApplicationJSON, bytes.NewBuffer(req))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var trainingDetail apidef.CreateTraining201JSONResponse
	err = json.NewDecoder(res.Body).Decode(&trainingDetail)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, request.DurationMin, trainingDetail.DurationMin)
	assert.Equal(t, 800, trainingDetail.TotalDistance)
	assert.Equal(t, request.Start.Year(), trainingDetail.Start.Year())
	assert.Equal(t, request.Start.Month(), trainingDetail.Start.Month())
	assert.Equal(t, request.Start.Day(), trainingDetail.Start.Day())
	assert.Equal(t, request.Start.Hour(), trainingDetail.Start.Hour())
	assert.Equal(t, request.Start.Minute(), trainingDetail.Start.Minute())

	result := struct {
		trainingCount int
		setCount      int
	}{}
	err = data.SqlWithResult(
		TH.pool,
		"select (select count(*) from trainings where id = $1), (select count(*) from sets where training_id = $1)",
		[]any{trainingDetail.Id},
		[]any{&result.trainingCount, &result.setCount},
	)
	require.NoError(t, err)
	assert.Equal(t, 1, result.trainingCount)
	assert.Equal(t, 2, result.setCount)
}

func createTraining(t *testing.T, training *apidef.CreateTrainingRequest) apidef.TrainingDetail {
	var request apidef.CreateTrainingRequest
	if training == nil {
		request = apidef.CreateTrainingRequest{
			DurationMin: 60,
			Sets: []apidef.NewTrainingSet{
				{
					DistanceMeters: 100,
					Repeat:         1,
					SetOrder:       0,
					StartType:      apidef.None,
				},
			},
			Start: time.Now(),
		}
	} else {
		request = *training
	}
	req, err := json.Marshal(request)
	require.NoError(t, err)

	url := TH.ts.URL + "/trainings"
	res, err := http.Post(url, server.ApplicationJSON, bytes.NewBuffer(req))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var trainingDeatail apidef.CreateTraining201JSONResponse
	err = json.NewDecoder(res.Body).Decode(&trainingDeatail)
	res.Body.Close()
	require.NoError(t, err)

	return apidef.TrainingDetail(trainingDeatail.CreateTrainingReponseJSONResponse)
}
