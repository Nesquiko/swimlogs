package app

import (
	"time"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/google/uuid"
)

func dataTrainingIntoApiTraining(t data.Training) openapi.Training {
	return openapi.Training{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
		Sets:          dataSetsIntoApiTrainingSets(t.Sets),
	}
}

func dataSetsIntoApiTrainingSets(ts []data.TrainingSet) []openapi.TrainingSet {
	apiSets := make([]openapi.TrainingSet, 0)

	for _, dataSet := range ts {
		apiSet := dataSetIntoApiTrainingSet(dataSet)
		if dataSet.SubSets == nil || len(*dataSet.SubSets) == 0 {
			apiSets = append(apiSets, apiSet)
			continue
		}

		for _, subSet := range *dataSet.SubSets {
			apiSubSet := dataSetIntoApiTrainingSet(subSet)

			if apiSet.SubSets == nil {
				apiSet.SubSets = &[]openapi.TrainingSet{}
			}

			newSubSets := append(*apiSet.SubSets, apiSubSet)
			apiSet.SubSets = &newSubSets
		}
		apiSets = append(apiSets, apiSet)
	}

	return apiSets
}

func dataSetIntoApiTrainingSet(s data.TrainingSet) openapi.TrainingSet {
	return openapi.TrainingSet{
		Id:             s.Id,
		SetOrder:       &s.SetOrder,
		SubSetOrder:    s.SubSetOrder,
		Repeat:         s.Repeat,
		Description:    s.Description,
		DistanceMeters: s.DistanceMeters,
		StartType:      openapi.StartType(s.StartType),
		StartSeconds:   s.StartSeconds,
		TotalDistance:  s.TotalDistance,
	}
}

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

func apiSessionFromDataSession(s data.Session) openapi.Session {
	return openapi.Session{
		Id:          s.Id,
		Day:         openapi.Day(s.Day),
		StartTime:   s.StartTime.Format("15:04"),
		DurationMin: s.DurationMin,
	}
}

func mapDataSessionsToApiSessions(sessions []data.Session) []openapi.Session {
	return mapSlice(sessions, apiSessionFromDataSession)
}

func dataSessionFromNewApiSession(s openapi.CreateSessionJSONBody) data.Session {
	startTime, _ := time.Parse("15:04", s.StartTime)
	return data.Session{
		Day:         string(s.Day),
		StartTime:   startTime,
		DurationMin: s.DurationMin,
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
