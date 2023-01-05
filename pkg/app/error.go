package app

import (
	"fmt"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/google/uuid"
)

func invalidSessionError(
	invalid map[string]string,
) oapiGen.InvalidSessionErrorResponseJSONResponse {
	errorDetail := oapiGen.InvalidSessionErrorResponseJSONResponse{
		Title:                "Invalid request",
		Detail:               "There were invalid session attributes",
		AdditionalProperties: invalid,
	}
	return errorDetail
}

func sessionNotFound(id uuid.UUID) oapiGen.SessionNotFoundErrorResponseJSONResponse {
	return oapiGen.SessionNotFoundErrorResponseJSONResponse{
		Title:  "Session wasn't found",
		Detail: fmt.Sprintf("Session with Id '%s' wasn't found", id.String()),
	}
}

func internalServerError() oapiGen.InternalServerErrorResponseJSONResponse {
	return oapiGen.InternalServerErrorResponseJSONResponse{
		Title:  "Internal server error",
		Detail: "Internal server error",
	}
}
