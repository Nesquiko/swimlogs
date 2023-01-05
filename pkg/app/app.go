package app

import (
	"context"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SwimLogs interface {
	GetAllSessions(
		ctx context.Context,
	) (oapiGen.GetAllSessionsResponseObject, error)

	CreateSession(
		ctx context.Context,
		request oapiGen.CreateSessionRequestObject,
	) (oapiGen.CreateSessionResponseObject, error)

	DeleteSession(ctx context.Context, id uuid.UUID) (oapiGen.DeleteSessionResponseObject, error)

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

func New(db data.DBConn, logger *zap.SugaredLogger) SwimLogs {
	return &swimLogsApp{db, logger}
}

type swimLogsApp struct {
	db     data.DBConn
	logger *zap.SugaredLogger
}
