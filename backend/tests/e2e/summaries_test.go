package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestSummariesPage_ReturnsMainSets(t *testing.T) {
	ClearDB()

	n := 10
	startStr := "2020-05-11T13:00:00+01:00"
	for i := 0; i < n; i++ {
		now, err := time.Parse(time.RFC3339, startStr)
		require.NoError(t, err)
		now = now.Add(time.Duration(i) * time.Minute)

		request := defaultRequest(t)
		request.Sets = []api.NewTrainingSet{request.Sets[0]}
		request.Sets[0].IsMain = asPtr(i%2 == 0)
		request.Sets = append(request.Sets, api.NewTrainingSet{
			SetOrder:       1,
			Repeat:         4,
			DistanceMeters: 100,
			IsMain:         asPtr(i%2 == 0),
			Type:           api.Normal,
		})
		request.Start = now

		mustCreateNewTraining(t, request)
	}

	from, err := time.Parse(time.RFC3339, "2020-05-11T12:59:00+01:00")
	require.NoError(t, err)
	until := from.Add(time.Duration(n) * time.Minute)

	assert := assert.New(t)
	summaries := mustReadSummariesPage(t, 0, n, &from, &until)
	for _, summary := range summaries.Summaries {
		if summary.MainSets == nil || len(*summary.MainSets) == 0 {
			continue
		}
		assert.Len(*summary.MainSets, 2)
	}
}

func TestSummariesPage_CorrectStartFiltering(t *testing.T) {
	ClearDB()

	n := 30
	startStr := "2020-08-11T09:00:00+01:00"
	start, err := time.Parse(time.RFC3339, startStr)
	require.NoError(t, err)

	created := make([]api.Training, 0, n)
	for i := 0; i < n; i++ {
		trainingStart := start.Add(time.Duration(i) * time.Minute)

		request := defaultRequest(t)
		request.Start = trainingStart

		training := mustCreateNewTraining(t, request)
		created = append(created, training)
	}

	expectedSummaries := 10
	from := start.Add(10 * time.Minute)
	require.NoError(t, err)
	until := from.Add(time.Duration(expectedSummaries-1) * time.Minute)

	pageSize := 5
	summaries := mustReadSummariesPage(t, 0, pageSize, &from, &until)
	assert := assert.New(t)
	assert.Len(summaries.Summaries, pageSize)
	for i := 0; i < pageSize; i++ {
		assert.Equalf(
			created[2*expectedSummaries-1-i].Id,
			summaries.Summaries[i].Id,
			"expected start: %s, actual start: %s",
			created[2*expectedSummaries-1-i].Start,
			summaries.Summaries[i].Start,
		)
	}
	assert.Equal(0, summaries.Pagination.Page)
	assert.Equal(pageSize, summaries.Pagination.PageSize)
	assert.Equal(n, summaries.Pagination.Total)

	summaries = mustReadSummariesPage(t, 1, pageSize, &from, &until)
	assert.Len(summaries.Summaries, pageSize)
	for i := pageSize; i < 2*pageSize; i++ {
		assert.Equalf(
			created[2*expectedSummaries-1-i].Id,
			summaries.Summaries[i-pageSize].Id,
			"expected start: %s, actual start: %s",
			created[2*expectedSummaries-1-i].Start,
			summaries.Summaries[i-pageSize].Start,
		)
	}
	assert.Equal(1, summaries.Pagination.Page)
	assert.Equal(pageSize, summaries.Pagination.PageSize)
	assert.Equal(n, summaries.Pagination.Total)

	summaries = mustReadSummariesPage(t, 0, n, &from, nil)
	assert.Len(summaries.Summaries, 20)
	assert.Equal(0, summaries.Pagination.Page)
	assert.Equal(20, summaries.Pagination.PageSize)
	assert.Equal(n, summaries.Pagination.Total)
}

func TestSummariesPage_CorrectPaging(t *testing.T) {
	ClearDB()

	n := 10
	for i := 0; i < n; i++ {
		mustCreateNewTraining(t, nil)
	}

	assert := assert.New(t)
	pageSize := 5

	summaries := mustReadSummariesPage(t, 0, pageSize, nil, nil)
	assert.Len(summaries.Summaries, pageSize)
	assert.Equal(0, summaries.Pagination.Page)
	assert.Equal(pageSize, summaries.Pagination.PageSize)
	assert.GreaterOrEqual(summaries.Pagination.Total, n)

	total := summaries.Pagination.Total
	summaries = mustReadSummariesPage(t, 1, total, nil, nil)
	assert.Empty(summaries.Summaries)
	assert.Equal(1, summaries.Pagination.Page)
	assert.Equal(0, summaries.Pagination.PageSize)

	summaries = mustReadSummariesPage(t, 1, total-1, nil, nil)
	assert.Len(summaries.Summaries, 1)
	assert.Equal(1, summaries.Pagination.Page)
	assert.Equal(1, summaries.Pagination.PageSize)
}

func TestSummariesPage_ReadPage(t *testing.T) {
	ClearDB()

	created := make([]api.Training, 0)
	for i := 0; i < 10; i++ {
		now := time.Now()
		now = now.AddDate(1, 0, i)

		request := defaultRequest(t)
		request.Start = now

		training := mustCreateNewTraining(t, request)
		created = append(created, training)
	}

	summaries := mustReadSummariesPage(t, 0, len(created), nil, nil)

	assert := assert.New(t)
	assert.Len(summaries.Summaries, len(created))

outer:

	for _, ts := range summaries.Summaries {
		for _, expected := range created {
			if expected.Id != ts.Id {
				continue
			}

			assert.Equal(expected.Id, ts.Id)
			continue outer
		}

		require.Failf(
			t,
			"no matching id",
			"There was no summary with matching id %s",
			ts.Id.String(),
		)
	}
}

func TestSummariesPage_InvalidPageSize(t *testing.T) {
	invalidPageSize := 0
	res, err := readSummariesPage(0, invalidPageSize, nil, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError api.ErrorDetail
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(server.ValidationErrorTitle, apiError.Title)
	assert.Equal(server.SchemaValidationErrorCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(
		fmt.Sprintf(server.ValidationErrorDetail, "number must be at least 1"),
		apiError.Detail,
	)

	require.NotNil(t, apiError.AdditionalProperties)
	path := apiError.AdditionalProperties["path"]
	reason := apiError.AdditionalProperties["reason"]
	assert.Equal("\"pageSize\"", path)
	assert.Equal("number must be at least 1", reason)
}

func mustReadSummariesPage(
	t *testing.T,
	page int,
	pageSize int,
	from *time.Time,
	until *time.Time,
) api.TrainingSummariesResponse {
	res, err := readSummariesPage(page, pageSize, from, until)

	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)

	var summaries api.TrainingSummariesResponse
	err = json.NewDecoder(res.Body).Decode(&summaries)
	res.Body.Close()
	require.NoErrorf(t, err, "response: %+v", res)

	return summaries
}

func readSummariesPage(
	page int,
	pageSize int,
	from *time.Time,
	until *time.Time,
) (*http.Response, error) {
	url := ServerUrl + "/trainings/summaries"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("readSummariesPage request: %w", err)
	}
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("pageSize", strconv.Itoa(pageSize))
	if from != nil {
		q.Add("from", from.Format(time.RFC3339))
	}
	if until != nil {
		q.Add("until", until.Format(time.RFC3339))
	}
	req.URL.RawQuery = q.Encode()

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("readSummariesPage do: %w", err)
	}

	return res, nil
}
