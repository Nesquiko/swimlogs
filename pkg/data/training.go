package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	InsertSet                 = "insert into set (id, repeat, distance, what, starting_rule, rule_seconds, block_id) values ($1, $2, $3, $4, $5, $6, $7)"
	InsertBlock               = "insert into block (id, repeat, name, training_id) values ($1, $2, $3, $4)"
	InsertTraining            = "insert into training (id, created_at, modified_at, date, day, starttime, duration) values ($1, $2, $3, $4, $5, $6, $7)"
	InsertTrainingFromSession = `insert into training (id, created_at, modified_at, date, day, starttime, duration)
	select $1, $2, $3, $4, s.day ,s.starttime ,s.duration from session as s where s.id = $5`

	SelectTrainings = "select t.*, b.*, s.* from training t left join block b on t.id = b.training_id left join set s on b.id = s.block_id group by t.id, b.id, s.id"
)

func (psql *postgresDbConn) GetTrainings(page, pageSize int) ([]Training, error) {
	return nil, errors.New("not implemented, not needed yet")
}

func (psql *postgresDbConn) SaveTraining(t Training, tx *sql.Tx) (*uuid.UUID, error) {
	base := createBase()
	t.Base = base

	_, err := tx.Exec(
		InsertTraining,
		t.Id,
		t.CreatedAt,
		t.ModifiedAt,
		t.Date,
		t.Day,
		t.StartTime,
		t.DurationMin,
	)
	if err != nil {
		return nil, fmt.Errorf("SaveTraining: %w", err)
	}

	for _, b := range t.Blocks {
		err = psql.saveBlock(b, t.Id, tx)
		if err != nil {
			return nil, fmt.Errorf("SaveTraining: %w", err)
		}
	}

	return &t.Id, nil
}

func (psql *postgresDbConn) SaveTrainingWithSesssionData(
	t Training,
	sId uuid.UUID,
	tx *sql.Tx,
) (*uuid.UUID, error) {
	base := createBase()
	t.Base = base

	_, err := tx.Exec(InsertTrainingFromSession, t.Id, t.CreatedAt, t.ModifiedAt, t.Date, sId)
	if err != nil {
		return nil, fmt.Errorf("SaveTrainingWithSesssionData: %w", err)
	}

	for _, b := range t.Blocks {
		err = psql.saveBlock(b, t.Id, tx)
		if err != nil {
			return nil, fmt.Errorf("SaveTrainingWithSesssionData: %w", err)
		}
	}

	return &t.Id, nil
}

func (psql *postgresDbConn) saveBlock(b Block, trainingId uuid.UUID, tx *sql.Tx) error {
	b.Id = uuid.New()

	_, err := tx.Exec(InsertBlock, b.Id, b.Repeat, b.Name, trainingId)
	if err != nil {
		return fmt.Errorf("saveBlock: %w", err)
	}

	for _, s := range b.Sets {
		err = psql.saveSet(s, b.Id, tx)
		if err != nil {
			return fmt.Errorf("saveBlock: %w", err)
		}
	}

	return nil
}

func (psql *postgresDbConn) saveSet(s Set, blockId uuid.UUID, tx *sql.Tx) error {
	s.Id = uuid.New()

	_, err := tx.Exec(
		InsertSet,
		s.Id,
		s.Repeat,
		s.Distance,
		s.What,
		s.StartingRule,
		s.RuleSeconds,
		blockId,
	)
	if err != nil {
		return fmt.Errorf("saveSet: %w", err)
	}

	return nil
}
