package app

import (
	"database/sql"
	"strings"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/google/uuid"
)

func (app *swimLogsApp) GetAllSessions() (oapiGen.GetAllSessionsResponseObject, error) {
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
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	newSession := request.Body
	if invalid := validateSession(*newSession); len(invalid) != 0 {
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

func (app *swimLogsApp) DeleteSession(id uuid.UUID) (oapiGen.DeleteSessionResponseObject, error) {
	err := app.db.InTx(func(tx *sql.Tx) error {
		return app.db.DeleteSession(id, tx)
	})

	if err == data.ErrRowNotFound {
		return oapiGen.DeleteSession404JSONResponse{
			SessionNotFoundErrorResponseJSONResponse: sessionNotFound(id),
		}, nil
	}

	if err != nil {
		app.logger.Error(err)
		return oapiGen.DeleteSession500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.DeleteSession200Response{}, nil
}

func (app *swimLogsApp) UpdateSession(
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {
	if invalid := validateSession(*request.Body); len(invalid) != 0 {
		return oapiGen.UpdateSession400JSONResponse{
			InvalidSessionErrorResponseJSONResponse: invalidSessionError(invalid),
		}, nil
	}
	updated := transformRestSession(*request.Body)

	var session oapiGen.Session
	err := app.db.InTx(func(tx *sql.Tx) error {
		sess, err := app.db.UpdateSession(request.Id, updated, tx)
		if err != nil {
			return err
		}
		session = transformDataSession(sess)
		return nil
	})
	if err == data.ErrRowNotFound {
		return oapiGen.UpdateSession404JSONResponse{
			SessionNotFoundErrorResponseJSONResponse: sessionNotFound(request.Id),
		}, nil
	} else if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateSession500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.UpdateSession200JSONResponse(session), nil
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
