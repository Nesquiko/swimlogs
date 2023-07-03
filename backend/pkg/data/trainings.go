package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	NoneStartType     = "None"
	IntervalStartType = "Interval"
	PauseStartType    = "Pause"
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
	ParentSetId    *uuid.UUID
	SetOrder       int
	SubSetOrder    *int
	TotalDistance  int
	Repeat         int
	DistanceMeters *int
	Description    *string
	StartType      string
	StartSeconds   *int
	SubSets        *[]TrainingSet
}

func (db *PostgresDbConn) SaveTraining(t Training) (Training, error) {
	return txWithResult(db.DB, func(tx *sql.Tx) (Training, error) {
		return db.saveTraining(t, tx)
	})
}

var insertTraining = `
insert into trainings (id, start, duration_min, total_distance, created_at, modified_at)
values ($1, $2, $3, $4, $5, $5)
returning id, start, duration_min, total_distance, created_at, modified_at
`

func (db *PostgresDbConn) saveTraining(t Training, tx *sql.Tx) (Training, error) {
	err := tx.QueryRow(insertTraining, t.Id, t.Start, t.DurationMin, t.TotalDistance, time.Now()).
		Scan(&t.Id, &t.Start, &t.DurationMin, &t.TotalDistance, &t.CreatedAt, &t.ModifiedAt)

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case ForeignKeyViolationCode:
			return Training{}, fmt.Errorf("saveTraining: %w", ErrForeignKeyViolation)
		case CheckViolationCode:
			return Training{}, fmt.Errorf("saveTraining: %w", ErrCheckViolation)
		case InvalidEnumTypeCode:
			return Training{}, fmt.Errorf("saveTraining: %w", ErrInvalidEnumType)
		default:
			return Training{}, fmt.Errorf("saveTraining: %w", err)
		}
	} else if err != nil {
		return Training{}, fmt.Errorf("saveTraining: %w", err)
	}

	for _, s := range t.Sets {
		_, err := db.saveSet(tx, s)
		if err != nil {
			return Training{}, fmt.Errorf("saveTraining: %w", err)
		}
	}

	return t, nil
}

var insertSet = `
insert into sets (id, parent_set_id, training_id, set_order, subset_order, repeat,
		distance_meters, description, start_type, start_seconds, total_distance)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
returning id, parent_set_id, training_id, set_order, subset_order, repeat, distance_meters, description, start_type,
    start_seconds, total_distance
`

func (db *PostgresDbConn) saveSet(tx *sql.Tx, s TrainingSet) (TrainingSet, error) {
	err := tx.QueryRow(
		insertSet,
		s.Id,
		s.ParentSetId,
		s.TrainingId,
		s.SetOrder,
		s.SubSetOrder,
		s.Repeat,
		s.DistanceMeters,
		s.Description,
		s.StartType,
		s.StartSeconds,
		s.TotalDistance,
	).Scan(
		&s.Id,
		&s.ParentSetId,
		&s.TrainingId,
		&s.SetOrder,
		&s.SubSetOrder,
		&s.Repeat,
		&s.DistanceMeters,
		&s.Description,
		&s.StartType,
		&s.StartSeconds,
		&s.TotalDistance,
	)

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case ForeignKeyViolationCode:
			return TrainingSet{}, fmt.Errorf("saveSet: %w", ErrForeignKeyViolation)
		case CheckViolationCode:
			return TrainingSet{}, fmt.Errorf("saveSet: %w", ErrCheckViolation)
		case InvalidEnumTypeCode:
			return TrainingSet{}, fmt.Errorf("saveSet: %w", ErrInvalidEnumType)
		default:
			return TrainingSet{}, fmt.Errorf("saveSet: %w", err)
		}
	} else if err != nil {
		return TrainingSet{}, fmt.Errorf("saveSet: %w", err)
	}

	return s, nil
}

var selectTraining = `
select
    t.id, t.start, t.duration_min, t.total_distance, t.created_at, t.modified_at,
    s.id, s.parent_set_id, s.training_id, s.set_order, s.subset_order, s.repeat,
    s.distance_meters, s.description, s.start_type, s.start_seconds, s.total_distance
from trainings t
         join sets s on t.id = s.training_id
where t.id = $1
order by s.set_order, s.subset_order nulls first
`

func (db *PostgresDbConn) GetTrainingById(id uuid.UUID) (Training, error) {
	var t Training
	rows, err := db.Query(selectTraining, id)
	if err != nil {
		return Training{}, fmt.Errorf("GetTrainingById: %w", err)
	}
	defer rows.Close()

	count := 0
	rootSets := make([]TrainingSet, 0)
	setsMap := make(map[uuid.UUID]*TrainingSet)
	for rows.Next() {
		count++
		var s TrainingSet
		err := rows.Scan(
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.TotalDistance,
			&t.CreatedAt,
			&t.ModifiedAt,
			&s.Id,
			&s.ParentSetId,
			&s.TrainingId,
			&s.SetOrder,
			&s.SubSetOrder,
			&s.Repeat,
			&s.DistanceMeters,
			&s.Description,
			&s.StartType,
			&s.StartSeconds,
			&s.TotalDistance,
		)
		if err != nil {
			return Training{}, fmt.Errorf("GetTrainingById: %w", err)
		}

		if s.ParentSetId == nil {
			rootSets = append(rootSets, s)
			setsMap[s.Id] = &rootSets[len(rootSets)-1]
			continue
		}

		parentSet, ok := setsMap[*s.ParentSetId]
		if !ok {
			log.Error().
				Str("training_id", t.Id.String()).
				Str("set_id", s.Id.String()).
				Str("parent_set_id", s.ParentSetId.String()).
				Msg("parent set not found")
			continue
		}

		if parentSet.SubSets == nil {
			parentSet.SubSets = &[]TrainingSet{}
		}

		newSubSets := append(*parentSet.SubSets, s)
		parentSet.SubSets = &newSubSets
		setsMap[s.Id] = &(*parentSet.SubSets)[len(*parentSet.SubSets)-1]
	}
	if count == 0 {
		return Training{}, fmt.Errorf("GetTrainingById: %w", ErrRowsNotFound)
	}
	t.Sets = rootSets

	return t, nil
}

var selectTrainingDetailsForThisWeek = `
select t.id, t.start, t.duration_min, t.total_distance, t.created_at, t.modified_at
from trainings t
where date_trunc('week', t.start) = date_trunc('week', $1::date)
order by t.start, t.duration_min, t.total_distance
`

func (db *PostgresDbConn) GetTrainingDetailsInWeek(week time.Time) ([]Training, error) {
	var ts = make([]Training, 0)

	rows, err := db.Query(selectTrainingDetailsForThisWeek, week)
	if err != nil {
		return nil, fmt.Errorf("GetTrainingDetailsForThisWeek: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Training
		err := rows.Scan(
			&t.Id,
			&t.Start,
			&t.DurationMin,
			&t.TotalDistance,
			&t.CreatedAt,
			&t.ModifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetTrainingDetailsForThisWeek: %w", err)
		}
		ts = append(ts, t)
	}

	return ts, nil
}
