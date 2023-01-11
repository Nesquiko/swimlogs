package app

import (
	"fmt"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
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

func invalidTrainingError(
	invalid map[string]string,
) oapiGen.InvalidTrainingErrorResponseJSONResponse {
	return oapiGen.InvalidTrainingErrorResponseJSONResponse{
		Title:                "Invalid request",
		Detail:               "There were invalid training attributes",
		AdditionalProperties: invalid,
	}
}

func trainingNotFound(id uuid.UUID) oapiGen.TrainingNotFoundErrorResponseJSONResponse {
	return oapiGen.TrainingNotFoundErrorResponseJSONResponse{
		Title:  "Training wasn't found",
		Detail: fmt.Sprintf("Training with Id '%s' wasn't found", id.String()),
	}
}
