package app

import (
	"fmt"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func New(pool *data.PostgresDbPool) SwimLogsApp {
	return SwimLogsApp{pool}
}

type SwimLogsApp struct {
	pool *data.PostgresDbPool
}

func (app SwimLogsApp) CreateTraining(
	newTraining apidef.NewTraining,
) (apidef.TrainingDetail, error) {
	t := newTrainingToDataTraining(newTraining)
	t, err := app.pool.PersistTraining(t)
	if err != nil {
		return apidef.TrainingDetail{}, fmt.Errorf("CreateTraining: %w", err)
	}

	return trainingToDeatil(t), nil
}
