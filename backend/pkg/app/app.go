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
) (apidef.TrainingDetail, error) {
	recalcDistanceOnNewTraining(&newTraining)
	t := newTrainingToDataTraining(newTraining)

	t.Start = t.Start.Truncate(time.Minute)
	t, err := app.pool.PersistTraining(t)
	if err != nil {
		return apidef.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", err)
	}

	return trainingToDetail(t), nil
}

func (app SwimLogsApp) DeleteTraining(id uuid.UUID) error {
	err := app.pool.DeleteTraining(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return fmt.Errorf("DeleteTraining: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}
	return nil
}

func (app SwimLogsApp) TrainingDetailsPage(
	page, pageSize int,
) ([]apidef.TrainingDetail, int, error) {
	detailsPage, total, err := app.pool.TrainingDetails(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingDetailsPage: %w", err)
	}

	details := make([]apidef.TrainingDetail, len(detailsPage))
	for i, d := range detailsPage {
		details[i] = trainingToDetail(d)
	}

	return details, total, nil
}

func (app SwimLogsApp) TrainingDetailsCurrentWeek() ([]apidef.TrainingDetail, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -(int(now.Weekday())+6)%7)
	endOfWeek := now.AddDate(0, 0, (7-int(now.Weekday()))%7)

	detailsInRange, err := app.pool.TrainingDetailsInRange(startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("TrainingDetailsCurrentWeek: %w", err)
	}

	details := make([]apidef.TrainingDetail, len(detailsInRange))
	for i, d := range detailsInRange {
		details[i] = trainingToDetail(d)
	}
	return details, nil
}

func (app SwimLogsApp) Training(id uuid.UUID) (apidef.Training, error) {
	t, err := app.pool.Training(id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("Training: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("Training: %w", err)
	}

	return dataTrainingToApiTraining(t), nil
}

func (app SwimLogsApp) EditTraining(
	id uuid.UUID,
	t apidef.Training,
) (apidef.TrainingDetail, error) {
	recalcDistanceOnTraining(&t)
	training := trainingToDataTraining(t)

	edited, err := app.pool.EditTraining(id, training)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.TrainingDetail{}, fmt.Errorf("EditTraining: %w", ErrNotFound)
	} else if err != nil {
		return apidef.TrainingDetail{}, fmt.Errorf("EditTraining: %w", err)
	}

	return trainingToDetail(edited), nil
}
