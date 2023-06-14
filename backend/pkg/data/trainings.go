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

func (db *PostgresDbConn) SaveTraining(newT openapi.NewTraining) (openapi.TrainingDetail, error) {
	setTotalDistances(&newT)

	t, err := txWithResult(db.DB, func(tx *sql.Tx) (openapi.Training, error) {
		return db.saveTraining(newT, tx)
	})

	if err != nil {
		return openapi.TrainingDetail{}, fmt.Errorf("SaveTraining: %w", err)
	}

	td := openapi.TrainingDetail{
		Id:            t.Id,
		Date:          t.Date,
		StartTime:     t.StartTime,
		DurationMin:   t.DurationMin,
		TotalDistance: t.TotalDistance,
	}

	return td, nil
}

var insertTraining = `
insert into trainings (date, start_time, duration, total_distance, created_at, modified_at)
values ($1, $2, $3, $4, now(), now())
returning id, date, start_time, duration, total_distance
`

func (db *PostgresDbConn) saveTraining(
	newT openapi.NewTraining,
	tx *sql.Tx,
) (openapi.Training, error) {
	var t openapi.Training
	var date time.Time
	err := tx.QueryRow(
		insertTraining,
		newT.Date.Time,
		newT.StartTime,
		newT.DurationMin,
		newT.TotalDistance,
	).Scan(
		&t.Id,
		&date,
		&t.StartTime,
		&t.DurationMin,
		&t.TotalDistance,
	)
	t.Date = types.Date{Time: date}

	if psErr, ok := err.(*pq.Error); ok {
		switch psErr.Code {
		case ForeignKeyViolationCode:
			return openapi.Training{}, fmt.Errorf("saveTraining: %w", ErrForeignKeyViolation)
		case CheckViolationCode:
			return openapi.Training{}, fmt.Errorf("saveTraining: %w", ErrCheckViolation)
		case InvalidEnumTypeCode:
			return openapi.Training{}, fmt.Errorf("saveTraining: %w", ErrInvalidEnumType)
		default:
			return openapi.Training{}, fmt.Errorf("saveTraining: %w", err)
		}
	} else if err != nil {
		return openapi.Training{}, fmt.Errorf("saveTraining: %w", err)
	}

	for _, b := range newT.Blocks {
		_, err := db.saveBlock(b, t.Id, tx)
		if err != nil {
			return openapi.Training{}, fmt.Errorf("saveTraining: %w", err)
		}
	}

	return t, nil
}

var insertBlock = `
insert into blocks (num, repeat, name, total_distance, training_id)
values ($1, $2, $3, $4, $5)
returning id, num, repeat, name, total_distance
`

func (db *PostgresDbConn) saveBlock(
	b openapi.NewBlock,
	trainingId uuid.UUID,
	tx *sql.Tx,
) (openapi.Block, error) {
	var block openapi.Block
	err := tx.QueryRow(
		insertBlock,
		b.Num,
		b.Repeat,
		b.Name,
		b.TotalDistance,
		trainingId,
	).Scan(
		&block.Id,
		&block.Num,
		&block.Repeat,
		&block.Name,
		&block.TotalDistance,
	)
	if psErr, ok := err.(*pq.Error); ok {
		switch psErr.Code {
		case ForeignKeyViolationCode:
			return openapi.Block{}, fmt.Errorf("saveBlock: %w", ErrForeignKeyViolation)
		case CheckViolationCode:
			return openapi.Block{}, fmt.Errorf("saveBlock: %w", ErrCheckViolation)
		default:
			return openapi.Block{}, fmt.Errorf("saveBlock: %w", err)
		}
	} else if err != nil {
		return openapi.Block{}, fmt.Errorf("saveBlock: %w", err)
	}

	for _, s := range b.Sets {
		_, err := db.saveSet(s, block.Id, tx)
		if err != nil {
			return openapi.Block{}, fmt.Errorf("saveBlock: %w", err)
		}
	}

	return block, nil
}

var insertSet = `
insert into sets (num, repeat, distance, what, starting_rule, rule_seconds, total_distance, block_id)
values ($1, $2, $3, $4, $5::starting_rule, $6, $7, $8)
returning id, num, repeat, distance, what, starting_rule, rule_seconds, total_distance
`

func (db *PostgresDbConn) saveSet(
	s openapi.NewTrainingSet,
	blockId uuid.UUID,
	tx *sql.Tx,
) (openapi.TrainingSet, error) {
	set := openapi.TrainingSet{StartingRule: openapi.StartingRule{}}

	err := tx.QueryRow(
		insertSet,
		s.Num,
		s.Repeat,
		s.Distance,
		s.What,
		s.StartingRule.Type,
		s.StartingRule.Seconds,
		s.Distance*s.Repeat,
		blockId,
	).Scan(
		&set.Id,
		&set.Num,
		&set.Repeat,
		&set.Distance,
		&set.What,
		&set.StartingRule.Type,
		&set.StartingRule.Seconds,
		&set.TotalDistance,
	)

	if psErr, ok := err.(*pq.Error); ok {
		switch psErr.Code {
		case ForeignKeyViolationCode:
			return openapi.TrainingSet{}, fmt.Errorf("saveSet: %w", ErrForeignKeyViolation)
		case CheckViolationCode:
			return openapi.TrainingSet{}, fmt.Errorf("saveSet: %w", ErrCheckViolation)
		case InvalidEnumTypeCode:
			return openapi.TrainingSet{}, fmt.Errorf("saveSet: %w", ErrInvalidEnumType)
		default:
			return openapi.TrainingSet{}, fmt.Errorf("saveSet: %w", err)
		}
	} else if err != nil {
		return openapi.TrainingSet{}, fmt.Errorf("saveSet: %w", err)
	}

	return set, nil
}

func setTotalDistances(t *openapi.NewTraining) {
	t.TotalDistance = 0
	for i := range t.Blocks {
		b := &t.Blocks[i]
		b.TotalDistance = 0

		for i := range b.Sets {
			s := &b.Sets[i]
			sDist := s.Distance * s.Repeat
			b.TotalDistance += sDist
		}

		b.TotalDistance *= b.Repeat
		t.TotalDistance += b.TotalDistance
	}
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
select
    t.id,
    t.date,
    t.start_time,
    t.duration,
    t.total_distance
from trainings t
where date_trunc('week', t.date) = date_trunc('week', current_date)
order by t.date, t.start_time
`

func (db *PostgresDbConn) GetTrainingDetailsForThisWeek() ([]openapi.TrainingDetail, error) {
	var tds = make([]openapi.TrainingDetail, 0)

	rows, err := db.Query(selectTrainingDetailsForThisWeek)
	if err != nil {
		return nil, fmt.Errorf("GetTrainingDetailsForThisWeek: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var td openapi.TrainingDetail
		var date time.Time

		err := rows.Scan(
			&td.Id,
			&date,
			&td.StartTime,
			&td.DurationMin,
			&td.TotalDistance,
		)
		if err != nil {
			return nil, fmt.Errorf("GetTrainingDetailsForThisWeek: %w", err)
		}

		td.Date = types.Date{Time: date}
		tds = append(tds, td)
	}

	return tds, nil
}
