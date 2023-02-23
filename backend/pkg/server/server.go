package server

import (
	"context"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

// SwimLogsServer implements interface oapiGen.StrictServerInterface,
// which is generated from OpenApi spec located in documentation.
type SwimLogsServer struct {
	swimLogs app.SwimLogs
}

func New(swimLogs app.SwimLogs) *SwimLogsServer {
	return &SwimLogsServer{swimLogs: swimLogs}
}

func (server *SwimLogsServer) GetAllSessions(
	ctx context.Context,
	request oapiGen.GetAllSessionsRequestObject,
) (oapiGen.GetAllSessionsResponseObject, error) {
	return server.swimLogs.GetAllSessions()
}

func (server *SwimLogsServer) CreateSession(
	ctx context.Context,
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	return server.swimLogs.CreateSession(request)
}

func (server *SwimLogsServer) DeleteSession(
	ctx context.Context,
	request oapiGen.DeleteSessionRequestObject,
) (oapiGen.DeleteSessionResponseObject, error) {
	return server.swimLogs.DeleteSession(request.Id)
}

func (server *SwimLogsServer) UpdateSession(
	ctx context.Context,
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {
	return server.swimLogs.UpdateSession(request)
}

func (server *SwimLogsServer) GetTrainings(
	ctx context.Context,
	request oapiGen.GetTrainingsRequestObject,
) (oapiGen.GetTrainingsResponseObject, error) {
	return server.swimLogs.GetTrainings(request)
}

func (server *SwimLogsServer) CreateTraining(
	ctx context.Context,
	request oapiGen.CreateTrainingRequestObject,
) (oapiGen.CreateTrainingResponseObject, error) {
	return server.swimLogs.CreateTraining(request)
}

func (server *SwimLogsServer) DeleteTraining(
	ctx context.Context,
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	return server.swimLogs.DeleteTraining(request)
}

func (server *SwimLogsServer) GetTrainingById(
	ctx context.Context,
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	return server.swimLogs.GetTrainingById(request)
}

func (server *SwimLogsServer) UpdateTraining(
	ctx context.Context,
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	return server.swimLogs.UpdateTraining(request)
}

func (server *SwimLogsServer) GetTrainingsDetails(
	ctx context.Context,
	request oapiGen.GetTrainingsDetailsRequestObject,
) (oapiGen.GetTrainingsDetailsResponseObject, error) {
	return server.swimLogs.GetTrainingsDetails(request)
}

func (server *SwimLogsServer) GetTrainingsDetailsCurrentWeek(
	ctx context.Context,
	request oapiGen.GetTrainingsDetailsCurrentWeekRequestObject,
) (oapiGen.GetTrainingsDetailsCurrentWeekResponseObject, error) {
	return server.swimLogs.GetTrainingsDetailsCurrentWeek(request)
}

func (server *SwimLogsServer) RecordError(
	ctx context.Context,
	request oapiGen.RecordErrorRequestObject,
) (oapiGen.RecordErrorResponseObject, error) {
	return server.swimLogs.RecordError(request)
}
