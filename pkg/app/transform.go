package app

import (
	"database/sql"
	"strings"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/deepmap/oapi-codegen/pkg/types"
)

func transformRestSession(session oapiGen.Session) data.Session {
	return data.Session{
		Day:         strings.ToLower(string(session.Day)),
		StartTime:   session.StartTime,
		DurationMin: session.DurationMin,
	}
}

func transformDataSession(session data.Session) oapiGen.Session {
	return oapiGen.Session{
		Id:          session.Id,
		Day:         oapiGen.Day(session.Day),
		StartTime:   session.StartTime,
		DurationMin: session.DurationMin,
	}
}

func transformRestTraining(t *oapiGen.Training) data.Training {
	training := data.Training{
		Date:        t.Date.Time,
		Day:         (*string)(t.Day),
		DurationMin: t.DurationMin,
		StartTime:   t.StartTime,
	}

	training.Blocks = make([]data.Block, len(t.Blocks))
	totDist := 0
	for i, b := range t.Blocks {
		block := transformRestBlock(b)
		training.Blocks[i] = block
		totDist += block.TotalDistance
	}
	training.TotalDistance = totDist

	return training
}

func transformRestBlock(b oapiGen.Block) data.Block {
	block := data.Block{
		Id:     b.Id,
		Num:    b.Num,
		Repeat: b.Repeat,
		Name:   b.Name,
	}

	block.Sets = make([]data.Set, len(b.Sets))
	totDist := 0
	for i, s := range b.Sets {
		set := transformRestSet(s)
		block.Sets[i] = set
		totDist += set.TotalDistance
	}
	block.TotalDistance = totDist

	return block
}

func transformRestSet(s oapiGen.Set) data.Set {
	ruleSeconds := sql.NullInt16{Int16: 0, Valid: false}

	if s.StartingRule.Seconds != nil {
		ruleSeconds.Int16 = int16(*s.StartingRule.Seconds)
		ruleSeconds.Valid = true
	}
	totDist := s.Distance * s.Repeat

	return data.Set{
		Id:            s.Id,
		Num:           s.Num,
		Repeat:        s.Repeat,
		Distance:      s.Distance,
		What:          s.What,
		StartingRule:  string(s.StartingRule.Rule),
		RuleSeconds:   ruleSeconds,
		TotalDistance: totDist,
	}
}

func transformDataTraining(t data.Training) oapiGen.Training {
	training := oapiGen.Training{
		Id:          t.Id,
		Date:        types.Date{Time: t.Date},
		Day:         (*oapiGen.Day)(t.Day),
		DurationMin: t.DurationMin,
		StartTime:   t.StartTime,
		TotalDist:   &t.TotalDistance,
	}

	training.Blocks = make([]oapiGen.Block, len(t.Blocks))
	for i, b := range t.Blocks {
		training.Blocks[i] = transformDataBlock(b)
	}

	return training
}

func transformDataBlock(b data.Block) oapiGen.Block {
	block := oapiGen.Block{
		Id:        b.Id,
		Num:       b.Num,
		Name:      b.Name,
		Repeat:    b.Repeat,
		TotalDist: &b.TotalDistance,
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
		Id:           s.Id,
		Num:          s.Num,
		Distance:     s.Distance,
		Repeat:       s.Repeat,
		StartingRule: sr,
		What:         s.What,
		TotalDist:    &s.TotalDistance,
	}
}

func transormToDetails(ts []data.Training) []oapiGen.TrainingDetail {
	details := make([]oapiGen.TrainingDetail, len(ts))

	for i, t := range ts {
		details[i] = oapiGen.TrainingDetail{
			Id:          t.Id,
			Date:        types.Date{t.Date},
			Day:         oapiGen.Day(*t.Day),
			StartTime:   *t.StartTime,
			DurationMin: *t.DurationMin,
			TotalDist:   t.TotalDistance,
		}
	}

	return details
}
