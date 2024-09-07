package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestDeleteSet_Successfully(t *testing.T) {
	t.Parallel()

	request := defaultRequest(t)
	trainingId := mustCreateNewTraining(t, request).Id
	training := mustReadTraining(t, trainingId)
	setId := training.Sets[0].Id

	res, err := deleteSet(setId)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	training = mustReadTraining(t, trainingId)
	assert.Len(t, training.Sets, len(request.Sets)-1)
	assert.Equal(t, training.Sets[0].SetOrder, 0)
}

func TestDeleteSet_ReorderRemainingSets(t *testing.T) {
	t.Parallel()

	request := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{SetOrder: 0, DistanceMeters: 1, Repeat: 1, Type: api.Normal},
			{SetOrder: 1, DistanceMeters: 2, Repeat: 1, Type: api.Normal},
			{SetOrder: 2, DistanceMeters: 3, Repeat: 1, Type: api.Normal},
			{SetOrder: 3, DistanceMeters: 4, Repeat: 1, Type: api.Normal},
			{SetOrder: 4, DistanceMeters: 5, Repeat: 1, Type: api.Normal},
		},
		Start: time.Now(),
	}
	trainingId := mustCreateNewTraining(t, &request).Id
	training := mustReadTraining(t, trainingId)
	setId := training.Sets[1].Id

	res, err := deleteSet(setId)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	training = mustReadTraining(t, trainingId)
	assert.Len(t, training.Sets, 4)
	for i := 0; i < 4; i++ {
		assert.Equal(t, i, training.Sets[i].SetOrder)
	}
}

func TestDeleteSet_AlsoDeleteTraining(t *testing.T) {
	t.Parallel()

	trainingId := mustCreateNewTraining(t, nil).Id
	training := mustReadTraining(t, trainingId)
	for _, set := range training.Sets {
		setId := set.Id
		res, err := deleteSet(setId)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	}

	res, err := readTraining(trainingId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(fmt.Sprintf(server.NotFoundTitleFormat, "training"), apiError.Title)
	assert.Equal(server.NotFoundCode, apiError.Code)
	assert.Equal(http.StatusNotFound, apiError.Status)
	assert.Equal(
		fmt.Sprintf(server.NotFoundDetailFormat, "training", trainingId.String()),
		apiError.Detail,
	)
	assert.Nil(apiError.AdditionalProperties)
}

func TestDeleteSet_NotFound(t *testing.T) {
	t.Parallel()
	setId := uuid.New()
	res, err := deleteSet(setId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(fmt.Sprintf(server.NotFoundTitleFormat, "set"), apiError.Title)
	assert.Equal(server.NotFoundCode, apiError.Code)
	assert.Equal(http.StatusNotFound, apiError.Status)
	assert.Equal(fmt.Sprintf(server.NotFoundDetailFormat, "set", setId.String()), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func TestDeleteTraining_Successfully(t *testing.T) {
	t.Parallel()
	id := mustCreateNewTraining(t, nil).Id
	res, err := deleteTraining(id)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = readTraining(id)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(fmt.Sprintf(server.NotFoundTitleFormat, "training"), apiError.Title)
	assert.Equal(server.NotFoundCode, apiError.Code)
	assert.Equal(http.StatusNotFound, apiError.Status)
	assert.Equal(fmt.Sprintf(server.NotFoundDetailFormat, "training", id.String()), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func TestDeleteTraining_NotFound(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	res, err := deleteTraining(id)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(fmt.Sprintf(server.NotFoundTitleFormat, "training"), apiError.Title)
	assert.Equal(server.NotFoundCode, apiError.Code)
	assert.Equal(http.StatusNotFound, apiError.Status)
	assert.Equal(fmt.Sprintf(server.NotFoundDetailFormat, "training", id.String()), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func deleteTraining(id uuid.UUID) (*http.Response, error) {
	url := ServerUrl + "/trainings/" + id.String()
	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("deleteTraining request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deleteTraining do: %w", err)
	}

	return res, nil
}

func deleteSet(id uuid.UUID) (*http.Response, error) {
	url := ServerUrl + "/sets/" + id.String()
	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("deleteSet request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deleteSet do: %w", err)
	}

	return res, nil
}
