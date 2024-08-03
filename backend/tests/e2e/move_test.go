package e2e

import (
	"bytes"
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

func TestMoveSet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		initialSetOrder int
		newSetOrder     int

		initialNearSetOrder int
		newNearSetOrder     int
	}{
		{
			name:                "move up",
			initialSetOrder:     4,
			newSetOrder:         1,
			initialNearSetOrder: 3,
			newNearSetOrder:     4,
		},
		{
			name:                "move down",
			initialSetOrder:     1,
			newSetOrder:         5,
			initialNearSetOrder: 2,
			newNearSetOrder:     1,
		},
		{
			name:                "move to first",
			initialSetOrder:     5,
			newSetOrder:         0,
			initialNearSetOrder: 0,
			newNearSetOrder:     1,
		},
		{
			name:                "move to last",
			initialSetOrder:     0,
			newSetOrder:         5,
			initialNearSetOrder: 1,
			newNearSetOrder:     0,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()

			newTrainingReq := apidef.CreateTrainingRequest{
				DurationMin: 60,
				Sets: []apidef.NewTrainingSet{
					{SetOrder: 0, DistanceMeters: 1, Repeat: 1},
					{SetOrder: 1, DistanceMeters: 2, Repeat: 1},
					{SetOrder: 2, DistanceMeters: 3, Repeat: 1},
					{SetOrder: 3, DistanceMeters: 4, Repeat: 1},
					{SetOrder: 4, DistanceMeters: 5, Repeat: 1},
					{SetOrder: 5, DistanceMeters: 6, Repeat: 1},
				},
				Start: time.Now(),
			}
			trainingId := mustCreateNewTraining(t, &newTrainingReq).Id
			sets := mustReadTraining(t, trainingId).Sets
			setId := sets[tC.initialSetOrder].Id

			request := apidef.MoveSetRequest{NewSetOrder: tC.newSetOrder}
			training := mustMoveSet(t, trainingId, setId, request)

			assert := assert.New(t)
			assert.Len(training.Sets, len(sets))
			for i := 0; i < len(training.Sets); i++ {
				assert.Equal(i, training.Sets[i].SetOrder)
			}

			assert.Equal(setId, training.Sets[request.NewSetOrder].Id)
			assert.Equal(sets[tC.initialNearSetOrder].Id, training.Sets[tC.newNearSetOrder].Id)
		})
	}
}

func mustMoveSet(
	t *testing.T,
	trainingId uuid.UUID,
	setId uuid.UUID,
	request apidef.MoveSetRequest,
) apidef.Training {
	res, err := moveSet(trainingId, setId, request)
	defer res.Body.Close()
	require.NoError(t, err)

	if !assert.Equal(t, http.StatusOK, res.StatusCode) {
		if res.Header.Get(server.ContentType) == server.ApplicationProblemJSON {
			var apiErr apidef.ErrorDetail
			err := json.NewDecoder(res.Body).Decode(&apiErr)
			require.NoError(t, err)
			require.Failf(t, server.ApplicationProblemJSON, "error: %+v", apiErr)
		}
		require.Failf(t, "error", "res: %+v", res)
	}

	var tr apidef.Training
	err = json.NewDecoder(res.Body).Decode(&tr)
	require.NoErrorf(t, err, "response: %+v", res)

	return tr
}

func moveSet(
	trainingId uuid.UUID,
	setId uuid.UUID,
	request apidef.MoveSetRequest,
) (*http.Response, error) {
	url := ServerUrl + "/trainings/" + trainingId.String() + "/sets/move/" + setId.String()
	client := http.Client{}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("moveSet marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("moveSet request: %w", err)
	}
	req.Header.Add(server.ContentType, server.ApplicationJSON)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moveSet do: %w", err)
	}

	return res, nil
}
