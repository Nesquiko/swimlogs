package app

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/validation"
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
		return oapiGen.GetTrainingsDetails500Response{}, nil
	}
	totalTrainings, err := app.db.GetTrainingCount()
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetTrainingsDetails500Response{}, nil
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
		return oapiGen.GetTrainingsDetailsCurrentWeek500Response{}, nil
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

	invalid := app.validateTraining(*newTraining)
	if invalid != nil {
		return oapiGen.CreateTraining400JSONResponse{
			InvalidTrainingErrorResponseJSONResponse: oapiGen.InvalidTrainingErrorResponseJSONResponse(
				*invalid,
			),
		}, nil
	}

	t := transformRestTraining(*newTraining)
	calculateTotalDistance(&t)
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
		return oapiGen.CreateTraining500Response{}, nil
	}

	return oapiGen.CreateTraining201JSONResponse(td), nil
}

func (app *swimLogsApp) validateTraining(t oapiGen.Training) *oapiGen.InvalidTraining {
	if t.SessionId == nil {
		return validation.ValidateTraining(t)
	}
	s, err := app.db.GetSessionById(*t.SessionId)
	if err != nil {
		app.logger.Errorf("received unknown session: %v", err)
	}
	return validation.ValidateTrainingWithSession(t, transformDataSession(s))
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
		return oapiGen.DeleteTraining500Response{}, nil
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
		return oapiGen.GetTrainingById500Response{}, nil
	}

	training := transformDataTraining(t)
	return oapiGen.GetTrainingById200JSONResponse(training), nil
}

func (app *swimLogsApp) UpdateTraining(
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {

	exists, err := app.db.TrainingExists(request.Id)
	if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateTraining500Response{}, nil
	}

	if !exists {
		return oapiGen.UpdateTraining404JSONResponse{
			TrainingNotFoundErrorResponseJSONResponse: trainingNotFound(request.Id),
		}, nil
	}

	newTraining := request.Body
	invalid := app.validateTraining(*newTraining)
	if invalid != nil {
		return oapiGen.UpdateTraining400JSONResponse{
			InvalidTrainingErrorResponseJSONResponse: oapiGen.InvalidTrainingErrorResponseJSONResponse(
				*invalid,
			),
		}, nil
	}

	t := transformRestTraining(*newTraining)
	updateTotalDist(newTraining, t)

	err = app.db.InTx(func(tx *sql.Tx) error {
		err := app.db.UpdateTrainingById(request.Id, t, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, data.ErrRowNotFound) {
		app.logger.Warn(err)
		return oapiGen.UpdateSession409Response{}, nil
	}
	if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateTraining500Response{}, nil
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

// calculateTotalDistance calculates distances from sets and up
func calculateTotalDistance(t *data.Training) {
	tTotDist := 0
	for i := range t.Blocks {
		bTotDist := 0

		for j := range t.Blocks[i].Sets {
			sTotDist := t.Blocks[i].Sets[j].Repeat * t.Blocks[i].Sets[j].Distance
			bTotDist += sTotDist

			t.Blocks[i].Sets[j].TotalDistance = sTotDist
		}

		tTotDist += t.Blocks[i].Repeat * bTotDist
		t.Blocks[i].TotalDistance = t.Blocks[i].Repeat * bTotDist
	}
	t.TotalDistance = tTotDist
}

func trainingNotFound(id uuid.UUID) oapiGen.TrainingNotFoundErrorResponseJSONResponse {
	return oapiGen.TrainingNotFoundErrorResponseJSONResponse{
		Title:  "Training wasn't found",
		Detail: fmt.Sprintf("Training with Id '%s' wasn't found", id.String()),
	}
}
