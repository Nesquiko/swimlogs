package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Training struct {
	Id          uuid.UUID
	Start       time.Time
	DurationMin int
	Sets        []TrainingSet

	TotalDistance int
	CreatedAt     time.Time
	ModifiedAt    time.Time
}

type TrainingSet struct {
	Id             uuid.UUID
	TrainingId     uuid.UUID
	SetOrder       int
	Repeat         int
	DistanceMeters int
	Description    *string
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	Group          *string
}

func (pool *PostgresDbPool) PersistTraining(ctx context.Context, t Training) (Training, error) {
	return TxWithResult(ctx, pool, func(ctx context.Context, tx pgx.Tx) (Training, error) {
		t, err := pool.persistTraining(ctx, t, tx)
		if err != nil {
			return t, fmt.Errorf("PersistTraining tx: %w", err)
		}
		return t, nil
	})
}

func (pool *PostgresDbPool) DeleteTraining(ctx context.Context, id uuid.UUID) error {
	return Tx(pool, func(tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, "delete from trainings where id = $1", id)
		if err != nil {
			return fmt.Errorf("DeleteTraining: %w", err)
		} else if ct.RowsAffected() == 0 {
			return fmt.Errorf("DeleteTraining training doesnt exist: %w", ErrRowsNotFound)
		}
		return nil
	})
}

var selectTrainingSummariesPage = `
select t.id, t.start, t.duration_min, t.created_at, t.modified_at,
    sum(s.repeat * s.distance_meters), count(t.id) over ()
from trainings t
    join sets s on t.id = s.training_id
group by t.id, t.start, t.duration_min, t.created_at, t.modified_at
order by t.start desc, t.duration_min, t.created_at
limit $1 offset $2
`

func (pool *PostgresDbPool) TrainingSummaries(
	ctx context.Context,
	page, pageSize int,
) ([]Training, int, error) {
	tds := make([]Training, 0)

	rows, err := pool.Query(
		ctx,
		selectTrainingSummariesPage,
		pageSize,
		page*pageSize,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("TrainingSummaries query error: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var t Training
		err := rows.Scan(
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.CreatedAt,
			&t.ModifiedAt,
			&t.TotalDistance,
			&count,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("TrainingSummaries scanning row: %w", err)
		}
		tds = append(tds, t)
	}

	return tds, count, nil
}

var selectTrainingSummariesInDateRange = `
select t.id, t.start, t.duration_min, t.created_at, t.modified_at, sum(s.repeat * s.distance_meters)
from trainings t
    join sets s on t.id = s.training_id
where date(t.start) between $1::date and $2::date
group by t.id, t.start, t.duration_min, t.created_at, t.modified_at
order by t.start, t.duration_min, t.created_at
`

func (pool *PostgresDbPool) TrainingSummariesInRange(
	ctx context.Context,
	start, end time.Time,
) ([]Training, error) {
	tds := make([]Training, 0)

	rows, err := pool.Query(ctx, selectTrainingSummariesInDateRange, start, end)
	if err != nil {
		return nil, fmt.Errorf(
			"TrainingSummariesInRange from %s to %s query error: %w",
			start,
			end,
			err,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var t Training
		err := rows.Scan(
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.CreatedAt,
			&t.ModifiedAt,
			&t.TotalDistance,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"TrainingSummariesInRange from %s to %s scanning error: %w",
				start,
				end,
				err,
			)
		}
		tds = append(tds, t)
	}

	return tds, nil
}

var selectTraining = `
select
    t.id, t.start, t.duration_min, t.created_at, t.modified_at,
    s.id, s.training_id, s.set_order, s.repeat, s.distance_meters, s.description,
    s.start_type, s.start_seconds, s.equipment, s.group
from trainings t join sets s on t.id = s.training_id
where t.id = $1
order by s.set_order
`

func (pool *PostgresDbPool) Training(ctx context.Context, id uuid.UUID) (Training, error) {
	t := Training{}
	rows, err := pool.Query(ctx, selectTraining, id)
	if err != nil {
		return Training{}, fmt.Errorf("Training query error: %w", err)
	}

	for rows.Next() {
		s := TrainingSet{}
		err := rows.Scan(
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.CreatedAt,
			&t.ModifiedAt,
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
		)
		if err != nil {
			return Training{}, fmt.Errorf("Training scanning error: %w", err)
		}
		t.Sets = append(t.Sets, s)
	}
	rows.Close()

	if rows.CommandTag().RowsAffected() == 0 {
		return Training{}, fmt.Errorf("Training id doesnt exist: %w", ErrRowsNotFound)
	}

	return t, nil
}

func (pool *PostgresDbPool) EditTraining(
	ctx context.Context,
	id uuid.UUID,
	t Training,
) (Training, error) {
	return TxWithResult(ctx, pool, func(ctx context.Context, tx pgx.Tx) (Training, error) {
		t, err := pool.editTraining(ctx, id, t, tx)
		if err != nil {
			return t, fmt.Errorf("EditTraining tx: %w", err)
		}
		return t, nil
	})
}

var insertTraining = `
insert into trainings (id, start, duration_min, created_at, modified_at)
values ($1, $2, $3, now(), now())
returning id, start, duration_min, created_at, modified_at
`

func (pool *PostgresDbPool) persistTraining(
	ctx context.Context,
	t Training,
	tx pgx.Tx,
) (Training, error) {
	err := tx.QueryRow(ctx, insertTraining, t.Id, t.Start, t.DurationMin).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.CreatedAt, &t.ModifiedAt)
	if err != nil {
		return Training{}, fmt.Errorf("persistTraining persisting training: %w", err)
	}

	for i, s := range t.Sets {
		ts, err := pool.persistSet(ctx, tx, s)
		if err != nil {
			return Training{}, fmt.Errorf("persistTraining set %d: %w", i, err)
		}
		t.Sets[i] = ts
	}

	return t, nil
}

