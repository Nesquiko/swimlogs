package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (app *SwimLogsApp) SaveTraining(training openapi.NewTraining) Result[openapi.TrainingDetail] {
	if tv := validateNewTraining(training); !tv.IsValid {
		log.Warn().Interface("trainingValidation", tv).Msg("invalid training")
		return errorResult[openapi.TrainingDetail](
			tv.InvalidTraining,
			http.StatusBadRequest,
		)
	}

	recalculateTotalDistances(&training)

	t, err := app.db.SaveTraining(apiNewTrainingIntoDataTraining(training))
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

	return resultWithBody(dataTrainingIntoApiTrainingDetail(t), http.StatusCreated)
}

func recalculateTotalDistances(t *openapi.NewTraining) {
	totalDistance := 0

	for i := range t.Sets {
		setTotalDistance := 0

		if t.Sets[i].SubSets != nil {
			for j := range *t.Sets[i].SubSets {
				subSetTotalDistance := *(*t.Sets[i].SubSets)[j].DistanceMeters * (*t.Sets[i].SubSets)[j].Repeat
				if (*t.Sets[i].SubSets)[j].TotalDistance != subSetTotalDistance {
					log.Debug().
						Int("set_number", i).
						Int("subset_number", j).
						Int("received_total_distance", (*t.Sets[i].SubSets)[j].TotalDistance).
						Int("calculated_total_distance", subSetTotalDistance).
						Msg("total distance in received subset does not match calculated total distance")
					(*t.Sets[i].SubSets)[j].TotalDistance = subSetTotalDistance
				}
				setTotalDistance += subSetTotalDistance
			}
			setTotalDistance *= t.Sets[i].Repeat
		} else {
			setTotalDistance = *t.Sets[i].DistanceMeters * t.Sets[i].Repeat
		}

		if t.Sets[i].TotalDistance != setTotalDistance {
			log.Debug().
				Int("set_number", i).
				Int("received_total_distance", t.Sets[i].TotalDistance).
				Int("calculated_total_distance", setTotalDistance).
				Msg("total distance in received set does not match calculated total distance")
			t.Sets[i].TotalDistance = setTotalDistance
		}
		totalDistance += setTotalDistance
	}

	if t.TotalDistance != totalDistance {
		log.Debug().
			Int("received_total_distance", t.TotalDistance).
			Int("calculated_total_distance", totalDistance).
			Msg("total distance in received training does not match calculated total distance")
		t.TotalDistance = totalDistance
	}
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

	return resultWithBody(dataTrainingIntoApiTraining(t), http.StatusOK)
}

func (app *SwimLogsApp) GetTrainingDetailsForCurrentWeek() Result[openapi.TrainingDetailsCurrentWeekResponse] {
	return app.GetTrainingDetailsInWeek(time.Now())
}

func (app *SwimLogsApp) GetTrainingDetailsInWeek(
	week time.Time,
) Result[openapi.TrainingDetailsCurrentWeekResponse] {
	trainings, err := app.db.GetTrainingDetailsInWeek(week)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get training details")
		return internalServerErrorResult[openapi.TrainingDetailsCurrentWeekResponse](
			"Failed to get training details",
		)
	}

	body := openapi.TrainingDetailsCurrentWeekResponse{
		Details: mapDataTrainingsToApiTrainingDetails(trainings),
	}
	return resultWithBody(body, http.StatusOK)
}
