package app

import (
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
)

func dataTrainingIntoApiTrainingDetail(t data.Training) openapi.TrainingDetail {
	return openapi.TrainingDetail{
		Id: t.Id,
	}
}

func apiNewTrainingIntoDataTraining(t openapi.NewTraining) data.Training {
	id := uuid.New()
	return data.Training{
		Id:            id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
		Sets:          apiNewSetsIntoDataSets(t.Sets, id),
	}
}

func apiNewSetsIntoDataSets(sets []openapi.NewTrainingSet, tId uuid.UUID) []data.TrainingSet {
	// will be bigger, because some sets will have subsets
	dataSets := make([]data.TrainingSet, 0, len(sets))
	for _, set := range sets {
		dataSets = append(dataSets, nestedApiNewSetIntoFlatDataSet(set, tId)...)
	}
	return dataSets
}

func nestedApiNewSetIntoFlatDataSet(s openapi.NewTrainingSet, tId uuid.UUID) []data.TrainingSet {
	parent := apiNewSetIntoDataSet(s, *s.SetOrder, tId)
	dataSets := []data.TrainingSet{parent}
	if s.SubSets == nil {
		return dataSets
	}

	for _, subSet := range *s.SubSets {
		dataSet := apiNewSetIntoDataSet(subSet, parent.SetOrder, tId)
		dataSet.ParentSetId = &parent.Id
		dataSets = append(dataSets, dataSet)
	}
	return dataSets
}

func apiNewSetIntoDataSet(
	s openapi.NewTrainingSet,
	parentSetOrder int,
	tId uuid.UUID,
) data.TrainingSet {
	if s.SetOrder == nil {
		s.SetOrder = asPtr(parentSetOrder)
	}

	return data.TrainingSet{
		Id:             uuid.New(),
		TrainingId:     tId,
		SetOrder:       *s.SetOrder,
		SubSetOrder:    s.SubSetOrder,
		TotalDistance:  s.TotalDistance,
		Repeat:         s.Repeat,
		DistanceMeters: s.DistanceMeters,
		Description:    s.Description,
		StartType:      string(s.StartType),
		StartSeconds:   s.StartSeconds,
	}
}

func asPtr[T any](t T) *T {
	return &t
}
