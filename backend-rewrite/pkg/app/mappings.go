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

func trainingToDetail(t data.Training) apidef.TrainingDetail {
	return apidef.TrainingDetail{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
	}
}

func dataTrainingToApiTraining(t data.Training) apidef.Training {
	return apidef.Training{
		Id:            t.Id,
		DurationMin:   t.DurationMin,
		Start:         t.Start,
		TotalDistance: t.TotalDistance,
		Sets:          dataSetsToApiSets(t.Sets),
	}
}

func dataSetsToApiSets(sets []data.TrainingSet) []apidef.TrainingSet {
	apiSets := make([]apidef.TrainingSet, 0, len(sets))
	for _, set := range sets {
		apiSets = append(apiSets, dataSetToApiSet(set))
	}
	return apiSets
}

func dataSetToApiSet(s data.TrainingSet) apidef.TrainingSet {
	set := apidef.TrainingSet{
		Id:             s.Id,
		SetOrder:       s.SetOrder,
		Repeat:         s.Repeat,
		Description:    s.Description,
		DistanceMeters: s.DistanceMeters,
		StartType:      apidef.StartTypeEnum(s.StartType),
		StartSeconds:   s.StartSeconds,
		TotalDistance:  s.TotalDistance,
	}

	if s.Equipment == nil {
		return set
	}

	var equipment []apidef.EquipmentEnum
	for _, e := range *s.Equipment {
		equipment = append(equipment, apidef.EquipmentEnum(e))
	}
	set.Equipment = &equipment
	return set
}

func trainingToDataTraining(t apidef.Training) data.Training {
	return data.Training{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
		Sets:          setsToDataSets(t.Sets, t.Id),
	}
}

func setsToDataSets(sets []apidef.TrainingSet, tId uuid.UUID) []data.TrainingSet {
	dataSets := make([]data.TrainingSet, 0, len(sets))
	for _, set := range sets {
		dataSets = append(dataSets, setToDataSet(set, tId))
	}
	return dataSets
}

func setToDataSet(set apidef.TrainingSet, tId uuid.UUID) data.TrainingSet {
	var equipment []string = nil
	if set.Equipment != nil {
		for _, e := range *set.Equipment {
			equipment = append(equipment, string(e))
		}
	}

	ts := data.TrainingSet{
		Id:             set.Id,
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
