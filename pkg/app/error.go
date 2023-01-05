package app

import "github.com/Nesquiko/swimlogs/generator/oapiGen"

func invalidSessionError(invalid map[string]string) oapiGen.CreateSession400JSONResponse {
	errorDetail := oapiGen.InvalidSessionErrorResponseJSONResponse{
		Title:                "Invalid request",
		Detail:               "There were invalid session attributes",
		AdditionalProperties: invalid,
	}
	return oapiGen.CreateSession400JSONResponse{
		InvalidSessionErrorResponseJSONResponse: errorDetail,
	}
}

func internalServerError() oapiGen.InternalServerErrorResponseJSONResponse {
	return oapiGen.InternalServerErrorResponseJSONResponse{
		Title:  "Internal server error",
		Detail: "Internal server error",
	}
}
