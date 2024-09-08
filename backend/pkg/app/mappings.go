package app

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"

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
		SetType:        string(s.Type),
		SetOrder:       s.SetOrder,
		StyleId:        s.StyleId,
	}

	if s.Group.IsSpecified() {
		group := s.Group.MustGet()
		set.Group = (*string)(&group)
	}

	if s.Intensity.IsSpecified() {
		intensity := s.Intensity.MustGet()
		set.Intensity = (*string)(&intensity)
	}

	if s.Progression.IsSpecified() {
		progression := s.Progression.MustGet()
		set.Progression = (*string)(&progression)
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

	if startType == string(api.LastFinishesTypeLastFinishes) {
		return
	}

	// this is AsInterval, because both pause and interval have seconds,
	// and if this was pause, then the pause type is set from discriminator
	secs, err := start.AsInterval()
	if err != nil {
		slog.Error(
			"received error from start.AsInterval, not setting seconds value",
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

	if typ == string(api.LastFinishesTypeLastFinishes) {
		start.FromLastFinishes(api.LastFinishes{Type: api.LastFinishesTypeLastFinishes})
		return &start
	}
	if seconds == nil {
		slog.Error("mapDataStart seconds should not be nil", slog.String("type", typ))
		return nil
	}

	if typ == string(api.PauseTypePause) {
		start.FromPause(api.Pause{Type: api.PauseTypePause, Seconds: *seconds})
	} else if typ == string(api.IntervalTypeInterval) {
		start.FromInterval(api.Interval{Type: api.IntervalTypeInterval, Seconds: *seconds})
	} else {
		slog.Error("mapDataStart unknown type", "type", typ)
		return nil
	}

	return &start
}

func newCompToDataComp(c api.NewSetComponent, setId uuid.UUID) data.SetComponent {
	comp := data.SetComponent{
		Id:             uuid.New(),
		SetId:          setId,
		Orders:         c.ComponentOrders,
		Repeat:         c.Repeat,
		DistanceMeters: c.DistanceMeters,
		Description:    c.Description,
		IterationOrder: c.IterationOrder,
		StyleId:        c.StyleId,
	}

	if c.Group.IsSpecified() {
		group := c.Group.MustGet()
		comp.Group = (*string)(&group)
	}

	if c.Intensity.IsSpecified() {
		intensity := c.Intensity.MustGet()
		comp.Intensity = (*string)(&intensity)
	}

	if c.Progression.IsSpecified() {
		progression := c.Progression.MustGet()
		comp.Progression = (*string)(&progression)
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
		IsMain:         s.IsMain,
		Type:           api.TypeEnum(s.SetType),
		TotalDistance:  s.Repeat * s.DistanceMeters,
	}

	if s.Group != nil {
		set.Group = nullable.NewNullableWithValue(api.GroupEnum(*s.Group))
	}

	if s.Intensity != nil {
		set.Intensity = nullable.NewNullableWithValue(api.IntensityEnum(*s.Intensity))
	}

	if s.Progression != nil {
		set.Progression = nullable.NewNullableWithValue(api.ProgressionEnum(*s.Progression))
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
	}

	if c.Group != nil {
		comp.Group = nullable.NewNullableWithValue(api.GroupEnum(*c.Group))
	}

	if c.Intensity != nil {
		comp.Intensity = nullable.NewNullableWithValue(api.IntensityEnum(*c.Intensity))
	}

	if c.Progression != nil {
		comp.Progression = nullable.NewNullableWithValue(api.ProgressionEnum(*c.Progression))
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

func editSetToEmptySet(edited api.EditSetRequest) data.EmptyTrainingSet {
	changes := data.EmptyTrainingSet{
		Repeat:         edited.Repeat,
		DistanceMeters: edited.DistanceMeters,
		IsMain:         edited.IsMain,
		SetType:        (*string)(edited.Type),
	}

	if edited.Description.IsNull() {
		changes.Description = &data.NullStringValue
	} else if edited.Description.IsSpecified() {
		desc := edited.Description.MustGet()
		changes.Description = &desc
	}

	if edited.Group.IsNull() {
		changes.Group = &data.NullStringValue
	} else if edited.Group.IsSpecified() {
		g := edited.Group.MustGet()
		changes.Group = (*string)(&g)
	}

	if edited.Equipment != nil && len(*edited.Equipment) != 0 {
		var equipment []string = nil
		for _, e := range *edited.Equipment {
			equipment = append(equipment, string(e))
		}
		changes.Equipment = &equipment
	}

	if edited.Start.IsNull() {
		changes.StartType = &data.NullStringValue
		changes.StartSeconds = &data.NullIntValue
	} else if edited.Start.IsSpecified() {
		start := edited.Start.MustGet()
		mapApiStart(&changes.StartType, &changes.StartSeconds, start)
	}

	if edited.Intensity.IsNull() {
		changes.Intensity = &data.NullStringValue
	} else if edited.Intensity.IsSpecified() {
		intensity := edited.Intensity.MustGet()
		changes.Intensity = (*string)(&intensity)
	}

	if edited.Progression.IsNull() {
		changes.Progression = &data.NullStringValue
	} else if edited.Progression.IsSpecified() {
		progression := edited.Progression.MustGet()
		changes.Progression = (*string)(&progression)
	}

	if edited.StyleId.IsNull() {
		changes.StyleId = &uuid.Nil
	} else if edited.StyleId.IsSpecified() {
		styleId := edited.StyleId.MustGet()
		changes.StyleId = &styleId
	}

	return changes
}
