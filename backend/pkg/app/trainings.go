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

func (app *SwimLogsApp) SaveTraining(training openapi.NewTraining) Result[openapi.TrainingDetail] {
	if errDetails := validateNewTraining(training); errDetails != nil {
		log.Warn().Msg("invalid training")
		return resultFromErrorDetails[openapi.TrainingDetail](
			*errDetails,
			http.StatusBadRequest,
		)
	}

	t, err := app.db.SaveTraining(training)
	if err != nil {
		switch err {
		case data.ErrCheckViolation:
			log.Warn().Err(err).Msg("training violates check constraints")
			return resultWithError[openapi.TrainingDetail](
				"Invalid training",
				"Training does not satisfy check constraints",
				http.StatusBadRequest,
				nil,
			)
		case data.ErrInvalidEnumType:
			log.Warn().Err(err).Msg("training has invalid enum type")
			return resultWithError[openapi.TrainingDetail](
				"Invalid training",
				"Training has invalid enum type",
				http.StatusBadRequest,
				nil,
			)
		case data.ErrForeignKeyViolation:
			log.Warn().Err(err).Msg("training violates foreign key constraints")
			return resultWithError[openapi.TrainingDetail](
				"Invalid training",
				"Training violates foreign key constraints",
				http.StatusBadRequest,
				nil,
			)
		default:
			log.Warn().Err(err).Msg("failed to save training")
			return internalServerErrorResult[openapi.TrainingDetail]("Failed to save training")
		}
	}

	return resultWithBody(t, http.StatusCreated)
}

func (app *SwimLogsApp) GetTrainingById(id uuid.UUID) Result[openapi.Training] {
	t, err := app.db.GetTrainingById(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		log.Warn().Err(err).Str("trainingId", id.String()).Msg("training not found")
		return resultWithError[openapi.Training](
			"Training not found",
			fmt.Sprintf("Training with id %s not found", id.String()),
			http.StatusNotFound,
			nil,
		)
	} else if err != nil {
		log.Warn().Err(err).Msg("failed to get training")
		return internalServerErrorResult[openapi.Training]("Failed to get training")
	}

	return resultWithBody(t, http.StatusOK)
}

func (app *SwimLogsApp) GetTrainingDetailsForCurrentWeek() Result[openapi.TrainingDetailsCurrentWeekResponse] {
	details, err := app.db.GetTrainingDetailsForThisWeek()
	if err != nil {
		log.Warn().Err(err).Msg("failed to get training details")
		return internalServerErrorResult[openapi.TrainingDetailsCurrentWeekResponse](
			"Failed to get training details",
		)
	}

	body := openapi.TrainingDetailsCurrentWeekResponse{Details: details}
	return resultWithBody(body, http.StatusOK)
}
