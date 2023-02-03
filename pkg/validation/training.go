package validation

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

var startingRules = map[string]bool{
	string(oapiGen.None):     true,
	string(oapiGen.Pause):    true,
	string(oapiGen.Interval): true,
}

func ValidateTraining(t oapiGen.Training) *oapiGen.InvalidTraining {
	invalid := &oapiGen.InvalidTraining{}

	if !isDayOnDate(*t.Day, t.Date.Time) {
		errMsg := fmt.Sprintf("Date '%s' isn't on '%s'", t.Date.Format("02.01.2006"), *t.Day)
		invalid.Date = &errMsg
	}

	invalidSessionData := ValidateSession(
		oapiGen.Session{Day: *t.Day, DurationMin: *t.DurationMin, StartTime: *t.StartTime},
	)
	invalid.Day = invalidSessionData.Day
	invalid.StartTime = invalidSessionData.StartTime
	invalid.DurationMin = invalidSessionData.DurationMin
	*invalid.Blocks = validateBlocks(t.Blocks)
	return invalid
}

func ValidateTrainingWithSession(t oapiGen.Training, s oapiGen.Session) *oapiGen.InvalidTraining {
	invalid := &oapiGen.InvalidTraining{}
	if !isDayOnDate(s.Day, t.Date.Time) {
		errMsg := fmt.Sprintf("Date '%s' isn't on '%s'", t.Date.Format("02.01.2006"), s.Day)
		invalid.Date = &errMsg
	}
	*invalid.Blocks = validateBlocks(t.Blocks)
	return invalid
}

func validateBlocks(blocks []oapiGen.Block) []oapiGen.InvalidBlock {
	invalid := make([]oapiGen.InvalidBlock, 0)
	for _, b := range blocks {
		invalidBlock := validateBlock(b)
		if invalidBlock == nil {
			continue
		}
		invalid = append(invalid, *invalidBlock)
	}
	return invalid
}

func isDayOnDate(day oapiGen.Day, date time.Time) bool {
	return strings.ToLower(date.Weekday().String()) != string(day)
}

func validateBlock(b oapiGen.Block) *oapiGen.InvalidBlock {
	invalid := &oapiGen.InvalidBlock{Num: &b.Num}
	if len(b.Name) > MaxNameLen {
		errMsg := fmt.Sprintf("Name must have less than %d characters", MaxNameLen)
		invalid.Name = &errMsg
	}

	if b.Repeat <= 0 {
		errMsg := "Repeat must be greater than 0"
		invalid.Repeat = &errMsg
	}

	*invalid.Sets = validateSets(b.Sets)
	return invalid
}

func validateSets(sets []oapiGen.Set) []oapiGen.InvalidSet {
	invalid := make([]oapiGen.InvalidSet, 0)
	for _, s := range sets {
		invalidSet := validateSet(s)
		if invalidSet == nil {
			continue
		}
		invalid = append(invalid, *invalidSet)
	}
	return invalid
}

func validateSet(s oapiGen.Set) *oapiGen.InvalidSet {
	invalid := &oapiGen.InvalidSet{Num: &s.Num}

	if !startingRules[strings.ToLower(string(s.StartingRule.Rule))] {
		errMsg := fmt.Sprintf("Unkwnown starting rule name '%s'", s.StartingRule.Rule)
		invalid.StartingRule.Rule = &errMsg
	}

	if s.StartingRule.Rule == oapiGen.Pause || s.StartingRule.Rule == oapiGen.Interval {
		if s.StartingRule.Seconds == nil {
			errMsg := fmt.Sprintf("Rule '%s' must have seconds set", string(s.StartingRule.Rule))
			invalid.StartingRule.Seconds = &errMsg
		}
	}

	if s.Distance <= 0 {
		errMsg := "Distance must be greater than 0"
		invalid.Distance = &errMsg
	}

	if s.Repeat <= 0 {
		errMsg := "Repeat must be greater than 0"
		invalid.Distance = &errMsg
	}

	return invalid
}
