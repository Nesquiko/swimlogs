package it

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/data"
)

func TestDeleteTraining_NotFound(t *testing.T) {
	client := http.Client{}
	url := TH.ts.URL + "/trainings/" + uuid.NewString()
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestDeleteTraining(t *testing.T) {
	tId := createTraining(t, nil).Id

	client := http.Client{}
	url := TH.ts.URL + "/trainings/" + tId.String()
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	result := struct {
		trainingCount int
		setCount      int
	}{}
	err = data.SqlWithResult(
		TH.pool,
		"select (select count(*) from trainings where id = $1), (select count(*) from sets where training_id = $1)",
		[]any{tId},
		[]any{&result.trainingCount, &result.setCount},
	)
	require.NoError(t, err)
	assert.Zero(t, result.trainingCount)
	assert.Zero(t, result.setCount)
}
