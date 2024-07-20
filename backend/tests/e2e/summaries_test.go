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

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/app"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestSummariesCurrentWeek_CorrectResponse(t *testing.T) {
	n := 10
	startOfWeek, endOfWeek := app.GetWeekRange(time.Now())

	for i := 0; i < n; i++ {
		start := startOfWeek.AddDate(0, 0, i%7)

		request := &apidef.CreateTrainingRequest{
			DurationMin: 60 + i,
			Sets: []apidef.NewTrainingSet{{
				SetOrder:       0,
				Repeat:         4,
				DistanceMeters: 100,
				Description:    asPtr("Some description"),
				Equipment:      &[]apidef.EquipmentEnum{apidef.Board, apidef.Fins},
				Group:          asPtr(apidef.Long),
				StartSeconds:   asPtr(60),
				StartType:      asPtr(apidef.Interval),
			}},
			Start: start,
		}

		mustCreateNewTraining(t, request)
	}

	assert := assert.New(t)
	summaries := mustReadSummariesCurrentWeek(t)
	assert.NotEmpty(summaries.Summaries)

	for _, summary := range summaries.Summaries {
		assert.WithinRange(summary.Start, startOfWeek, endOfWeek)
	}
}

func TestSummariesPage_CorrectPaging(t *testing.T) {
	n := 30
	for i := 0; i < n; i++ {
		mustCreateNewTraining(t, nil)
	}

	assert := assert.New(t)
	pageSize := 5

	summaries := mustReadSummariesPage(t, 0, pageSize)
	assert.Len(summaries.Summaries, pageSize)
	assert.Equal(0, summaries.Pagination.Page)
	assert.Equal(pageSize, summaries.Pagination.PageSize)
	assert.GreaterOrEqual(summaries.Pagination.Total, n)

	total := summaries.Pagination.Total
	summaries = mustReadSummariesPage(t, 1, total)
	assert.Empty(summaries.Summaries)
	assert.Equal(1, summaries.Pagination.Page)
	assert.Equal(0, summaries.Pagination.PageSize)

	summaries = mustReadSummariesPage(t, 1, total-1)
	assert.Len(summaries.Summaries, 1)
	assert.Equal(1, summaries.Pagination.Page)
	assert.Equal(1, summaries.Pagination.PageSize)
}

func TestSummariesPage_ReadPage(t *testing.T) {
	created := make([]apidef.TrainingSummary, 0)
	for i := 0; i < 10; i++ {
		// In order to make this test pass, the trainings are created
		// with start one year from now, so they will be returned first
		now := time.Now()
		now = now.AddDate(1, 0, i)

		request := &apidef.CreateTrainingRequest{
			DurationMin: 60 + i,
			Sets: []apidef.NewTrainingSet{{
				SetOrder:       0,
				Repeat:         4,
				DistanceMeters: 100,
				Description:    asPtr("Some description"),
				Equipment:      &[]apidef.EquipmentEnum{apidef.Board, apidef.Fins},
				Group:          asPtr(apidef.Long),
				StartSeconds:   asPtr(60),
				StartType:      asPtr(apidef.Interval),
			}},
			Start: now,
		}

		ts := mustCreateNewTraining(t, request)
		created = append(created, ts)
	}

	summaries := mustReadSummariesPage(t, 0, len(created))

	assert := assert.New(t)
	assert.Len(summaries.Summaries, len(created))

outer:
	for _, ts := range summaries.Summaries {
		for _, expected := range created {
			if expected.Id != ts.Id {
				continue
			}

			assert.Truef(compareSummaries(ts, expected), "Summaries aren't equal, expected: %s, actual: %s", expected, ts)
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
	res, err := readSummariesPage(0, invalidPageSize)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError apidef.NotFoundError
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(
		fmt.Sprintf(server.InvalidQueryParamTitleFormat, "pageSize", strconv.Itoa(invalidPageSize)),
		apiError.Title,
	)
	assert.Equal(server.InvalidQueryParamCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(
		fmt.Sprintf(
			server.InvalidQueryParamDetailFormat,
			"pageSize",
			strconv.Itoa(invalidPageSize),
		),
		apiError.Detail,
	)
	assert.Nil(apiError.AdditionalProperties)
}

func TestSummariesPage_InvalidPage(t *testing.T) {
	invalidPage := -1
	res, err := readSummariesPage(invalidPage, 10)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var apiError apidef.NotFoundError
	err = json.NewDecoder(res.Body).Decode(&apiError)
	require.NoError(t, err)

	assert := assert.New(t)
	assert.Equal(
		fmt.Sprintf(server.InvalidQueryParamTitleFormat, "page", strconv.Itoa(invalidPage)),
		apiError.Title,
	)
	assert.Equal(server.InvalidQueryParamCode, apiError.Code)
	assert.Equal(http.StatusBadRequest, apiError.Status)
	assert.Equal(
		fmt.Sprintf(server.InvalidQueryParamDetailFormat, "page", strconv.Itoa(invalidPage)),
		apiError.Detail,
	)
	assert.Nil(apiError.AdditionalProperties)
}

func mustReadSummariesPage(t *testing.T, page, pageSize int) apidef.TrainingSummariesResponse {
	res, err := readSummariesPage(page, pageSize)

	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)

	var summaries apidef.TrainingSummariesResponse
	err = json.NewDecoder(res.Body).Decode(&summaries)
	res.Body.Close()
	require.NoErrorf(t, err, "response: %+v", res)

	return summaries
}

func readSummariesPage(page, pageSize int) (*http.Response, error) {
	url := ServerUrl + "/trainings/summaries"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("readSummariesPage request: %w", err)
	}
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("pageSize", strconv.Itoa(pageSize))
	req.URL.RawQuery = q.Encode()

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("readSummariesPage do: %w", err)
	}

	return res, nil
}

func mustReadSummariesCurrentWeek(t *testing.T) apidef.TrainingSummariesCurrentWeekResponse {
	url := ServerUrl + "/trainings/summaries/current-week"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)

	require.Equalf(t, http.StatusOK, res.StatusCode, "response: %+v", res)
	var summaries apidef.TrainingSummariesCurrentWeekResponse
	err = json.NewDecoder(res.Body).Decode(&summaries)
	res.Body.Close()
	require.NoErrorf(t, err, "response: %+v", res)

	return summaries
}
