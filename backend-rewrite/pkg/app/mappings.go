package app

import (
	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func newTrainingToDataTraining(nt apidef.NewTraining) data.Training {
	id := uuid.New()
	return data.Training{
		Id:            id,
		Start:         nt.Start,
		DurationMin:   nt.DurationMin,
		TotalDistance: nt.TotalDistance,
		Sets:          newSetsToDataSets(nt.Sets, id),
	}
}

func newSetsToDataSets(sets []apidef.NewTrainingSet, tId uuid.UUID) []data.TrainingSet {
	dataSets := make([]data.TrainingSet, 0, len(sets))
	for _, set := range sets {
		dataSets = append(dataSets, newSetToDataSet(set, tId))
	}
	return dataSets
}

func newSetToDataSet(set apidef.NewTrainingSet, tId uuid.UUID) data.TrainingSet {
	var equipment []string = nil
	if set.Equipment != nil {
		for _, e := range *set.Equipment {
			equipment = append(equipment, string(e))
		}
	}

	ts := data.TrainingSet{
		Id:             uuid.New(),
		TrainingId:     tId,
		SetOrder:       set.SetOrder,
		TotalDistance:  set.TotalDistance,
		Repeat:         set.Repeat,
		DistanceMeters: set.DistanceMeters,
		Description:    set.Description,
		StartType:      string(set.StartType),
		StartSeconds:   set.StartSeconds,
		Equipment:      &equipment,
	}

	if set.StartType == apidef.None {
		ts.StartSeconds = nil
	}

	return ts
}

func trainingToDeatil(t data.Training) apidef.TrainingDetail {
	return apidef.TrainingDetail{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
	}
}
