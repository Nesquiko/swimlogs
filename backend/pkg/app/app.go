package app

import (
	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SwimLogs interface {
	GetAllSessions() (oapiGen.GetAllSessionsResponseObject, error)

	CreateSession(
		request oapiGen.CreateSessionRequestObject,
	) (oapiGen.CreateSessionResponseObject, error)

	DeleteSession(id uuid.UUID) (oapiGen.DeleteSessionResponseObject, error)

	UpdateSession(
		request oapiGen.UpdateSessionRequestObject,
	) (oapiGen.UpdateSessionResponseObject, error)

	CreateTraining(
		request oapiGen.CreateTrainingRequestObject,
	) (oapiGen.CreateTrainingResponseObject, error)

	GetTrainings(
		request oapiGen.GetTrainingsRequestObject,
	) (oapiGen.GetTrainingsResponseObject, error)

	DeleteTraining(
		request oapiGen.DeleteTrainingRequestObject,
	) (oapiGen.DeleteTrainingResponseObject, error)

	GetTrainingById(
		request oapiGen.GetTrainingByIdRequestObject,
	) (oapiGen.GetTrainingByIdResponseObject, error)

	UpdateTraining(
		request oapiGen.UpdateTrainingRequestObject,
	) (oapiGen.UpdateTrainingResponseObject, error)

	GetTrainingsDetails(
		request oapiGen.GetTrainingsDetailsRequestObject,
	) (oapiGen.GetTrainingsDetailsResponseObject, error)

	GetTrainingsDetailsCurrentWeek(
		request oapiGen.GetTrainingsDetailsCurrentWeekRequestObject,
	) (oapiGen.GetTrainingsDetailsCurrentWeekResponseObject, error)

	RecordError(
		request oapiGen.RecordErrorRequestObject,
	) (oapiGen.RecordErrorResponseObject, error)
}

func New(db data.DBConn, logger *zap.SugaredLogger) SwimLogs {
	return &swimLogsApp{db, logger}
}

type swimLogsApp struct {
	db     data.DBConn
	logger *zap.SugaredLogger
}
