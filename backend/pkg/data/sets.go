package data

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (pool *PostgresDbPool) DeleteSet(ctx context.Context, id uuid.UUID) error {
	err := Tx(ctx, pool, func(ctx context.Context, tx pgx.Tx) error {
		return deleteSet(ctx, id, tx)
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

func deleteSet(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {
	trainingId, err := trainingIdBySetId(ctx, id, tx)
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

	setsCount, err := countTrainingSets(ctx, trainingId, tx)
	if err != nil {
		return fmt.Errorf("deleteSet set count: %w", err)
	}

	if setsCount != 0 {
		return nil
	}

	err = deleteTraining(ctx, trainingId, tx)
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

func countTrainingSets(ctx context.Context, trainingId uuid.UUID, tx pgx.Tx) (int, error) {
	count := -1
	err := tx.QueryRow(ctx, setCountInTraining, trainingId).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("countTrainingSets query: %w", err)
	}
	return count, nil
}

const insertSet = `
insert into sets (id, training_id, set_order, repeat, distance_meters, description,
    start_type, start_seconds, equipment, "group", is_main, type, style_id,
    intensity, progression)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
`

func persistSet(ctx context.Context, tx pgx.Tx, s TrainingSet) error {
	args := []any{
		s.Id, s.TrainingId, s.SetOrder, s.Repeat, s.DistanceMeters,
		s.Description, s.StartType, s.StartSeconds, s.Equipment, s.Group,
		s.IsMain, s.SetType, s.StyleId, s.Intensity, s.Progression,
	}

	_, err := tx.Exec(ctx, insertSet, args...)
	if err != nil {
		return fmt.Errorf("persistSet: %w", err)
	}

	if s.Components == nil {
		return nil
	}

	for i, c := range *s.Components {
		err := persistComponent(ctx, tx, c)
		if err != nil {
			return fmt.Errorf("persistSet component %d: %w", i, err)
		}
	}

	return nil
}
