package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (app *SwimLogsApp) GetSessions(
	params openapi.GetSessionsParams,
) Result[openapi.SessionsResponse] {
	sessions, total, err := app.db.GetSessions(params.Page, params.PageSize)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get sessions")
		return internalServerErrorResult[openapi.SessionsResponse]("Failed to get sessions")
	}

	pagination := openapi.Pagination{Page: params.Page, PageSize: len(sessions), Total: total}
	body := openapi.SessionsResponse{
		Sessions:   mapDataSessionsToApiSessions(sessions),
		Pagination: pagination,
	}
	return resultWithBody(body, http.StatusOK)
}

func (app *SwimLogsApp) SaveSession(
	session openapi.CreateSessionJSONBody,
) Result[openapi.Session] {
	if sv := validateNewSession(session); !sv.IsValid {
		log.Warn().Interface("session_validation", sv).Msg("invalid session")
		return errorResult[openapi.Session](sv.InvalidSession, http.StatusBadRequest)
	}

	s, err := app.db.SaveSession(dataSessionFromNewApiSession(session))
	if err != nil {
		switch err {
		case data.ErrCheckViolation:
			log.Warn().Err(err).Msg("session does not pass check constraints")
			return resultWithError[openapi.Session](
				"Invalid session",
				"Session does not has invalid values",
				http.StatusBadRequest,
				nil,
			)
		case data.ErrUniqueViolation:
			log.Warn().Err(err).Msg("session violates unique constraints")
			return resultWithError[openapi.Session](
				"Duplicate session",
				"Session already exists",
				http.StatusConflict,
				nil,
			)
		case data.ErrInvalidEnumType:
			log.Warn().Err(err).Msg("session has invalid enum type")
			return resultWithError[openapi.Session](
				"Invalid session",
				"Session has invalid enum type",
				http.StatusBadRequest,
				nil,
			)
		default:
			log.Warn().Err(err).Msg("failed to save session")
			return internalServerErrorResult[openapi.Session]("Failed to save session")
		}
	}
	return resultWithBody(apiSessionFromDataSession(s), http.StatusCreated)
}

func (app *SwimLogsApp) DeleteSessionById(id uuid.UUID) Result[struct{}] {
	err := app.db.DeleteSessionById(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		log.Warn().Err(err).Msg("session not found")
		return resultWithError[struct{}](
			"Session not found",
			fmt.Sprintf("Session with id %s was not found", id),
			http.StatusNotFound,
			nil,
		)
	} else if err != nil {
		log.Warn().Err(err).Msg("failed to delete session")
		return internalServerErrorResult[struct{}]("Failed to delete session")
	}
	return resultWithoutBody(http.StatusNoContent)
}

func (app *SwimLogsApp) UpdateSession(
	id uuid.UUID,
	updated openapi.UpdateSessionJSONBody,
) Result[openapi.Session] {
	if errDetails := validateSessionUpdate(updated); errDetails != nil {
		log.Warn().Msg("invalid session")
		return resultFromErrorDetails[openapi.Session](
			*errDetails,
			http.StatusBadRequest,
		)
	}

	s, err := app.db.UpdateSession(id, updated)
	if err != nil {
		switch err {
		case data.ErrCheckViolation:
			log.Warn().Err(err).Msg("session does not pass check constraints")
			return resultWithError[openapi.Session](
				"Invalid session",
				"Session does not has invalid values",
				http.StatusBadRequest,
				nil,
			)
		case data.ErrUniqueViolation:
			log.Warn().Err(err).Msg("session violates unique constraints")
			return resultWithError[openapi.Session](
				"Duplicate session",
				"Session already exists",
				http.StatusConflict,
				nil,
			)
		case data.ErrInvalidEnumType:
			log.Warn().Err(err).Msg("session has invalid enum type")
			return resultWithError[openapi.Session](
				"Invalid session",
				"Session has invalid enum type",
				http.StatusBadRequest,
				nil,
			)
		case data.ErrRowsNotFound:
			log.Warn().Err(err).Msg("session not found")
			return resultWithError[openapi.Session](
				"Session not found",
				fmt.Sprintf("Session with id %s was not found", id),
				http.StatusNotFound,
				nil,
			)
		case data.ErrSerializationFailure:
			log.Error().Err(err).Msg("session edited by another user")
			return resultWithError[openapi.Session](
				"Session edited by another user",
				"Session was edited by another user",
				http.StatusConflict,
				nil,
			)
		default:
			log.Warn().Err(err).Msg("failed to save session")
			return internalServerErrorResult[openapi.Session]("Failed to save session")
		}
	}

	return resultWithBody(s, http.StatusOK)
}
