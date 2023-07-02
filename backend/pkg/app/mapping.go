package app

import (
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
)

func mapDataTrainingsToApiTrainingDetails(ts []data.Training) []openapi.TrainingDetail {
	return mapSlice[data.Training, openapi.TrainingDetail](ts, dataTrainingIntoApiTrainingDetail)
}

func dataTrainingIntoApiTrainingDetail(t data.Training) openapi.TrainingDetail {
	return openapi.TrainingDetail{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
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

// mapSlice is a generic function that maps a source slice of type S to a target
// slice of type T
func mapSlice[S, T any](source []S, f func(S) T) []T {
	result := make([]T, 0, len(source))
	for _, e := range source {
		result = append(result, f(e))
	}
	return result
}
