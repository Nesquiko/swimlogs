package app

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/google/uuid"
)

func (app *swimLogsApp) GetTrainingsDetails(
	request oapiGen.GetTrainingsDetailsRequestObject,
) (oapiGen.GetTrainingsDetailsResponseObject, error) {
	page, pageSize := request.Params.Page, request.Params.PageSize

	if page < 0 {
		errorDetail := oapiGen.ErrorDetail{
			Title:  "Page can't be less than 0",
			Detail: fmt.Sprintf("Page can't be less than 0, was '%d'", page),
		}
		return oapiGen.GetTrainingsDetails400JSONResponse(errorDetail), nil
	} else if pageSize < 1 {
		errorDetail := oapiGen.ErrorDetail{
			Title:  "Page size can't be less than 1",
			Detail: fmt.Sprintf("Page size can't be less than 1, was '%d'", pageSize),
		}
		return oapiGen.GetTrainingsDetails400JSONResponse(errorDetail), nil
	}

	trainings, err := app.db.GetDetailsOfTrainings(page, pageSize)
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetTrainingsDetails500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}
	totalTrainings, err := app.db.GetTrainingCount()
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetTrainingsDetails500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	details := transormToDetails(trainings)
	pagination := oapiGen.Pagination{
		Total:    totalTrainings,
		Page:     page,
		PageSize: len(details),
	}

	return oapiGen.GetTrainingsDetails200JSONResponse{
		Details:    &details,
		Pagination: &pagination,
	}, nil
}

func (app *swimLogsApp) GetTrainingsDetailsCurrentWeek(
	request oapiGen.GetTrainingsDetailsCurrentWeekRequestObject,
) (oapiGen.GetTrainingsDetailsCurrentWeekResponseObject, error) {
	trainings, err := app.db.GetDetailsOfTrainingsCurrentWeek()
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetTrainingsDetailsCurrentWeek500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	details := transormToDetails(trainings)
	return oapiGen.GetTrainingsDetailsCurrentWeek200JSONResponse{
		Details: &details,
	}, nil
}

func (app *swimLogsApp) CreateTraining(
	request oapiGen.CreateTrainingRequestObject,
) (oapiGen.CreateTrainingResponseObject, error) {
	newTraining := request.Body
	if invalid := validateTraining(*newTraining); len(invalid) != 0 {
		return oapiGen.CreateTraining400JSONResponse{
			InvalidTrainingErrorResponseJSONResponse: invalidTrainingError(invalid),
		}, nil
	}

	t := transformRestTraining(newTraining)
	updateTotalDist(newTraining, t)

	var td oapiGen.TrainingDetail
	err := app.db.InTx(func(tx *sql.Tx) error {
		var (
			id  *uuid.UUID
			err error

			day       = *newTraining.Day
			startTime = *newTraining.StartTime
			durMin    = *newTraining.DurationMin
		)

		if newTraining.SessionId == nil {
			id, err = app.db.SaveTraining(t, tx)
		} else {
			var populatedT *data.Training
			populatedT, err = app.db.SaveTrainingWithSesssionData(t, *newTraining.SessionId, tx)
			id = &populatedT.Id
			day = oapiGen.Day(*populatedT.Day)
			startTime = *populatedT.StartTime
			durMin = *populatedT.DurationMin
		}
		if err != nil {
			return err
		}

		td.Id = *id
		td.Date = newTraining.Date
		td.Day = day
		td.StartTime = startTime
		td.DurationMin = durMin
		td.TotalDist = t.TotalDistance

		return nil
	})
	if err != nil {
		app.logger.Error(err)
		return oapiGen.CreateTraining500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.CreateTraining201JSONResponse(td), nil
}

func (app *swimLogsApp) GetTrainings(
	request oapiGen.GetTrainingsRequestObject,
) (oapiGen.GetTrainingsResponseObject, error) {
	app.logger.Info("GetTrainings endpoint called, but it shouldn't have been")
	return nil, nil
}

func (app *swimLogsApp) DeleteTraining(
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	err := app.db.InTx(func(tx *sql.Tx) error {
		return app.db.DeleteTraining(request.Id, tx)
	})

	if err == data.ErrRowNotFound {
		return oapiGen.DeleteTraining404JSONResponse{
			TrainingNotFoundErrorResponseJSONResponse: trainingNotFound(request.Id),
		}, nil
	}

	if err != nil {
		app.logger.Error(err)
		return oapiGen.DeleteTraining500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.DeleteTraining200Response{}, nil
}

func (app *swimLogsApp) GetTrainingById(
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	t, err := app.db.GetTrainingById(request.Id)
	if err == data.ErrRowNotFound {
		return oapiGen.GetTrainingById404JSONResponse{
			TrainingNotFoundErrorResponseJSONResponse: trainingNotFound(request.Id),
		}, nil
	}

	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetTrainingById500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	training := transformDataTraining(t)
	return oapiGen.GetTrainingById200JSONResponse(training), nil
}

func (app *swimLogsApp) UpdateTraining(
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	newTraining := request.Body
	if invalid := validateTraining(*newTraining); len(invalid) != 0 {
		return oapiGen.UpdateTraining400JSONResponse{
			InvalidTrainingErrorResponseJSONResponse: invalidTrainingError(invalid),
		}, nil
	}

	t := transformRestTraining(newTraining)
	updateTotalDist(newTraining, t)

	err := app.db.InTx(func(tx *sql.Tx) error {
		err := app.db.UpdateTrainingById(request.Id, t, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, data.ErrRowNotFound) {
		app.logger.Warn(err)
		return oapiGen.UpdateTraining404JSONResponse{
			TrainingNotFoundErrorResponseJSONResponse: trainingNotFound(request.Id),
		}, nil
	}
	if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateTraining500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	td := oapiGen.TrainingDetail{
		Id:          request.Id,
		Date:        newTraining.Date,
		Day:         *newTraining.Day,
		StartTime:   *newTraining.StartTime,
		DurationMin: *newTraining.DurationMin,
		TotalDist:   *newTraining.TotalDist,
	}

	return oapiGen.UpdateTraining200JSONResponse(td), nil
}

func updateTotalDist(t *oapiGen.Training, data data.Training) {
	t.TotalDist = &data.TotalDistance

	for i, b := range data.Blocks {
		bTotDist := b.TotalDistance
		t.Blocks[i].TotalDist = &bTotDist
		for j, s := range b.Sets {
			sTotDist := s.TotalDistance
			t.Blocks[i].Sets[j].TotalDist = &sTotDist
		}
	}
}
