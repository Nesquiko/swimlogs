package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pool *PostgresDbPool) PersistTraining(ctx context.Context, t Training) error {
	return Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		err := persistTraining(ctx, t, tx)
		if err != nil {
			return fmt.Errorf("PersistTraining tx: %w", err)
		}
		return nil
	})
}

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

func (pool *PostgresDbPool) DeleteTraining(ctx context.Context, id uuid.UUID) error {
	err := Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		return deleteTraining(ctx, id, tx)
	})
	if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}
	return nil
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
	var row joinedRow
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

type joinedRow struct {
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

func (r joinedRow) ToTraining() Training {
	return Training{
		Id:          r.TrainingId,
		Start:       r.TrainingStart,
		DurationMin: r.TrainingDurationMin,
		CreatedAt:   r.TrainingCreatedAt,
		ModifiedAt:  r.TrainingModifiedAt,
	}
}

func (r joinedRow) ToTrainingSet() TrainingSet {
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

func (r joinedRow) ToSetComponent() SetComponent {
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

const insertTraining = `
insert into trainings (id, start, duration_min, created_at, modified_at)
values ($1, $2, $3, now(), now())
`

func persistTraining(ctx context.Context, t Training, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, insertTraining, t.Id, t.Start, t.DurationMin)
	if err != nil {
		return fmt.Errorf("persistTraining persisting training: %w", err)
	}

	for i, s := range t.Sets {
		err := persistSet(ctx, tx, s)
		if err != nil {
			return fmt.Errorf("persistTraining set %d: %w", i, err)
		}
	}

	return nil
}

func trainingIdBySetId(ctx context.Context, setId uuid.UUID, tx pgx.Tx) (uuid.UUID, error) {
	var trainingId uuid.UUID
	err := tx.QueryRow(ctx, "select training_id from sets where id = $1", setId).Scan(&trainingId)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, fmt.Errorf("trainingIdBySetId not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return uuid.UUID{}, fmt.Errorf("trainingIdBySetId: %w", err)
	}

	return trainingId, nil
}

func deleteTraining(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {
	ct, err := tx.Exec(ctx, "delete from trainings where id = $1", id)
	if err != nil {
		return fmt.Errorf("deleteTraining: %w", err)
	} else if ct.RowsAffected() == 0 {
		return fmt.Errorf("deleteTraining training doesnt exist: %w", ErrRowsNotFound)
	}
	return nil
}

const selectTrainingSummariesPage = `
with filtered as (
    select t.id, t.start, t.duration_min, t.created_at, t.modified_at, sum(s.repeat * s.distance_meters)
    from trainings t join sets s on t.id = s.training_id
    where t.start between $3 and $4
    group by t.id, t.start, t.duration_min, t.created_at, t.modified_at)
select t.*, count(t.id) over (), s.*
from filtered t left join sets s on s.training_id = t.id and s.is_main = true
order by t.start desc, t.duration_min, t.created_at
limit $1 offset $2
`

type EmptyTrainingSet struct {
	Id             *uuid.UUID
	TrainingId     *uuid.UUID
	SetOrder       *int
	Repeat         *int
	DistanceMeters *int
	Description    *string
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	Group          *string
	IsMain         *bool
}

func (s EmptyTrainingSet) isFilled() bool {
	return s.Id != nil && s.TrainingId != nil && s.SetOrder != nil && s.Repeat != nil &&
		s.DistanceMeters != nil && s.IsMain != nil
}

func (s EmptyTrainingSet) intoTrainingSet() TrainingSet {
	if s.Id == nil {
		slog.Error(
			"EmptyTrainingSet.intoTrainingSe: this should't be called, on non validated EmptyTrainingSet",
		)
		return TrainingSet{}
	}

	return TrainingSet{
		Id:             *s.Id,
		TrainingId:     *s.TrainingId,
		SetOrder:       *s.SetOrder,
		Repeat:         *s.Repeat,
		DistanceMeters: *s.DistanceMeters,
		Description:    s.Description,
		Equipment:      s.Equipment,
		StartType:      s.StartType,
		StartSeconds:   s.StartSeconds,
		Group:          s.Group,
		IsMain:         *s.IsMain,
	}
}

// TODO from and until filters
func (psg *PostgresDbPool) TrainingSummaries(
	ctx context.Context,
	page, pageSize int,
	from, until *time.Time,
) ([]Training, int, error) {
	ts := make([]Training, 0)

	fromStr := "-infinity"
	untilStr := "infinity"

	if from != nil {
		fromStr = from.Format(time.RFC3339)
	}
	if until != nil {
		untilStr = until.Format(time.RFC3339)
	}

	rows, err := psg.pool.Query(
		ctx,
		selectTrainingSummariesPage,
		pageSize,
		page*pageSize,
		fromStr,
		untilStr,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummaries query error: %w", err)
	}
	defer rows.Close()

	count := 0
	lastTrainingId := uuid.UUID{}
	for rows.Next() {
		t := Training{}
		s := EmptyTrainingSet{}

		scanArgs := []any{
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.CreatedAt,
			&t.ModifiedAt,
			&t.TotalDistance,
			&count,
			&s.Id,
			&s.TrainingId,
			&s.SetOrder,
			&s.Repeat,
			&s.DistanceMeters,
			&s.Description,
			&s.StartType,
			&s.StartSeconds,
			&s.Equipment,
			&s.Group,
			&s.IsMain,
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, fmt.Errorf("TrainingSummaries scanning row: %w", err)
		}

		if lastTrainingId != t.Id {
			if s.isFilled() {
				t.Sets = append(t.Sets, s.intoTrainingSet())
			}
			ts = append(ts, t)
			lastTrainingId = t.Id
		} else if s.isFilled() {
			ts[len(ts)-1].Sets = append(ts[len(ts)-1].Sets, s.intoTrainingSet())
		}
	}

	return ts, count, nil
}

func (pool *PostgresDbPool) EditTrainingSession(
	ctx context.Context,
	id uuid.UUID,
	session struct {
		DurationMin *int
		Start       *time.Time
	},
) (Training, error) {
	return TxWithResult(ctx, pool, func(ctx context.Context, tx pgx.Tx) (Training, error) {
		t, err := pool.editTrainingSession(ctx, id, session, tx)
		if err != nil {
			return t, fmt.Errorf("EditTrainingSession tx: %w", err)
		}
		return t, nil
	})
}

const updateTrainingSession = `
with updated as (
    update trainings
        set start = coalesce($2, start),
            duration_min = coalesce($3, duration_min),
            modified_at = now()
        where id = $1
        returning id, start, duration_min, created_at, modified_at)
select t.id, t.start, t.duration_min, t.created_at, t.modified_at,
    sum(s.repeat * s.distance_meters)
from updated t
         join sets s on t.id = s.training_id
group by t.id, t.start, t.duration_min, t.created_at, t.modified_at
`

func (pool *PostgresDbPool) editTrainingSession(
	ctx context.Context,
	id uuid.UUID,
	session struct {
		DurationMin *int
		Start       *time.Time
	},
	tx pgx.Tx,
) (Training, error) {
	t := Training{}
	err := tx.QueryRow(ctx, updateTrainingSession, id, session.Start, session.DurationMin).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.CreatedAt, &t.ModifiedAt, &t.TotalDistance)

	if errors.Is(err, pgx.ErrNoRows) {
		return Training{}, fmt.Errorf("editTrainingSession not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return Training{}, fmt.Errorf("editTrainingSession update training query error: %w", err)
	}

	return t, nil
}

func (pool *PostgresDbPool) EditSet(
	ctx context.Context,
	id uuid.UUID,
	s struct {
		Repeat         *int
		DistanceMeters *int
		Description    *string
		Equipment      *[]string
		StartType      *string
		StartSeconds   *int
		Group          *string
		IsMain         *bool
	},
) (TrainingSet, int, error) {
	result, err := TxWithResult(
		ctx,
		pool,
		func(ctx context.Context, tx pgx.Tx) (struct {
			set               TrainingSet
			trainingTotalDist int
		}, error,
		) {
			set, trainingTotalDist, err := pool.editSet(ctx, id, s, tx)
			return struct {
				set               TrainingSet
				trainingTotalDist int
			}{set, trainingTotalDist}, err
		},
	)
	if err != nil {
		return TrainingSet{}, 0, fmt.Errorf("EditSet: %w", err)
	}

	return result.set, result.trainingTotalDist, nil
}

const updateSet = `
update sets
   set repeat = coalesce($2, repeat),
       distance_meters = coalesce($3, distance_meters),
       description = coalesce($4, description),
       start_type = coalesce($5, start_type),
       start_seconds = coalesce($6, start_seconds),
       equipment = coalesce($7, equipment),
       "group" = coalesce($8, "group"),
       is_main = coalesce($9, is_main)
where id = $1
returning id, training_id, set_order, repeat, distance_meters, description, start_type, start_seconds, equipment, "group", is_main
`

const totalDistanceInTraining = `
select sum(s.repeat * s.distance_meters)
from trainings t join sets s on t.id = s.training_id
where t.id = $1;
`

func (pool *PostgresDbPool) editSet(
	ctx context.Context,
	id uuid.UUID,
	edited struct {
		Repeat         *int
		DistanceMeters *int
		Description    *string
		Equipment      *[]string
		StartType      *string
		StartSeconds   *int
		Group          *string
		IsMain         *bool
	},
	tx pgx.Tx,
) (TrainingSet, int, error) {
	s := TrainingSet{}
	err := tx.QueryRow(
		ctx,
		updateSet,
		id,
		edited.Repeat,
		edited.DistanceMeters,
		edited.Description,
		edited.StartType,
		edited.StartSeconds,
		edited.Equipment,
		edited.Group,
		edited.IsMain,
	).Scan(
		&s.Id,
		&s.TrainingId,
		&s.SetOrder,
		&s.Repeat,
		&s.DistanceMeters,
		&s.Description,
		&s.StartType,
		&s.StartSeconds,
		&s.Equipment,
		&s.Group,
		&s.IsMain,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return TrainingSet{}, 0, fmt.Errorf("editSet not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return TrainingSet{}, 0, fmt.Errorf("editSet update query error: %w, id: %s", err, s.Id)
	}

	trainingId, err := trainingIdBySetId(ctx, id, tx)
	if err != nil {
		return TrainingSet{}, 0, fmt.Errorf("editSet retrieve training id query: %w", err)
	}

	trainingTotalDist := 0
	err = tx.QueryRow(ctx, totalDistanceInTraining, trainingId).Scan(&trainingTotalDist)
	if err != nil {
		return TrainingSet{}, 0, fmt.Errorf("editSet sum query error: %w", err)
	}

	return s, trainingTotalDist, nil
}

const moveSets = `
with to_be_moved as (select * from sets s where s.id = $1)
update sets
set set_order = case when to_be_moved.set_order > $2 then sets.set_order + 1 else greatest(sets.set_order - 1, 0) end
from to_be_moved
where case
          when to_be_moved.set_order > $2 then
              sets.set_order between $2 and to_be_moved.set_order
          else
              sets.set_order between to_be_moved.set_order and $2 end
  and sets.training_id = to_be_moved.training_id;
`

const updateSetOrder = `
update sets
set set_order = $2
where id = $1
`

func (pool *PostgresDbPool) MoveSet(
	ctx context.Context,
	id uuid.UUID,
	newSetOrder int,
) error {
	err := Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		return pool.moveSet(ctx, id, newSetOrder, tx)
	})
	if err != nil {
		return fmt.Errorf("MoveSet: %w", err)
	}
	return nil
}

func (pool *PostgresDbPool) moveSet(
	ctx context.Context,
	id uuid.UUID,
	newSetOrder int,
	tx pgx.Tx,
) error {
	ct, err := tx.Exec(ctx, moveSets, id, newSetOrder)
	if err != nil {
		return fmt.Errorf("moveSet move query: %w", err)
	}

	ct, err = tx.Exec(ctx, updateSetOrder, id, newSetOrder)
	if err != nil {
		return fmt.Errorf("moveSet update query: %w", err)
	} else if ct.RowsAffected() == 0 {
		return fmt.Errorf("moveSet update query set not found: %w", ErrRowsNotFound)
	}
	return nil
}

func (pool *PostgresDbPool) TrainingIdBySetId(
	ctx context.Context,
	setId uuid.UUID,
) (uuid.UUID, error) {
	id, err := TxWithResult(ctx, pool, func(ctx context.Context, tx pgx.Tx) (uuid.UUID, error) {
		return trainingIdBySetId(ctx, setId, tx)
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("TrainingIdBySetId: %w", err)
	}
	return id, nil
}
