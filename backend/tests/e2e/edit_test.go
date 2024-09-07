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

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

//	func TestEditSet(t *testing.T) {
//		t.Parallel()
//
//		newTrainingRequest := api.CreateTrainingRequest{
//			DurationMin: 60,
//			Sets: []api.NewTrainingSet{
//				{
//					DistanceMeters: 400,
//					Repeat:         1,
//					SetOrder:       0,
//					Group:          asPtr(api.Mono),
//				},
//				{
//					Description:    asPtr("some Description"),
//					DistanceMeters: 50,
//					Equipment:      &[]api.EquipmentEnum{api.Fins},
//					Repeat:         8,
//					SetOrder:       1,
//					StartSeconds:   asPtr(90),
//					StartType:      asPtr(api.Interval),
//					IsMain:         asPtr(true),
//				},
//			},
//			Start: time.Now(),
//		}
//		ts := mustCreateNewTraining(t, &newTrainingRequest)
//		training := mustReadTraining(t, ts.Id)
//		set := training.Sets[0]
//
//		newDescription := "new description"
//		newDistance := 75
//		newRepeat := 7
//		newEquipment := []api.EquipmentEnum{api.Monofin, api.Fins, api.Board}
//		newGroup := api.Sprint
//		newStartType := api.Pause
//		newStartSeconds := 45
//		newIsMain := false
//		request := api.EditSetRequest{
//			Description:    &newDescription,
//			Repeat:         &newRepeat,
//			DistanceMeters: &newDistance,
//			Equipment:      &newEquipment,
//			Group:          &newGroup,
//			StartType:      &newStartType,
//			StartSeconds:   &newStartSeconds,
//			IsMain:         &newIsMain,
//		}
//
//		res := mustEditSet(t, set.Id, request)
//		editedSet := res.Set
//		assert := assert.New(t)
//
//		assert.Equal(ts.TotalDistance-400+newDistance*newRepeat, res.TotalDistance)
//		assert.Equal(newDescription, *editedSet.Description)
//		assert.Equal(newDistance, editedSet.DistanceMeters)
//		assert.Equal(newRepeat, editedSet.Repeat)
//		assert.Equal(newEquipment, *editedSet.Equipment)
//		assert.Equal(newGroup, *editedSet.Group)
//		assert.Equal(newStartType, *editedSet.StartType)
//		assert.Equal(newStartSeconds, *editedSet.StartSeconds)
//		assert.Equal(newIsMain, editedSet.IsMain)
//	}
//
//	func TestEditSet_Validation(t *testing.T) {
//		t.Parallel()
//		testCases := []struct {
//			name           string
//			request        api.EditSetRequest
//			expectedTitle  string
//			expectedDetail string
//		}{
//			{
//				name:           "Invalid Repeat",
//				request:        api.EditSetRequest{Repeat: asPtr(0)},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(app.RepeatErrorDetail, data.SmallIntMax, 0),
//			},
//			{
//				name:          "Invalid Distance",
//				request:       api.EditSetRequest{DistanceMeters: asPtr(data.SmallIntMax + 1)},
//				expectedTitle: app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(
//					app.DistanceErrorDetail,
//					data.SmallIntMax,
//					data.SmallIntMax+1,
//				),
//			},
//			{
//				name: "Unknown StartType",
//				request: api.EditSetRequest{
//					StartType: (*api.StartTypeEnum)(asPtr("unknown")),
//				},
//				expectedTitle: app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(
//					app.StartTypeUnknownErrorDetail,
//					api.Interval,
//					api.Pause,
//					"unknown",
//				),
//			},
//			{
//				name:           "Missing StartSeconds",
//				request:        api.EditSetRequest{StartType: asPtr(api.Interval)},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: app.StartSecondsRequiredErrorDetail,
//			},
//			{
//				name: "Invalid StartSeconds",
//				request: api.EditSetRequest{
//					StartType:    asPtr(api.Pause),
//					StartSeconds: asPtr(0),
//				},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(app.StartSecondsErrorDetail, data.SmallIntMax, 0),
//			},
//			{
//				name: "Unknown Equipment",
//				request: api.EditSetRequest{
//					Equipment: &[]api.EquipmentEnum{api.EquipmentEnum("unknown")},
//				},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(app.UnknownEnumErrorDetail, app.Equipment, "unknown"),
//			},
//			{
//				name: "Unknown Group",
//				request: api.EditSetRequest{
//					Group: (*api.GroupEnum)(asPtr("unknown")),
//				},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: fmt.Sprintf(app.UnknownEnumErrorDetail, app.Group, "unknown"),
//			},
//			{
//				name:           "No Change",
//				request:        api.EditSetRequest{},
//				expectedTitle:  app.InvalidSetErrorTitle,
//				expectedDetail: app.NoSetChangesDetail,
//			},
//		}
//
//		for _, tC := range testCases {
//			t.Run(tC.name, func(t *testing.T) {
//				ts := mustCreateNewTraining(t, nil)
//				training := mustReadTraining(t, ts.Id)
//				set := training.Sets[0]
//
//				res, err := editSet(set.Id, tC.request)
//				require.NoError(t, err)
//				require.Equalf(t, http.StatusBadRequest, res.StatusCode, "response: %+v", res)
//
//				var apiError api.InvalidSetError
//				err = json.NewDecoder(res.Body).Decode(&apiError)
//				require.NoError(t, err)
//
//				assert := assert.New(t)
//				assert.Equal(tC.expectedTitle, apiError.Title)
//				assert.Equal(app.InvalidSetErrorCode, apiError.Code)
//				assert.Equal(http.StatusBadRequest, apiError.Status)
//				assert.Equal(tC.expectedDetail, apiError.Detail)
//			})
//		}
//	}
//
//	func TestEditSet_NotFound(t *testing.T) {
//		t.Parallel()
//		setId := uuid.New()
//
//		res, err := editSet(setId, api.EditSetRequest{Description: asPtr("some desc")})
//		defer res.Body.Close()
//		require.NoError(t, err)
//		require.Equal(t, http.StatusNotFound, res.StatusCode)
//
//		var apiError api.NotFoundError
//		err = json.NewDecoder(res.Body).Decode(&apiError)
//		require.NoError(t, err)
//
//		assert := assert.New(t)
//		assert.Equal(fmt.Sprintf(server.NotFoundTitleFormat, "set"), apiError.Title)
//		assert.Equal(server.NotFoundCode, apiError.Code)
//		assert.Equal(http.StatusNotFound, apiError.Status)
//		assert.Equal(fmt.Sprintf(server.NotFoundDetailFormat, "set", setId.String()), apiError.Detail)
//		assert.Nil(apiError.AdditionalProperties)
//	}

func TestEditTrainingSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request api.EditSessionRequest
	}{
		{
			name:    "Duration",
			request: api.EditSessionRequest{DurationMinutes: asPtr(111)},
		},
		{
			name:    "Start",
			request: api.EditSessionRequest{Start: asPtr(time.Now())},
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
			if tC.request.DurationMinutes != nil {
				assert.Equal(*tC.request.DurationMinutes, edited.DurationMinutes)
			} else {
				assert.Equal(ts.DurationMinutes, edited.DurationMinutes)
			}
		})
	}
}

func TestEditTrainingSession_NoChange(t *testing.T) {
	t.Parallel()

	request := api.EditSessionRequest{}
	ts := mustCreateNewTraining(t, nil)

	res, err := editTrainingSession(ts.Id, request)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotModified, res.StatusCode, "response: %+v", res)
}

func TestEditTrainingSession_NotFound(t *testing.T) {
	t.Parallel()
	id := uuid.New()

	start := time.Now()
	res, err := editTrainingSession(id, api.EditSessionRequest{Start: &start})
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

func mustEditTrainingSession(
	t *testing.T,
	id uuid.UUID,
	r api.EditSessionRequest,
) api.TrainingSummary {
	res, err := editTrainingSession(id, r)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)

	var ts api.TrainingSummary
	err = json.NewDecoder(res.Body).Decode(&ts)
	res.Body.Close()
	require.NoErrorf(t, err, "response: %+v", res)

	return ts
}

func editTrainingSession(id uuid.UUID, r api.EditSessionRequest) (*http.Response, error) {
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

//
// func mustEditSet(t *testing.T, setId uuid.UUID, r api.EditSetRequest) api.EditSetResponse {
// 	res, err := editSet(setId, r)
// 	require.NoError(t, err)
// 	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)
//
// 	var esr api.EditSetResponse
// 	err = json.NewDecoder(res.Body).Decode(&esr)
// 	res.Body.Close()
// 	require.NoErrorf(t, err, "response: %+v", res)
//
// 	return esr
// }
//
// func editSet(id uuid.UUID, r api.EditSetRequest) (*http.Response, error) {
// 	url := ServerUrl + "/sets/" + id.String()
// 	client := http.Client{}
//
// 	body, err := json.Marshal(r)
// 	if err != nil {
// 		return nil, fmt.Errorf("editSet marshal: %w", err)
// 	}
//
// 	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, fmt.Errorf("editSet request: %w", err)
// 	}
// 	req.Header.Add(server.ContentType, server.ApplicationJSON)
//
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("editSet do: %w", err)
// 	}
//
// 	return res, nil
// }
