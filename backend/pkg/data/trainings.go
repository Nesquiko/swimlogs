package data

import (
	"context"
	"errors"
	"fmt"
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
    select t.id, t.start, t.duration_min, t.created_at, t.modified_at,
    sum(s.repeat * s.distance_meters) as total_distance
        from trainings t join sets s on t.id = s.training_id
    where t.start between $3 and $4
    group by t.id, t.start, t.duration_min, t.created_at, t.modified_at
    order by t.start desc, t.duration_min, t.created_at
    limit $1 offset $2)
select
    t.id, t.start, t.duration_min, t.created_at, t.modified_at, t.total_distance,
    (select count(*) from trainings) as count,

    s.id, s.training_id, s.set_order, s.repeat, s.distance_meters,  s.description,
    s.start_type, s.start_seconds, s.equipment, s."group", s.is_main,
    s.type, s.intensity, s.progression,

    set_style.id, set_style.name, set_style.exercise_name, set_style.exercise_description
from filtered t
    left join sets s on s.training_id = t.id and s.is_main = true
    left join styles set_style on s.style_id = set_style.id
order by t.start desc, t.duration_min, t.created_at
`

func (psg *PostgresDbPool) TrainingSummaries(
	ctx context.Context,
	page, pageSize int,
	from, until *time.Time,
) ([]Training, int, error) {
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
	defer rows.Close()
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummaries query error: %w", err)
	}

	ts := make([]Training, 0)
	count := 0
	lastTrainingId := uuid.UUID{}
	for rows.Next() {
		t := Training{}
		s := EmptyTrainingSet{}

		scanArgs := []any{
			&t.Id, &t.Start, &t.DurationMin, &t.CreatedAt, &t.ModifiedAt, &t.TotalDistance,
			&count,

			&s.Id, &s.TrainingId, &s.SetOrder, &s.Repeat, &s.DistanceMeters,
			&s.Description, &s.StartType, &s.StartSeconds,
			&s.Equipment, &s.Group, &s.IsMain, &s.SetType,
			&s.Intensity, &s.Progression,

			&s.Style.Id, &s.Style.Name, &s.Style.ExeciseName,
			&s.Style.ExeciseDescription,
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, fmt.Errorf("TrainingSummaries scanning row: %w", err)
		}

		if lastTrainingId != t.Id {
			if s.Id != nil {
				t.Sets = append(t.Sets, s.intoTrainingSet())
			}
			ts = append(ts, t)
			lastTrainingId = t.Id
		} else if s.Id != nil {
			ts[len(ts)-1].Sets = append(ts[len(ts)-1].Sets, s.intoTrainingSet())
		}
	}

	return ts, count, nil
}

func (pool *PostgresDbPool) EditTrainingSession(
	ctx context.Context,
	id uuid.UUID,
	session struct {
		DurationMinutes *int
		Start           *time.Time
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
		DurationMinutes *int
		Start           *time.Time
	},
	tx pgx.Tx,
) (Training, error) {
	t := Training{}
	err := tx.QueryRow(ctx, updateTrainingSession, id, session.Start, session.DurationMinutes).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.CreatedAt, &t.ModifiedAt, &t.TotalDistance)

	if errors.Is(err, pgx.ErrNoRows) {
		return Training{}, fmt.Errorf("editTrainingSession not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return Training{}, fmt.Errorf("editTrainingSession update training query error: %w", err)
	}

	return t, nil
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
	err := Tx(
		ctx,
		pool,
		func(ctx context.Context, tx pgx.Tx) error { return moveSet(ctx, id, newSetOrder, tx) },
	)
	if err != nil {
		return fmt.Errorf("MoveSet: %w", err)
	}
	return nil
}

func moveSet(ctx context.Context, id uuid.UUID, newSetOrder int, tx pgx.Tx) error {
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

func scan(training *Training, rows pgx.Row) error {
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

	if c.Id == uuid.Nil {
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

func scanAll(rows pgx.Row) (Training, TrainingSet, SetComponent, error) {
	var row joinedRow
	dests := []any{
		&row.training.Id, &row.training.Start, &row.training.DurationMin,
		&row.training.CreatedAt, &row.training.ModifiedAt,

		&row.set.Id, &row.set.SetOrder, &row.set.Repeat, &row.set.DistanceMeters,
		&row.set.Description, &row.set.StartType, &row.set.StartSeconds,
		&row.set.Equipment, &row.set.Group, &row.set.IsMain, &row.set.SetType,
		&row.set.Intensity, &row.set.Progression,

		&row.setStyle.Id, &row.setStyle.Name, &row.setStyle.ExeciseName,
		&row.setStyle.ExeciseDescription,

		&row.comp.Id, &row.comp.Orders, &row.comp.IterationOrder, &row.comp.Repeat,
		&row.comp.DistanceMeters, &row.comp.StartType, &row.comp.StartSeconds,
		&row.comp.Intensity, &row.comp.Progression, &row.comp.Description,
		&row.comp.Equipment, &row.comp.Group,

		&row.compStyle.Id, &row.compStyle.Name, &row.compStyle.ExeciseName,
		&row.compStyle.ExeciseDescription,
	}

	err := rows.Scan(dests...)
	if err != nil {
		return Training{}, TrainingSet{}, SetComponent{}, fmt.Errorf("scanAll: %w", err)
	}

	return row.toTraining(), row.toTrainingSet(), row.toSetComponent(), nil
}

type joinedRow struct {
	training Training

	set      TrainingSet
	setStyle emptyStyle

	comp      emptySetComponent
	compStyle emptyStyle
}

func (r joinedRow) toTraining() Training {
	return Training{
		Id:          r.training.Id,
		Start:       r.training.Start,
		DurationMin: r.training.DurationMin,
		CreatedAt:   r.training.CreatedAt,
		ModifiedAt:  r.training.ModifiedAt,
	}
}

func (r joinedRow) toTrainingSet() TrainingSet {
	s := TrainingSet{
		Id:             r.set.Id,
		TrainingId:     r.training.Id,
		SetOrder:       r.set.SetOrder,
		Repeat:         r.set.Repeat,
		DistanceMeters: r.set.DistanceMeters,
		Description:    r.set.Description,
		Equipment:      r.set.Equipment,
		StartType:      r.set.StartType,
		StartSeconds:   r.set.StartSeconds,
		Group:          r.set.Group,
		IsMain:         r.set.IsMain,
		SetType:        r.set.SetType,
		Intensity:      r.set.Intensity,
		Progression:    r.set.Progression,
	}

	if r.setStyle.Id != nil {
		s.Style = &Style{
			Id:                 *r.setStyle.Id,
			Name:               *r.setStyle.Name,
			ExeciseName:        r.setStyle.ExeciseName,
			ExeciseDescription: r.setStyle.ExeciseDescription,
		}
		s.StyleId = &s.Style.Id
	}

	return s
}

func (r joinedRow) toSetComponent() SetComponent {
	if r.comp.Id == nil {
		return SetComponent{}
	}

	c := SetComponent{
		Id:             *r.comp.Id,
		SetId:          r.set.Id,
		Orders:         r.comp.Orders,
		IterationOrder: r.comp.IterationOrder,
		Repeat:         *r.comp.Repeat,
		DistanceMeters: *r.comp.DistanceMeters,
		Equipment:      r.comp.Equipment,
		StartType:      r.comp.StartType,
		StartSeconds:   r.comp.StartSeconds,
		Group:          r.comp.Group,
		Intensity:      r.comp.Intensity,
		Progression:    r.comp.Progression,
		Description:    r.comp.Description,
	}

	if r.compStyle.Id != nil {
		c.Style = &Style{
			Id:                 *r.compStyle.Id,
			Name:               *r.compStyle.Name,
			ExeciseName:        r.compStyle.ExeciseName,
			ExeciseDescription: r.compStyle.ExeciseDescription,
		}
		c.StyleId = &c.Style.Id
	}

	return c
}
