package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

var ErrInternalServerError = errors.New("internal server error")

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
