package validation

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

func ValidateSession(s oapiGen.Session) *oapiGen.InvalidSession {
	invalid := &oapiGen.InvalidSession{}

	if !days[strings.ToLower(string(s.Day))] {
		errMsg := fmt.Sprintf("Unknown day name '%s'", s.Day)
		invalid.Day = &errMsg
	}

	if !isTimeValid(s.StartTime) {
		errMsg := fmt.Sprintf("Start time must be from 00:00 to 23:59, but was '%s'", s.StartTime)
		invalid.StartTime = &errMsg
	}

	if s.DurationMin <= 0 {
		errMsg := "Duration can't be less than 1"
		invalid.DurationMin = &errMsg
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
