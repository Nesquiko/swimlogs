package app

import (
	"context"
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
	ctx context.Context,
	newTraining apidef.NewTraining,
) (apidef.TrainingSummary, error) {
	validationErr := validateNewTraining(newTraining)
	if validationErr != nil {
		return apidef.TrainingSummary{}, validationErr
	}

	t := newTrainingToDataTraining(newTraining)
	t.Start = t.Start.Truncate(time.Minute)

	t, err := app.pool.PersistTraining(ctx, t)
	if err != nil {
		return apidef.TrainingSummary{}, fmt.Errorf("CreateTraining: %w", err)
	}

	return trainingToSummary(t), nil
}

func (app SwimLogsApp) DeleteTraining(ctx context.Context, id uuid.UUID) error {
	err := app.pool.DeleteTraining(ctx, id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return fmt.Errorf("DeleteTraining not found: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}
	return nil
}

func (app SwimLogsApp) TrainingSummariesPage(
	ctx context.Context,
	page, pageSize int,
) ([]apidef.TrainingSummary, int, error) {
	summariesPage, total, err := app.pool.TrainingSummaries(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummariesPage: %w", err)
	}
	summaries := make([]apidef.TrainingSummary, len(summariesPage))
	for i, d := range summariesPage {
		summaries[i] = trainingToSummary(d)
	}

	return summaries, total, nil
}

func (app SwimLogsApp) TrainingSummariesCurrentWeek(
	ctx context.Context,
) ([]apidef.TrainingSummary, error) {
	startOfWeek, endOfWeek := GetWeekRange(time.Now())

	summariesInRane, err := app.pool.TrainingSummariesInRange(ctx, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("TrainingSummariesCurrentWeek: %w", err)
	}

	summaries := make([]apidef.TrainingSummary, len(summariesInRane))
	for i, ts := range summariesInRane {
		summaries[i] = trainingToSummary(ts)
	}
	return summaries, nil
}

func (app SwimLogsApp) Training(ctx context.Context, id uuid.UUID) (apidef.Training, error) {
	t, err := app.pool.Training(ctx, id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("Training not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("Training: %w", err)
	}

	return dataTrainingToApiTraining(t), nil
}

func (app SwimLogsApp) EditTrainingSession(
	ctx context.Context,
	id uuid.UUID,
	session apidef.EditSessionRequest,
) (apidef.TrainingSummary, error) {
	validationErr := validateEditSessionRequest(session)
	if validationErr != nil {
		return apidef.TrainingSummary{}, validationErr
	}

	edited, err := app.pool.EditTrainingSession(ctx, id, struct {
		DurationMin *int
		Start       *time.Time
	}(session))
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.TrainingSummary{}, fmt.Errorf(
			"EditTrainingSession not found: %w",
			ErrNotFound,
		)
	} else if err != nil {
		return apidef.TrainingSummary{}, fmt.Errorf("EditTrainingSession: %w", err)
	}

	return trainingToSummary(edited), nil
}

func (app SwimLogsApp) DeleteSet(ctx context.Context, trainingId, setId uuid.UUID) error {
	err := app.pool.DeleteSet(ctx, trainingId, setId)
	if errors.Is(err, data.ErrRowsNotFound) {
		return fmt.Errorf("DeleteSet not found: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("DeleteSet: %w", err)
	}
	return nil
}

func GetWeekRange(t time.Time) (time.Time, time.Time) {
	year, week := t.ISOWeek()
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())

	// Find the first Monday of the year
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, 1)
	}

	// Add the number of weeks to the first Monday to get the start of the current week
	start = start.AddDate(0, 0, (week-1)*7)

	// The end of the week is 7 days after the start
	end := start.AddDate(0, 0, 7)

	return start, end
}
