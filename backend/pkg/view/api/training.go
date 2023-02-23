package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

func CreateTraining(t oapiGen.Training) (oapiGen.TrainingDetail, error) {
	trainingJson, err := json.Marshal(t)
	if err != nil {
		return oapiGen.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", err)
	}

	res, err := http.Post(BaseUrl+"/trainings", AppJsonType, bytes.NewBuffer(trainingJson))
	if err != nil {
		return oapiGen.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusBadRequest:
		dest := oapiGen.CreateTraining400JSONResponse{}
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&dest)
		if err != nil {
			return oapiGen.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", err)
		}
		return oapiGen.TrainingDetail{}, &BadRequestTraining{
			oapiGen.InvalidTraining(dest.InvalidTrainingErrorResponseJSONResponse),
		}
	case http.StatusInternalServerError:
		return oapiGen.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", ErrInternalServerError)
	}

	// training was successfully created
	var dest oapiGen.CreateTraining201JSONResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return oapiGen.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", err)
	}

	return oapiGen.TrainingDetail(dest), nil
}

func GetTrainingDetails(page, pageSize int) ([]oapiGen.TrainingDetail, oapiGen.Pagination, error) {
	path := fmt.Sprintf("/trainings/details?page=%d&pageSize=%d", page, pageSize)
	res, err := http.Get(BaseUrl + path)
	if err != nil {
		return nil, oapiGen.Pagination{}, fmt.Errorf("GetTrainingDetails: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, oapiGen.Pagination{}, fmt.Errorf(
			"GetTrainingDetails: %w",
			ErrInternalServerError,
		)
	}

	dest := oapiGen.GetTrainingsDetails200JSONResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return nil, oapiGen.Pagination{}, fmt.Errorf("GetTrainingDetails: %w", err)
	}

	return dest.Details, dest.Pagination, nil
}

func GetTrainingsInCurrentWeek() ([]oapiGen.TrainingDetail, error) {
	res, err := http.Get(BaseUrl + "/trainings/details/current-week")
	if err != nil {
		return nil, fmt.Errorf("GetTrainingsInCurrentWeek: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	var dest oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse
	switch res.StatusCode {
	case http.StatusOK:
		dest = oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse{}
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("GetTrainingsInCurrentWeek: %w", ErrInternalServerError)
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return nil, fmt.Errorf("GetTrainingsInCurrentWeek: %w", err)
	}

	return *dest.Details, nil
}

func FetchTraining(id uuid.UUID) (oapiGen.Training, error) {
	url := BaseUrl + "/trainings/" + id.String()
	res, err := http.Get(url)
	if err != nil {
		return oapiGen.Training{}, fmt.Errorf("FetchTraining: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	var dest oapiGen.GetTrainingById200JSONResponse
	switch res.StatusCode {
	case http.StatusOK:
		dest = oapiGen.GetTrainingById200JSONResponse{}
	case http.StatusNotFound:
		return oapiGen.Training{}, fmt.Errorf("FetchTraining: %w", ErrNotFound)
	case http.StatusInternalServerError:
		return oapiGen.Training{}, fmt.Errorf("FetchTraining: %w", ErrInternalServerError)
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return oapiGen.Training{}, fmt.Errorf("FetchTraining: %w", err)
	}

	return oapiGen.Training(dest), nil
}

type BadRequestTraining struct {
	oapiGen.InvalidTraining
}

func (r *BadRequestTraining) Error() string {
	return fmt.Sprintf("invalid training, %v", r.InvalidTraining)
}

func (r *BadRequestTraining) Is(tgt error) bool {
	_, ok := tgt.(*BadRequestTraining)
	return ok
}
