package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestEditSet(t *testing.T) {
	t.Parallel()

	start := api.Start{}
	start.FromInterval(api.Interval{
		Type:    api.IntervalTypeInterval,
		Seconds: 90,
	})
	newTrainingRequest := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				Group:          nullable.NewNullableWithValue(api.Mono),
				Type:           api.Normal,
			},
			{
				Description:    asPtr("some Description"),
				DistanceMeters: 50,
				Equipment:      &[]api.EquipmentEnum{api.Fins},
				Repeat:         8,
				SetOrder:       1,
				Start:          &start,
				IsMain:         asPtr(true),
				Type:           api.Normal,
				Progression:    nullable.NewNullableWithValue(api.Desc),
			},
		},
		Start: time.Now(),
	}
	ts := mustCreateNewTraining(t, &newTrainingRequest)
	training := mustReadTraining(t, ts.Id)
	set := training.Sets[0]

	newDescription := "new description"
	newDistance := 75
	newRepeat := 7
	newEquipment := []api.EquipmentEnum{api.Monofin, api.Fins, api.Board}
	newGroup := api.Sprint
	newStart := api.Start{}
	newStartType := api.PauseTypePause
	newStartSeconds := 45
	err := newStart.FromPause(api.Pause{
		Type:    newStartType,
		Seconds: newStartSeconds,
	})
	require.NoError(t, err)
	newIsMain := false
	newType := api.Compound
	newIntensity := api.En3
	request := api.EditSetRequest{
		Description:    nullable.NewNullableWithValue(newDescription),
		Repeat:         &newRepeat,
		DistanceMeters: &newDistance,
		Equipment:      &newEquipment,
		Group:          nullable.NewNullableWithValue(newGroup),
		Start:          nullable.NewNullableWithValue(newStart),
		IsMain:         &newIsMain,
		Type:           &newType,
		Intensity:      nullable.NewNullableWithValue(newIntensity),
		Progression:    nullable.NewNullNullable[api.ProgressionEnum](),
	}

	edited := mustEditSet(t, set.Id, request)
	editedSet := edited.Sets[0]
	assert := assert.New(t)
	require := require.New(t)

	assert.Equal(ts.TotalDistance-400+newDistance*newRepeat, edited.TotalDistance)
	assert.Equal(set.Id, editedSet.Id)
	assert.Equal(newDescription, *editedSet.Description)
	assert.Equal(newRepeat, editedSet.Repeat)
	assert.Equal(newDistance, editedSet.DistanceMeters)
	assert.Equal(newEquipment, *editedSet.Equipment)
	assert.Equal(newGroup, editedSet.Group.MustGet())

	newStartDist, err := newStart.Discriminator()
	require.NoError(err)
	require.Equal(string(newStartType), newStartDist)
	editedPause, err := newStart.AsPause()
	require.NoError(err)

	assert.Equal(newStartType, editedPause.Type)
	assert.Equal(newStartSeconds, editedPause.Seconds)

	assert.Equal(newIsMain, editedSet.IsMain)
	assert.Equal(newType, editedSet.Type)
	assert.Equal(newIntensity, editedSet.Intensity.MustGet())
	assert.False(editedSet.Progression.IsSpecified())
}

func TestEditSet_DontChangeNullableToNull(t *testing.T) {
	t.Parallel()

	group := api.Mono
	newTrainingRequest := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{
				Group:          nullable.NewNullableWithValue(group),
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				Type:           api.Normal,
			},
		},
		Start: time.Now(),
	}
	ts := mustCreateNewTraining(t, &newTrainingRequest)
	training := mustReadTraining(t, ts.Id)
	set := training.Sets[0]

	newDistance := 75
	request := api.EditSetRequest{DistanceMeters: &newDistance}

	edited := mustEditSet(t, set.Id, request)
	editedSet := edited.Sets[0]
	assert := assert.New(t)

	assert.Equal(set.Id, editedSet.Id)
	assert.Equal(newDistance, editedSet.DistanceMeters)
	require.NotNil(t, editedSet.Group)
	assert.Equal(group, editedSet.Group.MustGet())
}

func TestEditSet_ChangeNullableToNull(t *testing.T) {
	t.Parallel()

	group := api.Mono
	newTrainingRequest := api.CreateTrainingRequest{
		DurationMinutes: 60,
		Sets: []api.NewTrainingSet{
			{
				Group:          nullable.NewNullableWithValue(group),
				DistanceMeters: 400,
				Repeat:         1,
				SetOrder:       0,
				Type:           api.Normal,
			},
		},
		Start: time.Now(),
	}
	ts := mustCreateNewTraining(t, &newTrainingRequest)
	training := mustReadTraining(t, ts.Id)
	set := training.Sets[0]

	request := api.EditSetRequest{Group: nullable.NewNullNullable[api.GroupEnum]()}

	edited := mustEditSet(t, set.Id, request)
	editedSet := edited.Sets[0]
	assert := assert.New(t)

	assert.Equal(set.Id, editedSet.Id)
	assert.Nil(editedSet.Group)
}

func TestEditSet_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		request api.EditSetRequest
	}{
		{
			name:    "Invalid Repeat",
			request: api.EditSetRequest{Repeat: asPtr(0)},
		},
		{
			name:    "Invalid Distance",
			request: api.EditSetRequest{DistanceMeters: asPtr(data.SmallIntMax + 1)},
		},
		{
			name: "Unknown Equipment",
			request: api.EditSetRequest{
				Equipment: &[]api.EquipmentEnum{api.EquipmentEnum("unknown")},
			},
		},
		{
			name: "Unknown Group",
			request: api.EditSetRequest{
				Group: nullable.NewNullableWithValue(api.GroupEnum("unknown")),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			ts := mustCreateNewTraining(t, nil)
			training := mustReadTraining(t, ts.Id)
			set := training.Sets[0]

			res, err := editSet(set.Id, tC.request)
			require.NoError(t, err)
			require.Equalf(t, http.StatusBadRequest, res.StatusCode, "response: %+v", res)

			var apiError api.ErrorDetail
			err = json.NewDecoder(res.Body).Decode(&apiError)
			require.NoError(t, err)

			assert := assert.New(t)
			assert.Equal(server.SchemaValidationErrorCode, apiError.Code)
			assert.NotNil(apiError.AdditionalProperties)
		})
	}
}

func TestEditSet_NotFound(t *testing.T) {
	t.Parallel()
	setId := uuid.New()

	res, err := editSet(
		setId,
		api.EditSetRequest{Description: nullable.NewNullableWithValue("some desc")},
	)
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

func mustEditSet(t *testing.T, setId uuid.UUID, r api.EditSetRequest) api.Training {
	res, err := editSet(setId, r)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)

	var training api.Training
	err = json.NewDecoder(res.Body).Decode(&training)
	res.Body.Close()
	require.NoErrorf(t, err, "response: %+v", res)

	return training
}

func editSet(id uuid.UUID, r api.EditSetRequest) (*http.Response, error) {
	url := ServerUrl + "/sets/" + id.String()
	client := http.Client{}

	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("editSet marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("editSet request: %w", err)
	}
	req.Header.Add(server.ContentType, server.ApplicationJSON)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("editSet do: %w", err)
	}

	return res, nil
}
