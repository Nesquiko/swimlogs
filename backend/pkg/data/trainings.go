package data

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
    t.id, t.date, t.start_time, t.duration, t.total_distance,
    b.id, b.num, b.repeat, b.name, b.total_distance,
    s.id, s.num, s.repeat, s.distance, s.what, s.starting_rule, s.rule_seconds, s.total_distance
from trainings t
         join blocks b on t.id = b.training_id
         join sets s on b.id = s.block_id
where t.id = $1
order by b.num, s.num
`

func (db *PostgresDbConn) GetTrainingById(id uuid.UUID) (openapi.Training, error) {
	var t openapi.Training
	var date time.Time

	rows, err := db.Query(selectTraining, id)
	if err != nil {
		return openapi.Training{}, fmt.Errorf("GetTrainingById: %w", err)
	}
	defer rows.Close()

	count := 0
	blocks := make(map[uuid.UUID]*openapi.Block)
	for rows.Next() {
		count++
		var b openapi.Block
		var s openapi.TrainingSet

		err := rows.Scan(
			&t.Id,
			&date,
			&t.StartTime,
			&t.DurationMin,
			&t.TotalDistance,
			&b.Id,
			&b.Num,
			&b.Repeat,
			&b.Name,
			&b.TotalDistance,
			&s.Id,
			&s.Num,
			&s.Repeat,
			&s.Distance,
			&s.What,
			&s.StartingRule.Type,
			&s.StartingRule.Seconds,
			&s.TotalDistance,
		)
		if err != nil {
			return openapi.Training{}, fmt.Errorf("GetTrainingById: %w", err)
		}

		if _, ok := blocks[b.Id]; !ok {
			blocks[b.Id] = &b
		}
		blocks[b.Id].Sets = append(blocks[b.Id].Sets, s)
	}
	if count == 0 {
		return openapi.Training{}, fmt.Errorf("GetTrainingById: %w", ErrRowsNotFound)
	}
	t.Date = types.Date{Time: date}
	for _, b := range blocks {
		sort.Slice(b.Sets, func(i, j int) bool {
			return b.Sets[i].Num < b.Sets[j].Num
		})
		t.Blocks = append(t.Blocks, *b)
	}
	sort.Slice(t.Blocks, func(i, j int) bool {
		return t.Blocks[i].Num < t.Blocks[j].Num
	})

	return t, nil
}

var selectTrainingDetailsForThisWeek = `
select t.id, t.start, t.duration_min, t.total_distance, t.created_at, t.modified_at
from trainings t
where date_trunc('week', t.start) = date_trunc('week', current_date)
order by t.start, t.duration_min, t.total_distance
`

func (db *PostgresDbConn) GetTrainingDetailsInCurrentWeek() ([]Training, error) {
	var ts = make([]Training, 0)

	rows, err := db.Query(selectTrainingDetailsForThisWeek)
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
