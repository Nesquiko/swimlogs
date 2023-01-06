package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	InsertSet                 = "insert into set (id, repeat, distance, what, starting_rule, rule_seconds, total_dist, block_id) values ($1, $2, $3, $4, $5, $6, $7, $8)"
	InsertBlock               = "insert into block (id, repeat, name, total_dist, training_id) values ($1, $2, $3, $4, $5)"
	InsertTraining            = "insert into training (id, created_at, modified_at, date, day, starttime, duration, total_dist) values ($1, $2, $3, $4, $5, $6, $7, $8)"
	InsertTrainingFromSession = `insert into training (id, created_at, modified_at, date, total_dist, day, starttime, duration)
	select $1, $2, $3, $4, $5, s.day ,s.starttime ,s.duration from session as s where s.id = $6`

	SelectTrainings = "select t.*, b.*, s.* from training t left join block b on t.id = b.training_id left join set s on b.id = s.block_id group by t.id, b.id, s.id"

	DeleteTraining = "delete from training where id = $1"
)

func (psql *postgresDbConn) DeleteTraining(id uuid.UUID, tx *sql.Tx) error {
	res, err := tx.Exec(DeleteTraining, id)
	if err != nil {
		return fmt.Errorf("DeleteTraining: %w", err)
	}

	// Err can be ignored, because Postgres supports rows affected
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrRowNotFound
	}
	return nil
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
		t.TotalDistance,
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

	_, err := tx.Exec(
		InsertTrainingFromSession,
		t.Id,
		t.CreatedAt,
		t.ModifiedAt,
		t.Date,
		t.TotalDistance,
		sId,
	)
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

	_, err := tx.Exec(InsertBlock, b.Id, b.Repeat, b.Name, b.TotalDistance, trainingId)
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
		s.TotalDistance,
		blockId,
	)
	if err != nil {
		return fmt.Errorf("saveSet: %w", err)
	}

	return nil
}

func (psql *postgresDbConn) GetTrainings(page, pageSize int) ([]Training, error) {
	return nil, errors.New("not implemented, not needed yet")
}
