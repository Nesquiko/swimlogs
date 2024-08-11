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
	params apidef.SummariesPageParams,
) ([]apidef.TrainingSummary, int, error) {
	summariesPage, total, err := app.pool.TrainingSummaries(
		ctx,
		params.Page,
		params.PageSize,
		params.From,
		params.Until,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummariesPage: %w", err)
	}
	summaries := make([]apidef.TrainingSummary, len(summariesPage))
	for i, d := range summariesPage {
		summaries[i] = trainingToSummary(d)
	}

	return summaries, total, nil
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

func (app SwimLogsApp) DeleteSet(ctx context.Context, id uuid.UUID) error {
	err := app.pool.DeleteSet(ctx, id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return fmt.Errorf("DeleteSet not found: %w", ErrNotFound)
	} else if err != nil {
		return fmt.Errorf("DeleteSet: %w", err)
	}
	return nil
}

func (app SwimLogsApp) EditSet(
	ctx context.Context,
	id uuid.UUID,
	edited apidef.EditSetRequest,
) (apidef.TrainingSet, int, error) {
	validationErr := validateEditSetRequest(edited)
	if validationErr != nil {
		return apidef.TrainingSet{}, 0, validationErr
	}

	var equipment []string = nil
	if edited.Equipment != nil {
		for _, e := range *edited.Equipment {
			equipment = append(equipment, string(e))
		}
	}

	set, trainingTotalDist, err := app.pool.EditSet(ctx, id, struct {
		Repeat         *int
		DistanceMeters *int
		Description    *string
		Equipment      *[]string
		StartType      *string
		StartSeconds   *int
		Group          *string
		IsMain         *bool
	}{
		Repeat:         edited.Repeat,
		DistanceMeters: edited.DistanceMeters,
		Description:    edited.Description,
		Equipment:      &equipment,
		StartType:      (*string)(edited.StartType),
		StartSeconds:   edited.StartSeconds,
		Group:          (*string)(edited.Group),
		IsMain:         edited.IsMain,
	})
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.TrainingSet{}, 0, fmt.Errorf("EditSet not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.TrainingSet{}, 0, fmt.Errorf("EditSet: %w", err)
	}

	return dataSetToApiSet(set), trainingTotalDist, nil
}

func (app SwimLogsApp) MoveSet(
	ctx context.Context,
	id uuid.UUID,
	newSetOrder int,
) (apidef.Training, error) {
	validationErr := validateNewSetOrder(newSetOrder)
	if validationErr != nil {
		return apidef.Training{}, validationErr
	}

	err := app.pool.MoveSet(ctx, id, newSetOrder)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("MoveSet not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("MoveSet: %w", err)
	}

	trainingId, err := app.pool.TrainingIdBySetId(ctx, id)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("MoveSet training id not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("MoveSet: %w", err)
	}

	t, err := app.pool.Training(ctx, trainingId)
	if errors.Is(err, data.ErrRowsNotFound) {
		return apidef.Training{}, fmt.Errorf("MoveSet not found: %w", ErrNotFound)
	} else if err != nil {
		return apidef.Training{}, fmt.Errorf("MoveSet: %w", err)
	}

	return dataTrainingToApiTraining(t), nil
}
