package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const selectTraining = `
select
    t.id, t.start, t.duration_min, t.created_at, t.modified_at,

    s.id, s.set_order, s.repeat, s.distance_meters,  s.description,
    s.start_type, s.start_seconds, s.equipment, s."group", s.is_main,
    s.type, s.intensity, s.progression,
    
    set_style.id, set_style.name, set_style.exercise_name, set_style.exercise_description,

    c.id, c.component_orders, c.iteration_order, c.repeat,
    c.distance_meters, c.start_type, c.start_seconds, c.intensity,
    c.progression, c.description, c.equipment, c."group",
    
    comp_style.id,  comp_style.name,  comp_style.exercise_name,  comp_style.exercise_description
from trainings t
         join sets s on t.id = s.training_id
         left join styles set_style on s.style_id = set_style.id
         left join set_components c on s.id = c.set_id
         left join styles comp_style on c.style_id = comp_style.id
where t.id = $1
order by s.set_order
`

func (pgs *PostgresDbPool) TrainingById(ctx context.Context, id uuid.UUID) (Training, error) {
	rows, err := pgs.pool.Query(ctx, selectTraining, id)
	if err != nil {
		return Training{}, fmt.Errorf("TrainingById query error: %w", err)
	}

	var t Training
	for rows.Next() {
		err := scan(&t, rows)
		if err != nil {
			return Training{}, fmt.Errorf("Training scanning error: %w", err)
		}
	}

	if rows.Err() != nil {
		return Training{}, fmt.Errorf("TrainingById rows error: %w", rows.Err())
	}
	rows.Close()

	if rows.CommandTag().RowsAffected() == 0 {
		return Training{}, fmt.Errorf("TrainingById id not found: %w", ErrRowsNotFound)
	}

	return t, nil
}

func scan(training *Training, rows pgx.Rows) error {
	t, s, c, err := scanAll(rows)

	defer func() { *training = t }()

	if err != nil {
		return fmt.Errorf("scan scan all: %w", err)
	}
	t.Sets = training.Sets

	if len(t.Sets) == 0 {
		t.Sets = make([]TrainingSet, 1)
		t.Sets[0] = s
	}

	if s.Id != t.Sets[len(t.Sets)-1].Id {
		t.Sets = append(t.Sets, s)
	}

	if !c.Id.Valid {
		return nil
	}

	lastSet := &t.Sets[len(t.Sets)-1]
	if lastSet.Components == nil {
		lastSet.Components = &[]SetComponent{}
	}
	comps := *lastSet.Components
	comps = append(*lastSet.Components, c)
	lastSet.Components = &comps

	return nil
}

func scanAll(rows pgx.Rows) (Training, TrainingSet, SetComponent, error) {
	var row JoinedRow
	dests := []any{
		&row.TrainingId, &row.TrainingStart, &row.TrainingDurationMin,
		&row.TrainingCreatedAt, &row.TrainingModifiedAt, &row.SetId,
		&row.SetSetOrder, &row.SetRepeat, &row.SetDistanceMeters,
		&row.SetDescription, &row.SetStartType, &row.SetStartSeconds,
		&row.SetEquipment, &row.SetGroup, &row.SetIsMain, &row.SetSetType,
		&row.SetIntensity, &row.SetProgression, &row.SetStyleId,
		&row.SetStyleName, &row.SetStyleExeciseName, &row.SetStyleExeciseDescription,
		&row.CompId, &row.CompOrders, &row.CompIterationOrder, &row.CompRepeat,
		&row.CompDistanceMeters, &row.CompStartType, &row.CompStartSeconds,
		&row.CompIntensity, &row.CompProgression, &row.CompDescription,
		&row.CompEquipment, &row.CompGroup, &row.CompStyleId, &row.CompStyleName,
		&row.CompStyleExeciseName, &row.CompStyleExeciseDescription,
	}

	err := rows.Scan(dests...)
	if err != nil {
		return Training{}, TrainingSet{}, SetComponent{}, fmt.Errorf("scanAll: %w", err)
	}

	return row.ToTraining(), row.ToTrainingSet(), row.ToSetComponent(), nil
}

