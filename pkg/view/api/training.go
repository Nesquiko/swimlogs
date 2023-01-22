package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

var ErrInternalServerError = errors.New("internal server error")
var ErrNotFound = errors.New("resource not found")

func GetTrainingsInCurrentWeek() ([]oapiGen.TrainingDetail, error) {
	res, err := http.Get(BaseUrl + "/trainings/details/current-week")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var dest oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse
	switch res.StatusCode {
	case http.StatusOK:
		dest = oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse{}
	case http.StatusInternalServerError:
		return nil, ErrInternalServerError
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return nil, err
	}

	return *dest.Details, nil
}

func FetchTraining(id uuid.UUID) (oapiGen.Training, error) {
	url := BaseUrl + "/trainings/" + id.String()
	res, err := http.Get(url)
	if err != nil {
		return oapiGen.Training{}, err
	}
	defer res.Body.Close()

	var dest oapiGen.GetTrainingById200JSONResponse
	switch res.StatusCode {
	case http.StatusOK:
		dest = oapiGen.GetTrainingById200JSONResponse{}
	case http.StatusNotFound:
		return oapiGen.Training{}, ErrNotFound
	case http.StatusInternalServerError:
		return oapiGen.Training{}, ErrInternalServerError
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return oapiGen.Training{}, err
	}

	return oapiGen.Training(dest), nil
}
