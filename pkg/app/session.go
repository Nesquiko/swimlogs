package app

import (
	"database/sql"
	"fmt"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/validation"
	"github.com/google/uuid"
)

func (app *swimLogsApp) GetAllSessions() (oapiGen.GetAllSessionsResponseObject, error) {
	sessions, err := app.db.GetAllSessions()
	if err != nil {
		app.logger.Error(err)
		return oapiGen.GetAllSessions500Response{}, nil
	}

	ret := make([]oapiGen.Session, len(sessions))
	for i, s := range sessions {
		ret[i] = transformDataSession(s)
	}

	return oapiGen.GetAllSessions200JSONResponse{Sessions: ret}, nil
}

func (app *swimLogsApp) CreateSession(
	request oapiGen.CreateSessionRequestObject,
) (oapiGen.CreateSessionResponseObject, error) {
	newSession := request.Body

	invalid := validation.ValidateSession(*newSession)
	if invalid != nil {
		app.logger.Warnf("invalid session: %v", invalid)
		return oapiGen.CreateSession400JSONResponse{
			InvalidSessionErrorResponseJSONResponse: oapiGen.InvalidSessionErrorResponseJSONResponse(
				*invalid,
			),
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
		return oapiGen.CreateSession500Response{}, nil
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
		return oapiGen.DeleteSession500Response{}, nil
	}

	return oapiGen.DeleteSession200Response{}, nil
}

func (app *swimLogsApp) UpdateSession(
	request oapiGen.UpdateSessionRequestObject,
) (oapiGen.UpdateSessionResponseObject, error) {

	exists, err := app.db.SessionExists(request.Id)
	if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateSession500Response{}, nil
	}

	if !exists {
		return oapiGen.UpdateSession404JSONResponse{
			SessionNotFoundErrorResponseJSONResponse: sessionNotFound(request.Id),
		}, nil
	}

	invalid := validation.ValidateSession(*request.Body)
	if invalid != nil {
		return oapiGen.UpdateSession400JSONResponse{
			oapiGen.InvalidSessionErrorResponseJSONResponse(*invalid),
		}, nil
	}

	updated := transformRestSession(*request.Body)
	var session oapiGen.Session
	err = app.db.InTx(func(tx *sql.Tx) error {
		sess, err := app.db.UpdateSession(request.Id, updated, tx)
		if err != nil {
			return err
		}
		session = transformDataSession(sess)
		return nil
	})
	if err == data.ErrRowNotFound {
		return oapiGen.UpdateSession409Response{}, nil
	} else if err != nil {
		app.logger.Error(err)
		return oapiGen.UpdateSession500Response{}, nil
	}

	return oapiGen.UpdateSession200JSONResponse(session), nil
}

func sessionNotFound(id uuid.UUID) oapiGen.SessionNotFoundErrorResponseJSONResponse {
	return oapiGen.SessionNotFoundErrorResponseJSONResponse{
		Title:  "Session wasn't found",
		Detail: fmt.Sprintf("Session with Id '%s' wasn't found", id.String()),
	}
}
