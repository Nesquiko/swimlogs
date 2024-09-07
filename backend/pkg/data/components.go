package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

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
