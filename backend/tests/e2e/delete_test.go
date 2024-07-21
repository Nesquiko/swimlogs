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

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestDeleteSet_Successfully(t *testing.T) {
	t.Parallel()

	request := apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
			},
			{
				DistanceMeters: 50,
				Repeat:         8,
				SetOrder:       1,
			},
		},
		Start: time.Now(),
	}
	trainingId := mustCreateNewTraining(t, &request).Id
	training := mustReadTraining(t, trainingId)
	setId := training.Sets[0].Id

	res, err := deleteSet(trainingId, setId)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	training = mustReadTraining(t, trainingId)
	assert.Len(t, training.Sets, 1)
	assert.Equal(t, training.Sets[0].SetOrder, 1)
}

func TestDeleteSet_AlsoDeleteTraining(t *testing.T) {
	t.Parallel()

	trainingId := mustCreateNewTraining(t, nil).Id
	training := mustReadTraining(t, trainingId)
	setId := training.Sets[0].Id

	res, err := deleteSet(trainingId, setId)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	res, err = readTraining(trainingId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError apidef.NotFoundError
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
	trainingId := uuid.New()
	setId := uuid.New()
	res, err := deleteSet(trainingId, setId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var apiError apidef.NotFoundError
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

	var apiError apidef.NotFoundError
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

	var apiError apidef.NotFoundError
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

func deleteSet(trainingId, setId uuid.UUID) (*http.Response, error) {
	url := ServerUrl + "/trainings/" + trainingId.String() + "/sets/" + setId.String()
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
