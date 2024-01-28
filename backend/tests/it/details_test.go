package it

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
)

func TestTrainingDetails_Pagination(t *testing.T) {
	TH.CleanTrainings(t)
	for i := 0; i < 5; i++ {
		createTraining(t, nil)
	}

	url := fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 0, 2)
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var details apidef.TrainingDetailsResponse
	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, 0, details.Pagination.Page)
	assert.Equal(t, 2, details.Pagination.PageSize)
	assert.Equal(t, 5, details.Pagination.Total)

	url = fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 1, 2)
	res, err = http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, 1, details.Pagination.Page)
	assert.Equal(t, 2, details.Pagination.PageSize)
	assert.Equal(t, 5, details.Pagination.Total)
}

func TestTrainingDetails_Paging(t *testing.T) {
	TH.CleanTrainings(t)
	trainingIds := make([]uuid.UUID, 5)
	for i := 0; i < len(trainingIds); i++ {
		trainingIds[i] = createTraining(t, nil).Id
	}

	url := fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 0, 2)
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var details apidef.TrainingDetailsResponse
	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, trainingIds[0], details.Details[0].Id)
	assert.Equal(t, trainingIds[1], details.Details[1].Id)

	url = fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 1, 2)
	res, err = http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	url = fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 2, 2)
	assert.Equal(t, trainingIds[2], details.Details[0].Id)
	assert.Equal(t, trainingIds[3], details.Details[1].Id)

	res, err = http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, trainingIds[4], details.Details[0].Id)
}

func TestTrainingDetails(t *testing.T) {
	TH.CleanTrainings(t)
	trainingsCount := 5
	for i := 0; i < trainingsCount; i++ {
		createTraining(t, nil)
	}

	url := fmt.Sprintf("%s/trainings/details?page=%d&pageSize=%d", TH.ts.URL, 0, trainingsCount)
	res, err := http.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var details apidef.TrainingDetailsResponse
	err = json.NewDecoder(res.Body).Decode(&details)
	res.Body.Close()
	require.NoError(t, err)

	assert.Len(t, details.Details, trainingsCount)
}
