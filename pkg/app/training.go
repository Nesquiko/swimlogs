package app

import (
	"database/sql"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/google/uuid"
)

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

	return nil, nil
}

func updateTotalDist(t *oapiGen.Training, data data.Training) {
	t.TotalDist = &data.TotalDistance

	for i, b := range data.Blocks {
		t.Blocks[i].TotalDist = &b.TotalDistance
		for j, s := range b.Sets {
			t.Blocks[i].Sets[j].TotalDist = &s.TotalDistance
		}
	}
}