var insertSet = `
insert into sets (id, training_id, set_order, repeat, distance_meters,
    description, start_type, start_seconds, equipment, "group")
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning id, training_id, set_order, repeat, distance_meters,
    description, start_type, start_seconds, equipment, "group"
`

func (pool *PostgresDbPool) persistSet(
	ctx context.Context,
	tx pgx.Tx,
	s TrainingSet,
) (TrainingSet, error) {
	err := tx.QueryRow(
		ctx,
		insertSet,
		s.Id,
		s.TrainingId,
		s.SetOrder,
		s.Repeat,
		s.DistanceMeters,
		s.Description,
		s.StartType,
		s.StartSeconds,
		s.Equipment,
		s.Group,
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
	)
	if err != nil {
		return TrainingSet{}, fmt.Errorf("persistSet: %w", err)
	}

	return s, nil
}

var updateTraining = `
update trainings
set start          = $2,
    duration_min   = $3,
    total_distance = $4,
    modified_at    = now()
where id = $1
returning id, start, duration_min, total_distance, created_at, modified_at
`

func (pool *PostgresDbPool) editTraining(
	ctx context.Context,
	id uuid.UUID,
	t Training,
	tx pgx.Tx,
) (Training, error) {
	err := tx.QueryRow(ctx, updateTraining, id, t.Start, t.DurationMin).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.CreatedAt, &t.ModifiedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return Training{}, fmt.Errorf("editTraining not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return Training{}, fmt.Errorf("editTraining update training query error: %w", err)
	}

	for i, s := range t.Sets {
		ts, err := pool.editSet(ctx, tx, s)
		if err != nil {
			return Training{}, fmt.Errorf("editTraining set %d: %w", i, err)
		}
		t.Sets[i] = ts
	}

	return t, nil
}

var updateSet = `
update sets
set set_order       = $2,
    repeat          = $3,
    distance_meters = $4,
    description     = $5,
    start_type      = $6,
    start_seconds   = $7,
    total_distance  = $8,
    equipment       = $9,
    "group"             = $10
where id = $1
returning id, training_id, set_order, repeat, distance_meters, description,
    start_type, start_seconds, total_distance, equipment, "group"
`

func (pool *PostgresDbPool) editSet(
	ctx context.Context,
	tx pgx.Tx,
	s TrainingSet,
) (TrainingSet, error) {
	setExists, err := pool.setExists(ctx, tx, s)
	if err != nil {
		return TrainingSet{}, fmt.Errorf("editSet exists query: %w", err)
	}

	if !setExists {
		s.Id = uuid.New()
		return pool.persistSet(ctx, tx, s)
	}

	err = tx.QueryRow(
		ctx,
		updateSet,
		s.Id,
		s.SetOrder,
		s.Repeat,
		s.DistanceMeters,
		s.Description,
		s.StartType,
		s.StartSeconds,
		s.Equipment,
		s.Group,
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
	)
	if err != nil {
		return TrainingSet{}, fmt.Errorf("editSet query error: %w, id: %s", err, s.Id)
	}

	return s, nil
}

var setExists = "select exists(select 1 from sets where id = $1)"

func (pool *PostgresDbPool) setExists(ctx context.Context, tx pgx.Tx, s TrainingSet) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, setExists, s.Id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("isSetNew query error: %w, id: %s", err, s.Id)
	}
	return exists, nil
}
