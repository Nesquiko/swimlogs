package server

import (
	"context"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

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
	return nil, nil
}

func (server *SwimLogsServer) CreateSession(
	ctx context.Context,
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	return server.swimLogs.CreateSession(ctx, request)
}

func (server *SwimLogsServer) DeleteSession(
	ctx context.Context,
	request oapiGen.DeleteSessionRequestObject,
) (oapiGen.DeleteSessionResponseObject, error) {
	return nil, nil
}

func (server *SwimLogsServer) UpdateSession(
	ctx context.Context,
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {
	return nil, nil
}

func (server *SwimLogsServer) CreateTraining(
	ctx context.Context,
	request oapiGen.CreateTrainingRequestObject,
) (oapiGen.CreateTrainingResponseObject, error) {
	return nil, nil
}

func (server *SwimLogsServer) DeleteTraining(
	ctx context.Context,
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	return nil, nil
}

func (server *SwimLogsServer) GetTrainingById(
	ctx context.Context,
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	return nil, nil
}

func (server *SwimLogsServer) UpdateTraining(
	ctx context.Context,
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	return nil, nil
}
