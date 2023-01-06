package app

import (
	"context"
	"database/sql"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
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

	t := transformRestTraining(*newTraining)
	err := app.db.InTx(func(tx *sql.Tx) error {
		var id *uuid.UUID
		var err error

		if newTraining.SessionId == nil {
			id, err = app.db.SaveTraining(t, tx)
		} else {
			id, err = app.db.SaveTrainingWithSesssionData(t, *newTraining.SessionId, tx)
		}

		if err != nil {
			return err
		}
		newTraining.Id = *id

		return nil
	})
	if err != nil {
		app.logger.Error(err)
		return oapiGen.CreateTraining500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.CreateTraining201JSONResponse(*newTraining), nil
}

func (app *swimLogsApp) GetTrainings(
	request oapiGen.GetTrainingsRequestObject,
) (oapiGen.GetTrainingsResponseObject, error) {
	app.logger.Info("GetTrainings endpoint called, but it shouldn't have been")
	return nil, nil
}

func (app *swimLogsApp) DeleteTraining(
	ctx context.Context,
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	return nil, nil
}

func (app *swimLogsApp) GetTrainingById(
	ctx context.Context,
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	return nil, nil
}

func (app *swimLogsApp) UpdateTraining(
	ctx context.Context,
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	return nil, nil
}
