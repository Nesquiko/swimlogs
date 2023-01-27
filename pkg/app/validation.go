package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

const MaxNameLen = 255

// days is a set of day names
var days = map[string]bool{
	string(oapiGen.Monday):    true,
	string(oapiGen.Tuesday):   true,
	string(oapiGen.Wednesday): true,
	string(oapiGen.Thursday):  true,
	string(oapiGen.Friday):    true,
	string(oapiGen.Saturday):  true,
	string(oapiGen.Sunday):    true,
}

// startingRules is a set of allowed startingRule names
var startingRules = map[string]bool{
	string(oapiGen.None):     true,
	string(oapiGen.Pause):    true,
	string(oapiGen.Interval): true,
}

// ValidateSession returns a map with keys being names of session fields, which
// aren't valid, and values being reasons why they aren't valid.
func ValidateSession(s oapiGen.Session) map[string]string {
	return validateSessionData(string(s.Day), s.StartTime, s.DurationMin)
}

func validateSessionData(day, startTime string, durationMin int) map[string]string {
	invalid := make(map[string]string)

	if !days[strings.ToLower(day)] {
		invalid["day"] = fmt.Sprintf("Unknown day name '%s'", day)
	}

	if !isTimeValid(startTime) {
		invalid["startTime"] = fmt.Sprintf(
			"Start time must be from 00:00 to 23:59, but was '%s'",
			startTime,
		)
	}

	if durationMin <= 0 {
		invalid["durationMin"] = "Duration can't be less than 1"
	}

	return invalid
}

func isTimeValid(startTime string) bool {
	time := strings.Split(startTime, ":")
	if len(time) != 2 {
		return false
	}
	hours, err := strconv.Atoi(time[0])
	if err != nil {
		return false
	} else if hours < 0 || hours > 23 {
		return false
	}

	minutes, err := strconv.Atoi(time[1])
	if err != nil {
		return false
	} else if minutes < 0 || minutes > 59 {
		return false
	}

	return true
}

func validateTraining(t oapiGen.Training) map[string]string {
	var invalid map[string]string
	if t.SessionId == nil {
		invalid = validateSessionData(string(*t.Day), *t.StartTime, *t.DurationMin)
		if _, ok := invalid["day"]; !ok &&
			strings.ToLower(t.Date.Weekday().String()) != string(*t.Day) {
			invalid["day"] = fmt.Sprintf(
				"Date '%s' isn't on '%s'",
				t.Date.Format("02.01.2006"),
				string(*t.Day),
			)
		}
	} else {
		invalid = make(map[string]string)
	}

	if len(t.Blocks) == 0 {
		invalid["blocks"] = "No blocks in training"
	}

	for i, b := range t.Blocks {
		invalidBlock := validateBlock(b)
		for k, v := range invalidBlock {
			invalid[fmt.Sprintf("block#%d-%s", i, k)] = v
		}
	}
	return invalid
}

func (app *swimLogsApp) validateSessionInTraining(t oapiGen.Training) map[string]string {
	s, err := app.db.GetSessionById(*t.SessionId)
	if err != nil {
		return map[string]string{"session": "Unknown session"}
	}

	if strings.ToLower(t.Date.Weekday().String()) != string(s.Day) {
		return map[string]string{
			"day": fmt.Sprintf(
				"Date '%s' isn't on '%s'",
				t.Date.Format("02.01.2006"),
				string(*t.Day),
			),
		}
	}

	return map[string]string{}
}

func validateBlock(b oapiGen.Block) map[string]string {
	invalid := make(map[string]string)
	if len(b.Name) > MaxNameLen {
		invalid["name"] = fmt.Sprintf("Name must have less than %d characters", MaxNameLen)
	}

	if b.Repeat <= 0 {
		invalid["repeat"] = "Repeat must be greater than 0"
	}

	if len(b.Sets) == 0 {
		invalid["sets"] = "No sets in block"
	}

	for i, s := range b.Sets {
		invalidSet := validateSet(s)
		for k, v := range invalidSet {
			invalid[fmt.Sprintf("set#%d-%s", i, k)] = v
		}
	}
	return invalid
}

func validateSet(s oapiGen.Set) map[string]string {
	invalid := make(map[string]string)

	if !startingRules[strings.ToLower(string(s.StartingRule.Rule))] {
		invalid["startingRule"] = fmt.Sprintf(
			"Unkwnown starting rule name '%s'",
			s.StartingRule.Rule,
		)
	}

	if s.StartingRule.Rule == oapiGen.Pause || s.StartingRule.Rule == oapiGen.Interval {
		if s.StartingRule.Seconds == nil {
			invalid["startingRule"] = fmt.Sprintf(
				"Rule '%s' must have seconds set",
				string(s.StartingRule.Rule),
			)
		}
	}

	if s.Distance <= 0 {
		invalid["distance"] = "Distance must be greater than 0"
	}

	if s.Repeat <= 0 {
		invalid["repeat"] = "Repeat must be greater than 0"
	}

	return invalid
}
