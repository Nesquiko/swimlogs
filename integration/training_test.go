package integration

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
)

func TestGetTrainingsDetailsCurrentWeek(t *testing.T) {

	year, month, day := time.Now().Year(), int(time.Now().Month()), time.Now().Day()
	dates := []struct {
		date string
		day  oapiGen.Day
	}{
		{fmt.Sprintf("%d-%02d-%02d", year, month, day-7), getNameOfDay(year, month, day-7)},
		{fmt.Sprintf("%d-%02d-%02d", year, month, day), getNameOfDay(year, month, day)},
		{fmt.Sprintf("%d-%02d-%02d", year, month, day+7), getNameOfDay(year, month, day+7)},
	}

	for _, v := range dates {
		training := createValidTraining(nil)
		d, err := time.Parse("2006-01-02", v.date)
		if err != nil {
			t.Fatal(err)
		}
		training.Date = types.Date{Time: d}
		training.Day = &v.day
		saveTraining(training, t)
	}

	req := oapiGen.GetTrainingsDetailsCurrentWeekRequestObject{}
	res, err := SwimLogsApp.GetTrainingsDetailsCurrentWeek(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	trainings, ok := res.(oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse)
	if !ok {
		t.Fatalf("expected successfull response, but response was %+v", trainings)
	}
	details := *trainings.Details

	if len(details) != 1 {
		t.Fatalf("expected %d trainings, response was: %+v", 1, details)
	}

	cleanDB(DB)
}

func TestUpdateTraining(t *testing.T) {
	id := saveTraining(createValidTraining(nil), t)
	trainingResponse, _ := SwimLogsApp.GetTrainingById(
		oapiGen.GetTrainingByIdRequestObject{Id: *id},
	)
	training, _ := trainingResponse.(oapiGen.GetTrainingById200JSONResponse)

	startTime := "19:00"
	training.StartTime = &startTime
	updated := oapiGen.Training(training)

	req := oapiGen.UpdateTrainingRequestObject{Id: *id, Body: &updated}

	res, err := SwimLogsApp.UpdateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	detail, ok := res.(oapiGen.UpdateTraining200JSONResponse)
	if !ok {
		t.Fatalf("expected successfull response, but response was %+v", detail)
	}

	if detail.StartTime != startTime {
		t.Errorf("start time, expected %s, but was %s", startTime, detail.Date)
	}
	cleanDB(DB)
}

func TestUpdateTrainingInvalidTraining(t *testing.T) {
	id := saveTraining(createValidTraining(nil), t)

	updated := createValidTraining(nil)
	startTime := "99:00"
	updated.StartTime = &startTime

	req := oapiGen.UpdateTrainingRequestObject{Id: *id, Body: &updated}

	res, err := SwimLogsApp.UpdateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.UpdateTraining400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Invalid request"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := "There were invalid training attributes"
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	expectedAddProps := map[string]string{
		"startTime": fmt.Sprintf("Start time must be from 00:00 to 23:59, but was '%s'", startTime),
	}
	if !reflect.DeepEqual(expectedAddProps, errDetail.AdditionalProperties) {
		t.Fatalf(
			"error details additional properties, expected %v, but was %v",
			expectedAddProps,
			errDetail.AdditionalProperties,
		)
	}
	cleanDB(DB)
}

func TestUpdateTrainingNotFound(t *testing.T) {
	id := saveTraining(createValidTraining(nil), t)

	updated := createValidTraining(nil)
	startTime := "09:00"
	updated.StartTime = &startTime

	req := oapiGen.UpdateTrainingRequestObject{Id: *id, Body: &updated}

	res, err := SwimLogsApp.UpdateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.UpdateTraining404JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Training wasn't found"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Training with Id '%s' wasn't found", id.String())
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}
	cleanDB(DB)
}

func TestDeleteTraining(t *testing.T) {
	id := saveTraining(createValidTraining(nil), t)

	req := oapiGen.DeleteTrainingRequestObject{Id: *id}

	res, err := SwimLogsApp.DeleteTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	response, ok := res.(oapiGen.DeleteTraining200Response)
	if !ok {
		t.Fatalf("expected successfull response, but response was %+v", response)
	}
}

func TestDeleteTrainingNotFound(t *testing.T) {
	id := uuid.New()
	req := oapiGen.DeleteTrainingRequestObject{Id: id}

	res, err := SwimLogsApp.DeleteTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.DeleteTraining404JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Training wasn't found"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Training with Id '%s' wasn't found", id.String())
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}

}

