package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

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
