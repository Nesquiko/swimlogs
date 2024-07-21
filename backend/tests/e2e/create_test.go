//go:build integration

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestCreateTraining_ValidTrainingSummary(t *testing.T) {
	t.Parallel()
	request := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				Group:          asPtr(apidef.Mono),
			},
			{
				Description:    asPtr("some Description"),
				DistanceMeters: 50,
				Equipment:      &[]apidef.EquipmentEnum{apidef.Fins},
				Repeat:         8,
				SetOrder:       1,
				StartSeconds:   asPtr(90),
				StartType:      asPtr(apidef.Interval),
			},
		},
		Start: time.Now(),
	}
	ts := mustCreateNewTraining(t, &request)

	assert := assert.New(t)
	assert.Equal(request.DurationMin, ts.DurationMin)
	assert.Equal(800, ts.TotalDistance)
	assert.True(compareTimes(request.Start, ts.Start))
}

func TestCreateTraining_NonUniqueSetOrder(t *testing.T) {
	t.Parallel()
	nonUniqueSetOrder := 0
	request := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				SetOrder:       nonUniqueSetOrder,
				Repeat:         1,
				DistanceMeters: 400,
			},
			{
				SetOrder:       nonUniqueSetOrder,
				Repeat:         8,
				DistanceMeters: 50,
			},
		},
		Start: time.Now(),
	}

	res, err := createNewTraining(&request)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError apidef.InvalidTrainingError
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(app.InvalidSetErrorTitle, apiError.Title)
	assert.Equal(app.NonUniqueSetOrderCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(fmt.Sprintf(app.NonUniqueSetOrderDetail, nonUniqueSetOrder), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func TestCreateTraining_InvalidSet(t *testing.T) {
	t.Parallel()
	invalidStartType := apidef.StartTypeEnum("invalid")
	request := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				SetOrder:       0,
				Repeat:         1,
				DistanceMeters: 400,
			},
			{
				SetOrder:       1,
				Repeat:         1,
				DistanceMeters: 50,
				StartType:      &invalidStartType,
			},
		},
		Start: time.Now(),
	}

	res, err := createNewTraining(&request)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError apidef.InvalidSetResponse
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(app.InvalidSetErrorTitle, apiError.Title)
	assert.Equal(app.InvalidSetErrorCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(
		fmt.Sprintf(
			app.StartTypeUnknownErrorDetail,
			apidef.Interval,
			apidef.Pause,
			invalidStartType,
		),
		apiError.Detail,
	)
	assert.NotNil(apiError.AdditionalProperties)
	assert.Equal(request.Sets[1].SetOrder, int(apiError.AdditionalProperties["setOrder"].(float64)))
}

func mustCreateNewTraining(
	t *testing.T,
	request *apidef.CreateTrainingRequest,
) apidef.TrainingSummary {
	res, err := createNewTraining(request)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode, "response: %+v", res)

	var ts apidef.TrainingSummary
	err = json.NewDecoder(res.Body).Decode(&ts)
	res.Body.Close()
	require.NoError(t, err, "response: %+v", res)

	return ts
}

func createNewTraining(
	request *apidef.CreateTrainingRequest,
) (*http.Response, error) {
	if request == nil {
		request = defaultNewTraining()
	}
	req, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("createNewTraining marshal: %w", err)
	}

	url := ServerUrl + "/trainings"
	res, err := http.Post(url, server.ApplicationJSON, bytes.NewBuffer(req))
	if err != nil {
		return nil, fmt.Errorf("createNewTraining post: %w", err)
	}
	return res, nil
}

func defaultNewTraining() *apidef.CreateTrainingRequest {
	return &apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{{
			SetOrder:       0,
			Repeat:         4,
			DistanceMeters: 100,
			Description:    asPtr("Some description"),
			Equipment:      &[]apidef.EquipmentEnum{apidef.Board, apidef.Fins},
			Group:          asPtr(apidef.Long),
			StartSeconds:   asPtr(60),
			StartType:      asPtr(apidef.Interval),
		}},
		Start: time.Date(2024, 6, 24, 18, 0, 0, 0, time.UTC),
	}
}
