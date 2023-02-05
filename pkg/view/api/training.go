package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

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
