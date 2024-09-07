package data

import (
	"context"
	"fmt"

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

const insertComponent = `
insert into set_components (id, set_id, component_orders, iteration_order, repeat, distance_meters, start_type,
                            start_seconds, style_id, intensity, progression, description, equipment, "group")
values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
`

func persistComponent(ctx context.Context, tx pgx.Tx, c SetComponent) error {
	if !c.Id.Valid {
		return fmt.Errorf("persistComponent component without id")
	}

	args := []any{
		c.Id.UUID, c.SetId, c.Orders, c.IterationOrder, c.Repeat, c.DistanceMeters,
		c.StartType, c.StartSeconds, c.StyleId, c.Intensity, c.Progression,
		c.Description, c.Equipment, c.Group,
	}

	_, err := tx.Exec(ctx, insertComponent, args...)
	if err != nil {
		return fmt.Errorf("persistComponent: %w", err)
	}

	return nil
}
