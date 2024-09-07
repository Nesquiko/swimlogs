package data

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const styleIdsExists = `
with ids as (select unnest($1::uuid[]) as id)
select ids.id as id, s.id is not null as exists
from ids left join styles s on ids.id = s.id;
`

type IdCheck struct {
	Id     uuid.UUID
	Exists bool
}

func (psg *PostgresDbPool) StyleIdsExist(ctx context.Context, ids []uuid.UUID) ([]IdCheck, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("StyleIdsExist empty ids")
	}

	rows, err := psg.pool.Query(ctx, styleIdsExists, ids)
	if err != nil {
		return nil, fmt.Errorf("StyleIdsExist query error: %w", err)
	}
	checks := make([]IdCheck, 0, len(ids))
	for rows.Next() {
		var check IdCheck
		err = rows.Scan(&check.Id, &check.Exists)
		if err != nil {
			return nil, fmt.Errorf("StyleIdsExist scan error: %w", err)
		}
		checks = append(checks, check)
	}
	return checks, nil
}

func (pool *PostgresDbPool) DeleteTraining(ctx context.Context, id uuid.UUID) error {
	err := Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		return pool.deleteTraining(ctx, id, tx)
	})
	if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}
	return nil
}

func (pool *PostgresDbPool) deleteTraining(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {
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

func (pool *PostgresDbPool) DeleteSet(ctx context.Context, id uuid.UUID) error {
	err := Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		return pool.deleteSet(ctx, id, tx)
	})
	if err != nil {
		return fmt.Errorf("DeleteSet: %w", err)
	}
	return nil
}

const reorderRemainingSets = `
with to_be_deleted as (select * from sets s where s.id = $1)
update sets
set set_order = sets.set_order - 1
from to_be_deleted
where sets.set_order > to_be_deleted.set_order and sets.training_id = to_be_deleted.training_id
`

func (pool *PostgresDbPool) deleteSet(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {
	trainingId, err := pool.trainingIdBySetId(ctx, id, tx)
	if err != nil {
		return fmt.Errorf("deleteSet retrieve training id query: %w", err)
	}

	ct, err := tx.Exec(ctx, reorderRemainingSets, id)
	if err != nil {
		return fmt.Errorf("deleteSet reorder query: %w", err)
	}

	ct, err = tx.Exec(ctx, "delete from sets where id = $1", id)
	if err != nil {
		return fmt.Errorf("deleteSet delete query: %w", err)
	} else if ct.RowsAffected() == 0 {
		return fmt.Errorf("deleteSet delete query set not found: %w", ErrRowsNotFound)
	}

	setsCount, err := pool.countTrainingSets(ctx, trainingId, tx)
	if err != nil {
		return fmt.Errorf("deleteSet set count: %w", err)
	}

	if setsCount != 0 {
		return nil
	}

	err = pool.deleteTraining(ctx, trainingId, tx)
	if err != nil {
		return fmt.Errorf("deleteSet delete training: %w", err)
	}

	return nil
}

const setCountInTraining = `
select count(*)
    from trainings t
    join sets s on t.id = s.training_id
where t.id = $1
`

func (pool *PostgresDbPool) countTrainingSets(
	ctx context.Context,
	trainingId uuid.UUID,
	tx pgx.Tx,
) (int, error) {
	count := -1
	err := tx.QueryRow(ctx, setCountInTraining, trainingId).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("countTrainingSets query: %w", err)
	}
	return count, nil
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

	trainingId, err := pool.trainingIdBySetId(ctx, id, tx)
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
		return pool.trainingIdBySetId(ctx, setId, tx)
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("TrainingIdBySetId: %w", err)
	}
	return id, nil
}

func (pool *PostgresDbPool) trainingIdBySetId(
	ctx context.Context,
	setId uuid.UUID,
	tx pgx.Tx,
) (uuid.UUID, error) {
	var trainingId uuid.UUID
	err := tx.QueryRow(ctx, "select training_id from sets where id = $1", setId).Scan(&trainingId)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, fmt.Errorf("trainingIdBySetId not found: %w", ErrRowsNotFound)
	} else if err != nil {
		return uuid.UUID{}, fmt.Errorf("trainingIdBySetId: %w", err)
	}

	return trainingId, nil
}
