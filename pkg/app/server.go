package app

import (
	"context"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

type SwimLogsApp struct {
	db data.DBConn
}

func NewApp(db data.DBConn) *SwimLogsApp {
	return &SwimLogsApp{db}
}

func (app *SwimLogsApp) GetAllSessions(
	ctx context.Context,
	request oapiGen.GetAllSessionsRequestObject,
) (oapiGen.GetAllSessionsResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) CreateSession(
	ctx context.Context,
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) DeleteSession(
	ctx context.Context,
	request oapiGen.DeleteSessionRequestObject,
) (oapiGen.DeleteSessionResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) UpdateSession(
	ctx context.Context,
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) CreateTraining(
	ctx context.Context,
	request oapiGen.CreateTrainingRequestObject,
) (oapiGen.CreateTrainingResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) DeleteTraining(
	ctx context.Context,
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) GetTrainingById(
	ctx context.Context,
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	return nil, nil
}

func (app *SwimLogsApp) UpdateTraining(
	ctx context.Context,
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	return nil, nil
}
