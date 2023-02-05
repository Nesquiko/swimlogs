package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var SelectTrainingDetailsCurrentWeek = "select t.id, t.date, t.day, t.starttime, t.duration, t.total_dist from training t where date_trunc('week', t.date ) = date_trunc('week', current_date) order by t.modified_at"

func (psql *postgresDbConn) GetDetailsOfTrainingsCurrentWeek() ([]Training, error) {
	rows, err := psql.Query(SelectTrainingDetailsCurrentWeek)
	if err != nil {
		return nil, fmt.Errorf("GetDetailsOfTrainingsCurrentWeek: %w", err)
	}
	defer rows.Close()

	trainings := make([]Training, 0)
	for rows.Next() {
		var t Training
		var startTime string

		err = rows.Scan(&t.Id, &t.Date, &t.Day, &startTime, &t.DurationMin, &t.TotalDistance)
		if err != nil {
			return nil, fmt.Errorf("GetDetailsOfTrainingsCurrentWeek: %w", err)
		}

		st, err := time.Parse(TimeLayout, startTime)
		if err != nil {
			return nil, fmt.Errorf("GetDetailsOfTrainingsCurrentWeek: %w", err)
		}
		startTime = st.Format("15:04")
		t.StartTime = &startTime

		trainings = append(trainings, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDetailsOfTrainingsCurrentWeek: %w", err)
	}

	return trainings, nil
}

var SelectTrainingDetails = "select t.id, t.date, t.day, t.starttime, t.duration, t.total_dist from training t order by t.modified_at limit $2 offset $1;"

func (psql *postgresDbConn) GetDetailsOfTrainings(page, pageSize int) ([]Training, error) {
	rows, err := psql.Query(SelectTrainingDetails, page*pageSize, pageSize)
	if err != nil {
		return nil, fmt.Errorf("GetDetailsOfTrainings: %w", err)
	}
	defer rows.Close()

	trainings := make([]Training, 0, pageSize)
	for rows.Next() {
		var t Training
		err = rows.Scan(&t.Id, &t.Date, &t.Day, &t.StartTime, &t.DurationMin, &t.TotalDistance)
		if err != nil {
			return nil, fmt.Errorf("GetDetailsOfTrainings: %w", err)
		}
		trainings = append(trainings, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetDetailsOfTrainings: %w", err)
	}

	return trainings, nil
}

var TrainingCount = "select count(id) from training"

func (psql *postgresDbConn) GetTrainingCount() (int, error) {
	var count int
	err := psql.QueryRow(TrainingCount).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetTrainingCount: %w", err)
	}
	return count, nil
}

var TrainingExists = "select count(1) from training where id = $1"

func (psql *postgresDbConn) TrainingExists(id uuid.UUID) (bool, error) {
	var exists int
	err := psql.QueryRow(TrainingExists, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

var UpdateTrainig = "update training set modified_at = now(), version = version + 1, date = $2, day = $3, starttime = $4, duration = $5, total_dist = $6 where id = $1 and version = $7"

func (psql *postgresDbConn) UpdateTrainingById(id uuid.UUID, t Training, tx *sql.Tx) error {
	res, err := tx.Exec(
		UpdateTrainig,
		id,
		t.Date,
		t.Day,
		t.StartTime,
		t.DurationMin,
		t.TotalDistance,
		t.Version,
	)
	if err != nil {
		return fmt.Errorf("UpdateTrainingById id='%s': %w", id.String(), err)
	}

	// Err can be ignored, because Postgres supports rows affected
	if res, _ := res.RowsAffected(); res == 0 {
		return fmt.Errorf("UpdateTrainingById id='%s': %w", id.String(), ErrRowNotFound)
	}

	for _, b := range t.Blocks {
		err := psql.updateBlock(b.Id, b, tx)
		if err != nil {
			return fmt.Errorf("UpdateTrainingById id='%s': %w", id.String(), err)
		}
	}

	return nil
}

var UpdateBlock = "update block set num = $2, repeat = $3, name = $4, total_dist = $5 where id = $1"

func (psql *postgresDbConn) updateBlock(id uuid.UUID, b Block, tx *sql.Tx) error {
	res, err := tx.Exec(UpdateBlock, id, b.Num, b.Repeat, b.Name, b.TotalDistance)
	if err != nil {
		return fmt.Errorf("updateBlock id='%s': %w", id.String(), err)
	}
	// Err can be ignored, because Postgres supports rows affected
	if res, _ := res.RowsAffected(); res == 0 {
		return fmt.Errorf("updateBlock id='%s': %w", id.String(), ErrRowNotFound)
	}

	for _, s := range b.Sets {
		err := psql.updateSet(s.Id, s, tx)
		if err != nil {
			return fmt.Errorf("updateBlock id='%s': %w", id.String(), err)
		}
	}

	return nil
}

var UpdateSet = "update set set num = $2, repeat = $3, distance = $4, what = $5, starting_rule = $6, rule_seconds = $7, total_dist = $8 where id = $1"

func (psql *postgresDbConn) updateSet(id uuid.UUID, s Set, tx *sql.Tx) error {
	res, err := tx.Exec(
		UpdateSet,
		id,
		s.Num,
		s.Repeat,
		s.Distance,
		s.What,
		s.StartingRule,
		s.RuleSeconds,
		s.TotalDistance,
	)
	if err != nil {
		return fmt.Errorf("updateSet id='%s': %w", id.String(), err)
	}
	// Err can be ignored, because Postgres supports rows affected
	if res, _ := res.RowsAffected(); res == 0 {
		return fmt.Errorf("updateSet id='%s': %w", id.String(), ErrRowNotFound)
	}

	return nil
}

var SelectTrainingById = "select t.id, t.version, t.date, t.day, t.starttime, t.duration, t.total_dist, b.id, b.num, b.repeat, b.name, b.total_dist, s.id, s.num, s.repeat , s.distance , s.what , s.starting_rule , s.rule_seconds , s.total_dist from training t left join block b on t.id = b.training_id left join set s on b.id = s.block_id where t.id = $1"

func (psql *postgresDbConn) GetTrainingById(id uuid.UUID) (Training, error) {
	var training Training

	rows, err := psql.Query(SelectTrainingById, id)
	if err != nil {
		return training, fmt.Errorf("GetTrainingById: %w", err)
	}
	defer rows.Close()

	cts := make([]completeTraining, 0)
	for rows.Next() {
		var ct completeTraining
		var startTime string
		err := rows.Scan(
			&ct.tId,
			&ct.tVersion,
			&ct.tDate,
			&ct.tDay,
			&startTime,
			&ct.tDur,
			&ct.tTotDist,
			&ct.bId,
			&ct.bNum,
			&ct.bRepeat,
			&ct.bName,
			&ct.bTotDist,
			&ct.sId,
			&ct.sNum,
			&ct.sRepeat,
			&ct.sDist,
			&ct.sWhat,
			&ct.sStartRule,
			&ct.sRuleSecs,
			&ct.sTotDist,
		)
		if err != nil {
			return training, fmt.Errorf("GetTrainingById: %w", err)
		}

		st, err := time.Parse(TimeLayout, startTime)
		if err != nil {
			ct.tStartT = startTime
		} else {
			ct.tStartT = st.Format("15:04")
		}
		cts = append(cts, ct)
	}
	if err := rows.Err(); err != nil {
		return training, fmt.Errorf("GetTrainingById: %w", err)
	}

	if len(cts) == 0 {
		return training, ErrRowNotFound
	}

	return createTraining(cts), nil
}

var DeleteTraining = "delete from training where id = $1"

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

var InsertTraining = "insert into training (id, created_at, modified_at, version, date, day, starttime, duration, total_dist) values ($1, $2, $3, 0, $4, $5, $6, $7, $8)"

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

var InsertTrainingFromSession = `insert into training
	(id, created_at, modified_at, version, date, total_dist, day, starttime, duration)
	select $1, $2, $3, 0, $4, $5, s.day ,s.starttime ,s.duration from session as s
	where s.id = $6
	returning id, day, starttime, duration`

func (psql *postgresDbConn) SaveTrainingWithSesssionData(
	t Training,
	sId uuid.UUID,
	tx *sql.Tx,
) (*Training, error) {
	base := createBase()
	t.Base = base

	var startTime string
	err := tx.QueryRow(
		InsertTrainingFromSession,
		t.Id,
		t.CreatedAt,
		t.ModifiedAt,
		t.Date,
		t.TotalDistance,
		sId,
	).Scan(&t.Id, &t.Day, &startTime, &t.DurationMin)
	if err != nil {
		return nil, fmt.Errorf("SaveTrainingWithSesssionData: %w", err)
	}
	st, err := time.Parse(TimeLayout, startTime)
	if err != nil {
		t.StartTime = &startTime
	} else {
		formatted := st.Format("15:04")
		t.StartTime = &formatted
	}

	for _, b := range t.Blocks {
		err = psql.saveBlock(b, t.Id, tx)
		if err != nil {
			return nil, fmt.Errorf("SaveTrainingWithSesssionData: %w", err)
		}
	}

	return &t, nil
}

var InsertBlock = "insert into block (id, num, repeat, name, total_dist, training_id) values ($1, $2, $3, $4, $5, $6)"

func (psql *postgresDbConn) saveBlock(b Block, trainingId uuid.UUID, tx *sql.Tx) error {
	b.Id = uuid.New()

	_, err := tx.Exec(InsertBlock, b.Id, b.Num, b.Repeat, b.Name, b.TotalDistance, trainingId)
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

var InsertSet = "insert into set (id, num, repeat, distance, what, starting_rule, rule_seconds, total_dist, block_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

func (psql *postgresDbConn) saveSet(s Set, blockId uuid.UUID, tx *sql.Tx) error {
	s.Id = uuid.New()

	_, err := tx.Exec(
		InsertSet,
		s.Id,
		s.Num,
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

func createTraining(cts []completeTraining) Training {
	var t Training
	blocks := make(map[uuid.UUID]*Block)

	t.Id = cts[0].tId
	t.Date = cts[0].tDate
	t.Day = &cts[0].tDay
	t.DurationMin = &cts[0].tDur
	t.StartTime = &cts[0].tStartT
	t.TotalDistance = cts[0].tTotDist
	t.Blocks = make([]Block, 0)

	for _, ct := range cts {
		if _, ok := blocks[ct.bId]; !ok {
			blocks[ct.bId] = &Block{
				Id:            ct.bId,
				Num:           ct.bNum,
				Repeat:        ct.bRepeat,
				Name:          ct.bName,
				TotalDistance: ct.bTotDist,
				Sets:          make([]Set, 0),
			}
		}

		b := blocks[ct.bId]
		s := Set{
			Id:            ct.sId,
			Num:           ct.sNum,
			Repeat:        ct.sRepeat,
			Distance:      ct.sDist,
			What:          ct.sWhat,
			StartingRule:  ct.sStartRule,
			RuleSeconds:   ct.sRuleSecs,
			TotalDistance: ct.sTotDist,
		}
		b.Sets = append(b.Sets, s)
	}

	for _, b := range blocks {
		t.Blocks = append(t.Blocks, *b)
	}

	return t
}

var SelectTrainings = "select t.*, b.*, s.* from training t left join block b on t.id = b.training_id left join set s on b.id = s.block_id group by t.id, b.id, s.id"

func (psql *postgresDbConn) GetTrainings(page, pageSize int) ([]Training, error) {
	return nil, errors.New("not implemented, not needed yet")
}
