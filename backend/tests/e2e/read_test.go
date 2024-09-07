package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestTrainingById_MatchingResponse(t *testing.T) {
	t.Parallel()
	request := defaultRequest(t)
	id := mustCreateNewTraining(t, request).Id
	training := mustReadTraining(t, id)

	assert := assert.New(t)
	require := require.New(t)
	assert.Equal(request.DurationMinutes, training.DurationMinutes)
	assert.Equal(4600, training.TotalDistance)
	assert.True(compareTimes(request.Start, training.Start))
	assert.Len(training.Sets, 10)

	setTotalDists := []int{400, 12 * 50, 8 * 50, 16 * 50, 100, 5 * 100, 100, 5 * 100, 100, 1100}
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

		if expectedSet.Start != nil {
			require.NotNilf(
				set.Start,
				"required set %d to have start",
				expectedSet.SetOrder,
				expectedSet.Start,
			)
			expectedDisc, _ := expectedSet.Start.Discriminator()
			disc, _ := set.Start.Discriminator()
			assert.Equal(expectedDisc, disc)

			switch disc {
			case string(api.LastFinishes):
			case string(api.Interval), string(api.Pause):
			default:
				assert.Failf("unknown start type: %q", disc)
			}
		}

		if expectedSet.IsMain != nil {
			assert.Equal(*expectedSet.IsMain, set.IsMain)
		}

		if expectedSet.Components != nil {
			require.NotNil(set.Components)
			assert.Equal(len(*expectedSet.Components), len(*set.Components))
		}
	}
}

func TestTrainingById_NotFound(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	res, err := readTraining(id)
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

func TestTrainingById_InvalidUUID(t *testing.T) {
	t.Parallel()
	invalidId := "invalid-uuid"
	res, err := _readTraining(invalidId)
	defer res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(server.InvalidParamErrorTitle, apiError.Title)
	assert.Equal(server.InvalidParamErrorCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(fmt.Sprintf(server.InvalidParamErrorDetail, "id", invalidId), apiError.Detail)
	assert.Nil(apiError.AdditionalProperties)
}

func mustReadTraining(t *testing.T, id uuid.UUID) api.Training {
	res, err := readTraining(id)
	defer res.Body.Close()
	require.NoError(t, err)

	if !assert.Equal(t, http.StatusOK, res.StatusCode) {
		if res.Header.Get(server.ContentType) == server.ApplicationProblemJSON {
			var apiErr api.ErrorDetail
			err := json.NewDecoder(res.Body).Decode(&apiErr)
			require.NoError(t, err)
			require.Fail(t, server.ApplicationProblemJSON, "error: %+v", apiErr)
		}
		require.Failf(t, "error", "res: %+v", res)
	}

	var tr api.Training
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
