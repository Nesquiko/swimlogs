package app

import (
	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func newTrainingToDataTraining(nt apidef.NewTraining) data.Training {
	id := uuid.New()
	return data.Training{
		Id:          id,
		Start:       nt.Start,
		DurationMin: nt.DurationMin,
		Sets:        newSetsToDataSets(nt.Sets, id),
	}
}

func newSetsToDataSets(sets []apidef.NewTrainingSet, tId uuid.UUID) []data.TrainingSet {
	dataSets := make([]data.TrainingSet, 0, len(sets))
	for _, set := range sets {
		s := newSetToDataSet(set, tId)
		dataSets = append(dataSets, s)
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

	isMain := false
	if set.IsMain != nil && *set.IsMain {
		isMain = true
	}
	ts := data.TrainingSet{
		Id:             uuid.New(),
		TrainingId:     tId,
		SetOrder:       set.SetOrder,
		Repeat:         set.Repeat,
		DistanceMeters: set.DistanceMeters,
		Description:    set.Description,
		StartType:      (*string)(set.StartType),
		StartSeconds:   set.StartSeconds,
		Equipment:      &equipment,
		Group:          (*string)(set.Group),
		IsMain:         isMain,
	}

	if set.StartType == nil {
		ts.StartSeconds = nil
	}

	return ts
}

func trainingToSummary(t data.Training) apidef.TrainingSummary {
	totalDistance := 0

	mainSets := make([]apidef.TrainingSet, 0)
	if len(t.Sets) == 0 {
		totalDistance = t.TotalDistance
	} else {
		for _, s := range t.Sets {
			if s.IsMain {
				mainSets = append(mainSets, dataSetToApiSet(s))
			}
			totalDistance += s.Repeat * s.DistanceMeters
		}
	}

	return apidef.TrainingSummary{
		Id:            t.Id,
		Start:         t.Start,
		DurationMin:   t.DurationMin,
		TotalDistance: totalDistance,
		MainSets:      &mainSets,
	}
}

func dataTrainingToApiTraining(t data.Training) apidef.Training {
	sets := dataSetsToApiSets(t.Sets)
	totalDistance := 0
	for _, s := range sets {
		totalDistance += (s.Repeat * s.DistanceMeters)
	}

	return apidef.Training{
		Id:            t.Id,
		DurationMin:   t.DurationMin,
		Start:         t.Start,
		TotalDistance: totalDistance,
		Sets:          sets,
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
		StartType:      (*apidef.StartTypeEnum)(s.StartType),
		StartSeconds:   s.StartSeconds,
		TotalDistance:  s.Repeat * s.DistanceMeters,
		Group:          (*apidef.GroupEnum)(s.Group),
		IsMain:         s.IsMain,
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
		Repeat:         set.Repeat,
		DistanceMeters: set.DistanceMeters,
		Description:    set.Description,
		StartType:      (*string)(set.StartType),
		StartSeconds:   set.StartSeconds,
		Equipment:      &equipment,
		Group:          (*string)(set.Group),
		IsMain:         set.IsMain,
	}

	if set.StartType == nil || *set.StartType == "" {
		ts.StartSeconds = nil
	}

	return ts
}
