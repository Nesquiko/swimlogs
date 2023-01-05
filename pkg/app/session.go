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
	return nil, nil
}

func (app *swimLogsApp) CreateSession(
	ctx context.Context,
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	newSession := request.Body
	if invalid := validateSession(newSession); len(invalid) != 0 {
		return invalidSessionError(invalid), nil
	}

	session := transformSession(newSession)

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

	return oapiGen.CreateSession201JSONResponse{
		CreateSessionResponseJSONResponse: oapiGen.CreateSessionResponseJSONResponse(*newSession),
	}, nil
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

func transformSession(session *oapiGen.Session) data.Session {
	s := data.Session{
		Day:         strings.ToLower(string(session.Day)),
		StartTime:   session.StartTime,
		DurationMin: session.DurationMin,
	}

	return s
}
