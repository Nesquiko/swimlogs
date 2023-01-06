package app

import (
	"context"
	"database/sql"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
)

func (app *swimLogsApp) CreateTraining(
	request oapiGen.CreateTrainingRequestObject,
) (oapiGen.CreateTrainingResponseObject, error) {
	newTraining := request.Body
	if invalid := validateTraining(*newTraining); len(invalid) != 0 {
		return oapiGen.CreateTraining400JSONResponse{
			InvalidTrainingErrorResponseJSONResponse: invalidTrainingError(invalid),
		}, nil
	}

	t := transformRestTraining(*newTraining)
	err := app.db.InTx(func(tx *sql.Tx) error {
		var id *uuid.UUID
		var err error

		if newTraining.SessionId == nil {
			id, err = app.db.SaveTraining(t, tx)
		} else {
			id, err = app.db.SaveTrainingWithSesssionData(t, *newTraining.SessionId, tx)
		}

		if err != nil {
			return err
		}
		newTraining.Id = *id

		return nil
	})
	if err != nil {
		app.logger.Error(err)
		return oapiGen.CreateTraining500JSONResponse{
			InternalServerErrorResponseJSONResponse: internalServerError(),
		}, nil
	}

	return oapiGen.CreateTraining201JSONResponse(*newTraining), nil
}

func (app *swimLogsApp) GetTrainings(
	request oapiGen.GetTrainingsRequestObject,
) (oapiGen.GetTrainingsResponseObject, error) {
	app.logger.Info("GetTrainings endpoint called, but it shouldn't have been")
	return nil, nil
}

func (app *swimLogsApp) DeleteTraining(
	ctx context.Context,
	request oapiGen.DeleteTrainingRequestObject,
) (oapiGen.DeleteTrainingResponseObject, error) {
	return nil, nil
}

func (app *swimLogsApp) GetTrainingById(
	ctx context.Context,
	request oapiGen.GetTrainingByIdRequestObject,
) (oapiGen.GetTrainingByIdResponseObject, error) {
	return nil, nil
}

func (app *swimLogsApp) UpdateTraining(
	ctx context.Context,
	request oapiGen.UpdateTrainingRequestObject,
) (oapiGen.UpdateTrainingResponseObject, error) {
	return nil, nil
}

func transformRestTraining(t oapiGen.Training) data.Training {
	training := data.Training{
		Date:        t.Date.Time,
		Day:         (*string)(t.Day),
		DurationMin: t.DurationMin,
		StartTime:   t.StartTime,
	}

	training.Blocks = make([]data.Block, len(t.Blocks))
	for i, b := range t.Blocks {
		training.Blocks[i] = transformRestBlock(b)
	}

	return training
}

func transformRestBlock(b oapiGen.Block) data.Block {
	block := data.Block{
		Repeat: b.Repeat,
		Name:   b.Name,
	}

	block.Sets = make([]data.Set, len(b.Sets))
	for i, s := range b.Sets {
		block.Sets[i] = transformRestSet(s)
	}

	return block
}

func transformRestSet(s oapiGen.Set) data.Set {
	ruleSeconds := sql.NullInt16{Int16: 0, Valid: false}

	if s.StartingRule.Seconds != nil {
		ruleSeconds.Int16 = int16(*s.StartingRule.Seconds)
		ruleSeconds.Valid = true
	}

	return data.Set{
		Repeat:       s.Repeat,
		Distance:     s.Distance,
		What:         s.What,
		StartingRule: string(s.StartingRule.Rule),
		RuleSeconds:  ruleSeconds,
	}

}

func transformDataTraining(t data.Training) oapiGen.Training {
	training := oapiGen.Training{
		Id:          t.Id,
		Date:        types.Date{Time: t.Date},
		Day:         (*oapiGen.Day)(t.Day),
		DurationMin: t.DurationMin,
		StartTime:   t.StartTime,
	}

	training.Blocks = make([]oapiGen.Block, len(t.Blocks))
	for i, b := range t.Blocks {
		training.Blocks[i] = transformDataBlock(b)
	}

	return training
}

func transformDataBlock(b data.Block) oapiGen.Block {
	block := oapiGen.Block{
		Name:   b.Name,
		Repeat: b.Repeat,
	}

	block.Sets = make([]oapiGen.Set, len(b.Sets))
	for i, s := range b.Sets {
		block.Sets[i] = transformDataSet(s)
	}

	return block
}

func transformDataSet(s data.Set) oapiGen.Set {

	var seconds *int
	switch s.StartingRule {
	case string(oapiGen.Pause), string(oapiGen.Interval):
		s := int(s.RuleSeconds.Int16)
		seconds = &s
	default:
		seconds = nil
	}

	sr := oapiGen.StartingRule{
		Rule:    oapiGen.StartingRuleRule(s.StartingRule),
		Seconds: seconds,
	}

	return oapiGen.Set{
		Distance:     s.Distance,
		Repeat:       s.Repeat,
		StartingRule: sr,
		What:         s.What,
	}
}
