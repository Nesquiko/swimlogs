package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
)

var days = map[string]bool{
	string(oapiGen.Monday):    true,
	string(oapiGen.Tuesday):   true,
	string(oapiGen.Wednesday): true,
	string(oapiGen.Thursday):  true,
	string(oapiGen.Friday):    true,
	string(oapiGen.Saturday):  true,
	string(oapiGen.Sunday):    true,
}

func validateSession(session *oapiGen.Session) map[string]string {
	invalid := make(map[string]string)

	if !days[strings.ToLower(string(session.Day))] {
		invalid["day"] = fmt.Sprintf("Unkwnown day name %q", session.Day)
	}

	if !isTimeValid(session.StartTime) {
		invalid["startTime"] = fmt.Sprintf(
			"Start time must be from 00:00 to 23:59, but was %q",
			session.StartTime,
		)
	}

	if session.DurationMin <= 0 {
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
	} else if minutes < 0 || hours > 59 {
		return false
	}

	return true
}
