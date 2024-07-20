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

func TestTrainingById_MatchingResponse(t *testing.T) {
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
	id := mustCreateNewTraining(t, &request).Id
	training := mustReadTraining(t, id)

	assert := assert.New(t)
	assert.Equal(request.DurationMin, training.DurationMin)
	assert.Equal(800, training.TotalDistance)
	assert.True(compareTimes(request.Start, training.Start))
	assert.Len(training.Sets, 2)

	setTotalDists := []int{400, 400}
	for i := range training.Sets {
		expectedSet := request.Sets[i]
		set := training.Sets[i]

		assert.Equal(setTotalDists[i], set.TotalDistance)
		assert.Equal(expectedSet.SetOrder, set.SetOrder)
		assert.Equal(expectedSet.Repeat, set.Repeat)
		assert.Equal(expectedSet.DistanceMeters, set.DistanceMeters)
		assert.Equal(expectedSet.Description, set.Description)
		assert.Equal(expectedSet.Equipment, set.Equipment)
		assert.Equal(expectedSet.Group, set.Group)
		assert.Equal(expectedSet.StartSeconds, set.StartSeconds)
		assert.Equal(expectedSet.StartType, set.StartType)
	}
}

func TestTrainingById_NotFound(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	res, err := readTraining(id)
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

func TestTrainingById_InvalidUUID(t *testing.T) {
	t.Parallel()
	invalidId := "invalid-uuid"
	res, err := _readTraining(invalidId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError apidef.NotFoundError
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(fmt.Sprintf(server.InvalidPathParamTitleFormat, "id", invalidId), apiError.Title)
	assert.Equal(server.InvalidPathParamCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(fmt.Sprintf(server.InvalidPathParamDetailFormat, "id", invalidId), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func mustReadTraining(t *testing.T, id uuid.UUID) apidef.Training {
	res, err := readTraining(id)
	defer res.Body.Close()
	require.NoError(t, err)

	if !assert.Equal(t, http.StatusOK, res.StatusCode) {
		if res.Header.Get(server.ContentType) == server.ApplicationProblemJSON {
			var apiErr apidef.ErrorDetail
			err := json.NewDecoder(res.Body).Decode(&apiErr)
			require.NoError(t, err)
			require.Fail(t, server.ApplicationProblemJSON, "error: %+v", apiErr)
		}
		require.Failf(t, "error", "res: %+v", res)
	}

	var tr apidef.Training
	err = json.NewDecoder(res.Body).Decode(&tr)
	require.NoError(t, err, "response: %+v", res)

	return tr
}

func readTraining(id uuid.UUID) (*http.Response, error) {
	return _readTraining(id.String())
}

func _readTraining(id string) (*http.Response, error) {
	url := ServerUrl + "/trainings/" + id
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("readTraining post: %w", err)
	}
	return res, nil
}
