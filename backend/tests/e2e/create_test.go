//go:build integration

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestCreateTraining_ValidTraining(t *testing.T) {
	t.Parallel()

	request := defaultRequest(t)
	training := mustCreateNewTraining(t, request)

	assert := assert.New(t)
	assert.Equal(request.DurationMinutes, training.DurationMinutes)
	assert.Equal(4600, training.TotalDistance)
	assert.True(compareTimes(request.Start, training.Start))
	assert.Len(training.Sets, 10)
}

func TestCreateTraining_NonUniqueSetOrder(t *testing.T) {
	t.Parallel()
	nonUniqueSetOrder := 0
	request := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{
				SetOrder:       nonUniqueSetOrder,
				Repeat:         1,
				DistanceMeters: 400,
				Type:           api.Normal,
			},
			{
				SetOrder:       nonUniqueSetOrder,
				Repeat:         8,
				DistanceMeters: 50,
				Type:           api.Normal,
			},
		},
		Start: time.Now(),
	}

	res, err := createNewTraining(t, &request)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError api.ErrorDetail
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
	invalidSetOrder := -1
	request := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{
				SetOrder:       0,
				Repeat:         1,
				DistanceMeters: 400,
				Type:           api.Normal,
			},
			{
				SetOrder:       invalidSetOrder,
				Repeat:         1,
				DistanceMeters: 400,
				Type:           api.Normal,
			},
		},
		Start: time.Now(),
	}

	res, err := createNewTraining(t, &request)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(server.SchemaValidationErrorCode, apiError.Code)
	assert.Equal(server.ValidationErrorTitle, apiError.Title)
	assert.Equal(
		fmt.Sprintf(server.ValidationErrorDetail, "number must be at least 0"),
		apiError.Detail,
	)
	assert.Equal(http.StatusBadRequest, apiError.Status)

	assert.NotNil(apiError.AdditionalProperties)
	path := apiError.AdditionalProperties["path"]
	reason := apiError.AdditionalProperties["reason"]
	schema := apiError.AdditionalProperties["schema"]
	assert.Equal("/sets/1/setOrder", path)
	assert.Equal("number must be at least 0", reason)
	assert.Equal("#/components/schemas/NewTraining", schema)
}

func mustCreateNewTraining(
	t *testing.T,
	request *api.CreateTrainingRequest,
) api.Training {
	res, err := createNewTraining(t, request)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode, "response: %+v", res)

	var ts api.Training
	err = json.NewDecoder(res.Body).Decode(&ts)
	res.Body.Close()
	require.NoError(t, err, "response: %+v", res)

	return ts
}

func createNewTraining(t *testing.T, request *api.CreateTrainingRequest) (*http.Response, error) {
	if request == nil {
		request = defaultRequest(t)
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

func defaultRequest(t *testing.T) *api.CreateTrainingRequest {
	jsonContents, err := os.ReadFile("../data/default_new_training.json")
	require.NoError(t, err)

	var request api.CreateTrainingRequest
	err = json.Unmarshal(jsonContents, &request)
	require.NoError(t, err)

	return &request
}
