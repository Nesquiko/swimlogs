package data

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Training struct {
	Id            uuid.UUID
	Start         time.Time
	DurationMin   int
	TotalDistance int
	Sets          []TrainingSet

	CreatedAt  time.Time
	ModifiedAt time.Time
}

type TrainingSet struct {
	Id             uuid.UUID
	TrainingId     uuid.UUID
	SetOrder       int
	TotalDistance  int
	Repeat         int
	DistanceMeters int
	StartType      string
	Description    *string
	StartSeconds   *int
	Equipment      *[]string
}

func (pool *PostgresDbPool) PersistTraining(t Training) (Training, error) {
	return TxWithResult(pool, func(tx pgx.Tx) (Training, error) {
		t, err := pool.persistTraining(t, tx)
		if err != nil {
			return t, fmt.Errorf("PersistTraining tx: %w", err)
		}
		return t, nil
	})
}

var insertTraining = `
insert into trainings (id, start, duration_min, total_distance, created_at, modified_at)
values ($1, $2, $3, $4, now(), now())
returning id, start, duration_min, total_distance, created_at, modified_at
`

func (pool *PostgresDbPool) persistTraining(t Training, tx pgx.Tx) (Training, error) {
	err := tx.QueryRow(context.Background(), insertTraining, t.Id, t.Start, t.DurationMin, t.TotalDistance).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.TotalDistance, &t.CreatedAt, &t.ModifiedAt)
	if err != nil {
		return Training{}, fmt.Errorf("persistTraining persisting training: %w", err)
	}

	for i, s := range t.Sets {
		ts, err := pool.persistSet(tx, s)
		if err != nil {
			return Training{}, fmt.Errorf("persistTraining set %d: %w", i, err)
		}
		t.Sets[i] = ts
	}

	return t, nil
}

var insertSet = `
insert into sets (id, training_id, set_order, repeat, distance_meters,
    description, start_type, start_seconds, total_distance, equipment)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning id, training_id, set_order, repeat, distance_meters,
    description, start_type, start_seconds, total_distance, equipment
`

func (pool *PostgresDbPool) persistSet(tx pgx.Tx, s TrainingSet) (TrainingSet, error) {
	err := tx.QueryRow(
		context.Background(),
		insertSet,
		s.Id,
		s.TrainingId,
		s.SetOrder,
		s.Repeat,
		s.DistanceMeters,
		s.Description,
		s.StartType,
		s.StartSeconds,
		s.TotalDistance,
		s.Equipment,
	).Scan(
		&s.Id,
		&s.TrainingId,
		&s.SetOrder,
		&s.Repeat,
		&s.DistanceMeters,
		&s.Description,
		&s.StartType,
		&s.StartSeconds,
		&s.TotalDistance,
		s.Equipment,
	)
	if err != nil {
		return TrainingSet{}, fmt.Errorf("persistSet: %w", err)
	}

	return s, nil
}
