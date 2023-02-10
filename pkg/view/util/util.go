package util

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

var DaysOrder = map[oapiGen.Day]int{
	oapiGen.Monday:    0,
	oapiGen.Tuesday:   1,
	oapiGen.Wednesday: 2,
	oapiGen.Thursday:  3,
	oapiGen.Friday:    4,
	oapiGen.Saturday:  5,
	oapiGen.Sunday:    6,
}

func OrderByDays(sessions *[]oapiGen.Session) {
	if sessions == nil {
		return
	}

	s := *sessions
	sort.Slice(s, func(i, j int) bool {
		if DaysOrder[s[i].Day] == DaysOrder[s[j].Day] {
			return isStartTimeLess(s[i].StartTime, s[j].StartTime)
		}
		return DaysOrder[s[i].Day] < DaysOrder[s[j].Day]
	})
}

func isStartTimeLess(st1, st2 string) bool {
	time1Parts := strings.Split(st1, ":")
	time2Parts := strings.Split(st2, ":")

	hours1, _ := strconv.Atoi(time1Parts[0])
	minutes1, _ := strconv.Atoi(time1Parts[1])

	hours2, _ := strconv.Atoi(time2Parts[0])
	minutes2, _ := strconv.Atoi(time2Parts[1])

	totalMinutes1 := hours1*60 + minutes1
	totalMinutes2 := hours2*60 + minutes2

	return totalMinutes1 < totalMinutes2
}
