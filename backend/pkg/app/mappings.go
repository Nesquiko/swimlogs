package app

import (
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func newTrainingToDataTraining(nt api.NewTraining) data.Training {
	t := data.Training{
		Id:          uuid.New(),
		Start:       nt.Start.Truncate(time.Minute),
		DurationMin: nt.DurationMinutes,
	}

	sets := make([]data.TrainingSet, len(nt.Sets))
	for i, s := range nt.Sets {
		sets[i] = newSetToDataSet(s, t.Id)
	}

	t.Sets = sets

	return t
}

func newSetToDataSet(s api.NewTrainingSet, tId uuid.UUID) data.TrainingSet {
	set := data.TrainingSet{
		Id:             uuid.New(),
		TrainingId:     tId,
		Repeat:         s.Repeat,
		DistanceMeters: s.DistanceMeters,
		Description:    s.Description,
		Group:          (*string)(s.Group),
		SetType:        string(s.Type),
		Intensity:      (*string)(s.Intensity),
		Progression:    (*string)(s.Progression),
		SetOrder:       s.SetOrder,
		StyleId:        s.StyleId,
	}

	if s.Equipment != nil && len(*s.Equipment) != 0 {
		equipment := make([]string, len(*s.Equipment))
		for i, e := range *s.Equipment {
			equipment[i] = string(e)
		}
		set.Equipment = &equipment
	}

	if s.Start != nil {
		mapApiStart(&set.StartType, &set.StartSeconds, *s.Start)
	}

	if s.IsMain != nil {
		set.IsMain = *s.IsMain
	} else {
		set.IsMain = false
	}

	if s.Components != nil {
		comps := make([]data.SetComponent, len(*s.Components))
		for i, comp := range *s.Components {
			comps[i] = newCompToDataComp(comp, set.Id)
		}
		set.Components = &comps
	}

	return set
}

func mapApiStart(typ **string, seconds **int, start api.Start) {
	startType, err := start.Discriminator()
	if err != nil {
		slog.Error(
			"received error from start.Discriminator, not setting any start values",
			slog.String("error", err.Error()),
		)
		return
	}

	if *typ == nil {
		*typ = new(string)
	}
	**typ = startType

	if startType == string(api.LastFinishes) {
		return
	}

	secs, err := start.AsStartWithSeconds()
	if err != nil {
		slog.Error(
			"received error from start.AsStartWithSeconds, not setting seconds value",
			slog.String("error", err.Error()),
		)
		return
	}

	if *seconds == nil {
		*seconds = new(int)
	}
	**seconds = secs.Seconds
}

func mapDataStart(typ string, seconds *int) *api.Start {
	start := api.Start{}

	if typ == string(api.LastFinishes) {
		start.FromStartWithoutSeconds(api.StartWithoutSeconds{Type: api.LastFinishes})
		return &start
	}

	withSeconds := api.StartWithSeconds{
		Type: api.StartWithSecondsType(typ),
	}
	if seconds == nil {
		slog.Error("mapDataStart seconds should not be nil", slog.String("type", typ))
	} else {
		withSeconds.Seconds = *seconds
	}

	start.FromStartWithSeconds(withSeconds)
	return &start
}

func newCompToDataComp(c api.NewSetComponent, setId uuid.UUID) data.SetComponent {
	comp := data.SetComponent{
		Id:             uuid.New(),
		SetId:          setId,
		Orders:         c.ComponentOrders,
		Repeat:         c.Repeat,
		DistanceMeters: c.DistanceMeters,
		Intensity:      (*string)(c.Intensity),
		Progression:    (*string)(c.Progression),
		Description:    c.Description,
		IterationOrder: c.IterationOrder,
		StyleId:        c.StyleId,
	}

	if c.Start != nil {
		mapApiStart(&comp.StartType, &comp.StartSeconds, *c.Start)
	}

	if c.Equipment != nil && len(*c.Equipment) != 0 {
		equipment := make([]string, len(*c.Equipment))
		for i, e := range *c.Equipment {
			equipment[i] = string(e)
		}
		comp.Equipment = &equipment
	}

	return comp
}

func dataTrainingToApiTraining(t data.Training) api.Training {
	training := api.Training{
		Id:              t.Id,
		DurationMinutes: t.DurationMin,
		Start:           t.Start,
	}

	sets := make([]api.TrainingSet, len(t.Sets))
	totalDistance := 0
	for i, s := range t.Sets {
		set := dataSetToApiSet(s)
		sets[i] = set
		totalDistance += (s.Repeat * s.DistanceMeters)
	}
	training.Sets = sets
	training.TotalDistance = totalDistance

	return training
}

func dataSetToApiSet(s data.TrainingSet) api.TrainingSet {
	set := api.TrainingSet{
		Id:             s.Id,
		SetOrder:       s.SetOrder,
		Repeat:         s.Repeat,
		DistanceMeters: s.DistanceMeters,
		Description:    s.Description,
		Group:          (*api.GroupEnum)(s.Group),
		IsMain:         s.IsMain,
		Type:           api.TypeEnum(s.SetType),
		Intensity:      (*api.IntensityEnum)(s.Intensity),
		Progression:    (*api.ProgressionEnum)(s.Progression),
		TotalDistance:  s.Repeat * s.DistanceMeters,
	}

	if s.StartType != nil {
		set.Start = mapDataStart(*s.StartType, s.StartSeconds)
	}

	if s.Equipment != nil && len(*s.Equipment) != 0 {
		equipment := make([]api.EquipmentEnum, len(*s.Equipment))
		for i, e := range *s.Equipment {
			equipment[i] = api.EquipmentEnum(e)
		}
		set.Equipment = &equipment
	}

	if s.Components != nil && len(*s.Components) != 0 {
		comps := make([]api.SetComponent, len(*s.Components))
		for i, c := range *s.Components {
			comps[i] = dataCompToApiComp(c)
		}
		set.Components = &comps
	}

	if s.Style != nil {
		set.Style = dataStyleToApiStyle(*s.Style)
	}

	return set
}

func dataCompToApiComp(c data.SetComponent) api.SetComponent {
	comp := api.SetComponent{
		Id:              c.Id,
		ComponentOrders: c.Orders,
		IterationOrder:  c.IterationOrder,
		Repeat:          c.Repeat,
		DistanceMeters:  c.DistanceMeters,
		Description:     c.Description,
		Group:           (*api.GroupEnum)(c.Group),
		Intensity:       (*api.IntensityEnum)(c.Intensity),
		Progression:     (*api.ProgressionEnum)(c.Progression),
	}

	if c.StartType != nil {
		comp.Start = mapDataStart(*c.StartType, c.StartSeconds)
	}

	if c.Equipment != nil && len(*c.Equipment) != 0 {
		equipment := make([]api.EquipmentEnum, len(*c.Equipment))
		for i, e := range *c.Equipment {
			equipment[i] = api.EquipmentEnum(e)
		}
		comp.Equipment = &equipment
	}

	if c.Style != nil {
		comp.Style = dataStyleToApiStyle(*c.Style)
	}

	return comp
}

func dataStyleToApiStyle(s data.Style) *api.Style {
	return &api.Style{
		Id:                  s.Id,
		Name:                s.Name,
		ExecriseName:        s.ExeciseName,
		ExecriseDescription: s.ExeciseDescription,
	}
}

func trainingToSummary(t data.Training, withMainSets bool) api.TrainingSummary {
	totalDistance := 0

	ts := api.TrainingSummary{
		Id:              t.Id,
		Start:           t.Start,
		DurationMinutes: t.DurationMin,
	}

	if len(t.Sets) == 0 {
		totalDistance = t.TotalDistance
	} else {
		mainSets := make([]api.TrainingSet, 0)
		for _, s := range t.Sets {
			totalDistance += s.Repeat * s.DistanceMeters
			if s.IsMain && withMainSets {
				mainSets = append(mainSets, dataSetToApiSet(s))
			}
		}

		if withMainSets {
			ts.MainSets = &mainSets
		}
	}
	ts.TotalDistance = totalDistance

	return ts
}