type JoinedRow struct {
	TrainingId          uuid.UUID
	TrainingStart       time.Time
	TrainingDurationMin int
	TrainingCreatedAt   time.Time
	TrainingModifiedAt  time.Time

	SetId             uuid.UUID
	SetSetOrder       int
	SetRepeat         int
	SetDistanceMeters int
	SetDescription    *string
	SetStartType      *string
	SetStartSeconds   *int
	SetEquipment      *[]string
	SetGroup          *string
	SetIsMain         bool
	SetSetType        string
	SetIntensity      *string
	SetProgression    *string

	SetStyleId                 uuid.UUID
	SetStyleName               sql.NullString
	SetStyleExeciseName        *string
	SetStyleExeciseDescription *string

	CompId             uuid.NullUUID
	CompOrders         []int
	CompIterationOrder *int
	CompRepeat         sql.NullInt64
	CompDistanceMeters sql.NullInt64
	CompStartType      *string
	CompStartSeconds   *int
	CompIntensity      *string
	CompProgression    *string
	CompDescription    *string
	CompEquipment      *[]string
	CompGroup          *string

	CompStyleId                 uuid.UUID
	CompStyleName               sql.NullString
	CompStyleExeciseName        *string
	CompStyleExeciseDescription *string
}

func (r JoinedRow) ToTraining() Training {
	return Training{
		Id:          r.TrainingId,
		Start:       r.TrainingStart,
		DurationMin: r.TrainingDurationMin,
		CreatedAt:   r.TrainingCreatedAt,
		ModifiedAt:  r.TrainingModifiedAt,
	}
}

func (r JoinedRow) ToTrainingSet() TrainingSet {
	s := TrainingSet{
		Id:             r.SetId,
		TrainingId:     r.TrainingId,
		SetOrder:       r.SetSetOrder,
		Repeat:         r.SetRepeat,
		DistanceMeters: r.SetDistanceMeters,
		Description:    r.SetDescription,
		Equipment:      r.SetEquipment,
		StartType:      r.SetStartType,
		StartSeconds:   r.SetStartSeconds,
		Group:          r.SetGroup,
		IsMain:         r.SetIsMain,
		SetType:        r.SetSetType,
		Intensity:      r.SetIntensity,
		Progression:    r.SetProgression,
	}

	if r.SetStyleId != uuid.Nil {
		s.Style = &Style{
			Id:                 r.SetStyleId,
			Name:               r.SetStyleName.String,
			ExeciseName:        r.SetStyleExeciseName,
			ExeciseDescription: r.SetStyleExeciseDescription,
		}
		s.StyleId = &s.Style.Id
	}

	return s
}

func (r JoinedRow) ToSetComponent() SetComponent {
	c := SetComponent{
		Id:             r.CompId,
		SetId:          r.SetId,
		Orders:         r.CompOrders,
		IterationOrder: r.CompIterationOrder,
		Repeat:         int(r.CompRepeat.Int64),
		DistanceMeters: int(r.CompDistanceMeters.Int64),
		Equipment:      r.CompEquipment,
		StartType:      r.CompStartType,
		StartSeconds:   r.CompStartSeconds,
		Group:          r.CompGroup,
		Intensity:      r.CompIntensity,
		Progression:    r.CompProgression,
		Description:    r.CompDescription,
	}

	if r.CompStyleId != uuid.Nil {
		c.Style = &Style{
			Id:                 r.CompStyleId,
			Name:               r.CompStyleName.String,
			ExeciseName:        r.CompStyleExeciseName,
			ExeciseDescription: r.CompStyleExeciseDescription,
		}
		c.StyleId = &c.Style.Id
	}

	return c
}
