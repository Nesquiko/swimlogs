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
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestEditTrainingSession(t *testing.T) {
	tests := []struct {
		name    string
		request apidef.EditSessionRequest
	}{
		{
			name:    "Duration",
			request: apidef.EditSessionRequest{DurationMin: asPtr(111)},
		},
		{
			name:    "Start",
			request: apidef.EditSessionRequest{Start: asPtr(time.Now())},
		},
	}

	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			ts := mustCreateNewTraining(t, nil)

			edited := mustEditTrainingSession(t, ts.Id, tC.request)

			assert := assert.New(t)
			if tC.request.Start != nil {
				assert.True(compareTimes(*tC.request.Start, edited.Start))
			} else {
				assert.True(compareTimes(ts.Start, edited.Start))
			}
			if tC.request.DurationMin != nil {
				assert.Equal(*tC.request.DurationMin, edited.DurationMin)
			} else {
				assert.Equal(ts.DurationMin, edited.DurationMin)
			}
		})
	}
}

func TestEditTrainingSession_Validation(t *testing.T) {
	testCases := []struct {
		name           string
		request        apidef.EditSessionRequest
		expectedCode   int
		expectedTitle  string
		expectedDetail string
	}{
		{
			name:           "Invalid Duration",
			request:        apidef.EditSessionRequest{DurationMin: asPtr(0)},
			expectedCode:   http.StatusBadRequest,
			expectedTitle:  app.InvalidEditSessionRequestTitle,
			expectedDetail: fmt.Sprintf(app.DurationErrorDetail, data.SmallIntMax, 0),
		},
		{
			name:           "Invalid Start",
			request:        apidef.EditSessionRequest{Start: &time.Time{}},
			expectedCode:   http.StatusBadRequest,
			expectedTitle:  app.InvalidEditSessionRequestTitle,
			expectedDetail: fmt.Sprintf(app.StartErrorDetail, &time.Time{}),
		},
		{
			name:           "No Change",
			request:        apidef.EditSessionRequest{},
			expectedCode:   http.StatusBadRequest,
			expectedTitle:  app.InvalidEditSessionRequestTitle,
			expectedDetail: app.NoSessionChangesDetail,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			ts := mustCreateNewTraining(t, nil)

			res, err := editTrainingSession(ts.Id, tC.request)
			require.NoError(t, err)
			require.Equal(t, tC.expectedCode, res.StatusCode, "response: %+v", res)

			var apiError apidef.InvalidSessionResponse
			err = json.NewDecoder(res.Body).Decode(&apiError)
			require.NoError(t, err)

			assert := assert.New(t)
			assert.Equal(tC.expectedTitle, apiError.Title)
			assert.Equal(app.InvalidEditSessionRequestCode, apiError.Code)
			assert.Equal(tC.expectedCode, apiError.Status)
			assert.Equal(tC.expectedDetail, apiError.Detail)
		})
	}
}

func mustEditTrainingSession(
	t *testing.T,
	id uuid.UUID,
	r apidef.EditSessionRequest,
) apidef.TrainingSummary {
	res, err := editTrainingSession(id, r)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode, "response: %+v", res)

	var ts apidef.TrainingSummary
	err = json.NewDecoder(res.Body).Decode(&ts)
	res.Body.Close()
	require.NoError(t, err, "response: %+v", res)

	return ts
}

func editTrainingSession(id uuid.UUID, r apidef.EditSessionRequest) (*http.Response, error) {
	url := ServerUrl + "/trainings/" + id.String()
	client := http.Client{}

	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("editTrainingSession marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("editTrainingSession request: %w", err)
	}
	req.Header.Add(server.ContentType, server.ApplicationJSON)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("editTrainingSession do: %w", err)
	}

	return res, nil
}
