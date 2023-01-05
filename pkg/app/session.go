package app

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func (app *swimLogsApp) GetAllSessions(
	ctx context.Context,
) (oapiGen.GetAllSessionsResponseObject, error) {
	sessions, err := app.db.GetAllSessions()
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetAllSessions500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	ret := make([]oapiGen.Session, len(sessions))
	for i, s := range sessions {
		ret[i] = transformDataSession(s)
	}

	return oapiGen.GetAllSessions200JSONResponse{Sessions: &ret}, nil
}

func (app *swimLogsApp) CreateSession(
	ctx context.Context,
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	newSession := request.Body
	if invalid := validateSession(newSession); len(invalid) != 0 {
		return oapiGen.CreateSession400JSONResponse{
			InvalidSessionErrorResponseJSONResponse: invalidSessionError(invalid),
		}, nil
	}

	session := transformRestSession(*newSession)

	err := app.db.InTx(func(tx *sql.Tx) error {
		uuid, err := app.db.SaveSession(session, tx)
		if err != nil {
			return err
		}
		newSession.Id = *uuid
		return nil
	})
	if err != nil {
		app.logger.Error(err)
		return oapiGen.CreateSession500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.CreateSession201JSONResponse(*newSession), nil
}

func (app *swimLogsApp) DeleteSession(
	ctx context.Context,
	request oapiGen.DeleteSessionRequestObject,
) (oapiGen.DeleteSessionResponseObject, error) {
	return nil, nil
}

func (app *swimLogsApp) UpdateSession(
	ctx context.Context,
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {
	return nil, nil
}

func transformRestSession(session oapiGen.Session) data.Session {
	return data.Session{
		Day:         strings.ToLower(string(session.Day)),
		StartTime:   session.StartTime,
		DurationMin: session.DurationMin,
	}
}

func transformDataSession(session data.Session) oapiGen.Session {
	return oapiGen.Session{
		Id:          session.Id,
		Day:         oapiGen.Day(session.Day),
		StartTime:   session.StartTime,
		DurationMin: session.DurationMin,
	}
}
