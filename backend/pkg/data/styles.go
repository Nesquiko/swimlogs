package data

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
		return []IdCheck{}, nil
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