func TestGetTrainingById(t *testing.T) {
	id := saveTraining(createValidTraining(nil), t)

	req := oapiGen.GetTrainingByIdRequestObject{Id: *id}

	res, err := SwimLogsApp.GetTrainingById(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	training, ok := res.(oapiGen.GetTrainingById200JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", training)
	}

	if training.Id != *id {
		t.Errorf("id, expected %s, but was %s", id.String(), training.Id.String())
	}
	cleanDB(DB)
}

func TestGetTrainingByIdNotFound(t *testing.T) {
	id := uuid.New()
	req := oapiGen.GetTrainingByIdRequestObject{Id: id}

	res, err := SwimLogsApp.GetTrainingById(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.GetTrainingById404JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Training wasn't found"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Training with Id '%s' wasn't found", id.String())
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}
}

func TestGetTrainingsDetailsPagination(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		training := createValidTraining(nil)
		saveTraining(training, t)
	}

	page, pageSize := 0, 3
	req := oapiGen.GetTrainingsDetailsRequestObject{
		Params: oapiGen.GetTrainingsDetailsParams{Page: page, PageSize: pageSize},
	}

	res, err := SwimLogsApp.GetTrainingsDetails(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	details, ok := res.(oapiGen.GetTrainingsDetails200JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", details)
	}
	pagination := details.Pagination

	if pagination.Total != n {
		t.Errorf("total, expected %d, but was %d", n, pagination.Total)
	}
	if pagination.Page != page {
		t.Errorf("page, expected %d, but was %d", page, pagination.Page)
	}
	if pagination.PageSize != pageSize {
		t.Errorf("page size, expected %d, but was %d", pageSize, pagination.PageSize)
	}

	req.Params.Page++
	res, err = SwimLogsApp.GetTrainingsDetails(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}
	details, ok = res.(oapiGen.GetTrainingsDetails200JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", details)
	}
	pagination = details.Pagination

	if pagination.Page != page+1 {
		t.Errorf("next page, expected %d, but was %d", page+1, pagination.Page)
	}
	if pagination.PageSize != n-pageSize {
		t.Errorf("next page size, expected %d, but was %d", n-pageSize, pagination.PageSize)
	}
	cleanDB(DB)
}

func TestGetTrainingsDetailsInvalidPageSize(t *testing.T) {
	pageSize := 0
	req := oapiGen.GetTrainingsDetailsRequestObject{
		Params: oapiGen.GetTrainingsDetailsParams{Page: 1, PageSize: pageSize},
	}

	res, err := SwimLogsApp.GetTrainingsDetails(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.GetTrainingsDetails400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Page size can't be less than 1"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Page size can't be less than 1, was '%d'", pageSize)
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}
}

func TestGetTrainingsDetailsInvalidPage(t *testing.T) {
	page := -1
	req := oapiGen.GetTrainingsDetailsRequestObject{
		Params: oapiGen.GetTrainingsDetailsParams{Page: page, PageSize: 1},
	}

	res, err := SwimLogsApp.GetTrainingsDetails(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.GetTrainingsDetails400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Page can't be less than 0"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := fmt.Sprintf("Page can't be less than 0, was '%d'", page)
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	if len(errDetail.AdditionalProperties) != 0 {
		t.Errorf("expected no additional properties, received %v", errDetail.AdditionalProperties)
	}
}

func TestCreateTraining(t *testing.T) {
	training := createValidTraining(nil)

	req := oapiGen.CreateTrainingRequestObject{Body: &training}
	res, err := SwimLogsApp.CreateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	trainingDetail, ok := res.(oapiGen.CreateTraining201JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", trainingDetail)
	}

	if trainingDetail.Day != *training.Day {
		t.Errorf("day, expected %s, but was %s", *training.Day, trainingDetail.Day)
	}
	if trainingDetail.DurationMin != *training.DurationMin {
		t.Errorf(
			"duration, expected %d, but was %d",
			*training.DurationMin,
			trainingDetail.DurationMin,
		)
	}
	if trainingDetail.StartTime != *training.StartTime {
		t.Errorf(
			"start time, expected %s, but was %s",
			*training.StartTime,
			trainingDetail.StartTime,
		)
	}
	cleanDB(DB)
}

func TestCreateTrainingFromSession(t *testing.T) {
	sId := saveSession(validSession, t)
	training := createValidTraining(sId)

	req := oapiGen.CreateTrainingRequestObject{Body: &training}
	res, err := SwimLogsApp.CreateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	trainingDetail, ok := res.(oapiGen.CreateTraining201JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", trainingDetail)
	}

	if trainingDetail.Day != validSession.Day {
		t.Errorf("day, expected %s, but was %s", validSession.Day, trainingDetail.Day)
	}
	if trainingDetail.DurationMin != validSession.DurationMin {
		t.Errorf(
			"duration, expected %d, but was %d",
			validSession.DurationMin,
			trainingDetail.DurationMin,
		)
	}
	if trainingDetail.StartTime != validSession.StartTime {
		t.Errorf(
			"start time, expected %s, but was %s",
			validSession.StartTime,
			trainingDetail.StartTime,
		)
	}
	cleanDB(DB)
}

func TestCreateTrainingInvalid(t *testing.T) {
	training := createTrainingWithNoBlocks()

	req := oapiGen.CreateTrainingRequestObject{Body: &training}

	res, err := SwimLogsApp.CreateTraining(req)
	if err != nil {
		t.Fatalf("expected no error, but was %v", err)
	}

	errDetail, ok := res.(oapiGen.CreateTraining400JSONResponse)
	if !ok {
		t.Fatalf("expected error details, but response was %+v", errDetail)
	}

	expectedTitle := "Invalid request"
	if errDetail.Title != expectedTitle {
		t.Errorf("error deatils title, expected %q, but was %q", expectedTitle, errDetail.Title)
	}
	expectedDetail := "There were invalid training attributes"
	if errDetail.Detail != expectedDetail {
		t.Errorf("error deatils detail, expected %q, but was %q", expectedDetail, errDetail.Detail)
	}
	expectedAddProps := map[string]string{"blocks": "No blocks in training"}
	if !reflect.DeepEqual(expectedAddProps, errDetail.AdditionalProperties) {
		t.Fatalf(
			"error details additional properties, expected %v, but was %v",
			expectedAddProps,
			errDetail.AdditionalProperties,
		)
	}
}
