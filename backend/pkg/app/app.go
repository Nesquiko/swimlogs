package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

var ErrNotFound = errors.New("resource not found")

func New(pool *data.PostgresDbPool) SwimLogsApp {
	return SwimLogsApp{pool}
}

type SwimLogsApp struct {
	pool *data.PostgresDbPool
}

func (app SwimLogsApp) CreateTraining(
	newTraining apidef.NewTraining,
) (apidef.TrainingSummary, error) {
	t := newTrainingToDataTraining(newTraining)
	t.Start = t.Start.Truncate(time.Minute)

	t, err := app.pool.PersistTraining(t)
	if err != nil {
		return apidef.TrainingSummary{}, fmt.Errorf("CreateTraining: %w", err)
	}

	return trainingToSummary(t), nil
}

func (app SwimLogsApp) DeleteTraining(id uuid.UUID) error {
	err := app.pool.DeleteTraining(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return fmt.Errorf("DeleteTraining not found: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}
	return nil
}

func (app SwimLogsApp) TrainingSummariesPage(
	page, pageSize int,
) ([]apidef.TrainingSummary, int, error) {
	summariesPage, total, err := app.pool.TrainingSummaries(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummariesPage: %w", err)
	}
	summaries := make([]apidef.TrainingSummary, len(summariesPage))
	for i, d := range summariesPage {
		summaries[i] = trainingToSummary(d)
	}

	return summaries, total, nil
}

func (app SwimLogsApp) TrainingDetailsCurrentWeek() (apidef.TrainingSummariesCurrentWeekResponse, error) {
	// now := time.Now()
	// startOfWeek := now.AddDate(0, 0, -(int(now.Weekday())+6)%7)
	// endOfWeek := now.AddDate(0, 0, (7-int(now.Weekday()))%7)
	//
	// detailsInRange, err := app.pool.TrainingDetailsInRange(startOfWeek, endOfWeek)
	// if err != nil {
	// 	return nil, fmt.Errorf("TrainingDetailsCurrentWeek: %w", err)
	// }
	//
	// details := make([]apidef.TrainingDetail, len(detailsInRange))
	// for i, d := range detailsInRange {
	// 	details[i] = trainingToSummary(d)
	// }
	// return details, nil
	panic("not implemented")
}

func (app SwimLogsApp) Training(id uuid.UUID) (apidef.Training, error) {
	t, err := app.pool.Training(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("Training not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("Training: %w", err)
	}

	return dataTrainingToApiTraining(t), nil
}

func (app SwimLogsApp) EditTraining(
	id uuid.UUID,
	t apidef.Training,
) (apidef.TrainingSummary, error) {
	// recalcDistanceOnTraining(&t)
	// training := trainingToDataTraining(t)
	//
	// edited, err := app.pool.EditTraining(id, training)
	// if errors.Is(err, data.ErrRowsNotFound) {
	// 	return apidef.TrainingDetail{}, fmt.Errorf("EditTraining: %w", ErrNotFound)
	// } else if err != nil {
	// 	return apidef.TrainingDetail{}, fmt.Errorf("EditTraining: %w", err)
	// }
	//
	// return trainingToSummary(edited), nil
	panic("not implemented")
}
