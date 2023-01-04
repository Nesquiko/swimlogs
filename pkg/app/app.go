package app

import (
	"context"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

type SwimLogs interface {
	GetAllSessions(
		ctx context.Context,
	) (oapiGen.GetAllSessionsResponseObject, error)

	CreateSession(
		ctx context.Context,
		request oapiGen.CreateSessionRequestObject,
	) (oapiGen.CreateSessionResponseObject, error)

	DeleteSession(
		ctx context.Context,
		request oapiGen.DeleteSessionRequestObject,
	) (oapiGen.DeleteSessionResponseObject, error)

	UpdateSession(
		ctx context.Context,
		request oapiGen.UpdateSessionRequestObject,
	) (oapiGen.UpdateSessionResponseObject, error)

	CreateTraining(
		ctx context.Context,
		request oapiGen.CreateTrainingRequestObject,
	) (oapiGen.CreateTrainingResponseObject, error)

	DeleteTraining(
		ctx context.Context,
		request oapiGen.DeleteTrainingRequestObject,
	) (oapiGen.DeleteTrainingResponseObject, error)

	GetTrainingById(
		ctx context.Context,
		request oapiGen.GetTrainingByIdRequestObject,
	) (oapiGen.GetTrainingByIdResponseObject, error)

	UpdateTraining(
		ctx context.Context,
		request oapiGen.UpdateTrainingRequestObject,
	) (oapiGen.UpdateTrainingResponseObject, error)
}

func New(db data.DBConn) SwimLogs {
	return &swimLogsApp{db}
}

type swimLogsApp struct {
	db data.DBConn
}
