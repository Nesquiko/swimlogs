package pages

import (
	"sort"

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

// splitIntoDays splits passed details in to a map, in which keys are days of
// the week, and values are TrainingDetails which are in that day. Also returns
// sorted slice of days which are in map.
func splitIntoDays(
	details []oapiGen.TrainingDetail,
) ([]oapiGen.Day, map[oapiGen.Day][]oapiGen.TrainingDetail) {
	days := make([]oapiGen.Day, 0)
	daysMap := make(map[oapiGen.Day][]oapiGen.TrainingDetail)

	for _, td := range details {
		if _, ok := daysMap[td.Day]; !ok {
			days = append(days, td.Day)
			daysMap[td.Day] = make([]oapiGen.TrainingDetail, 0)
		}

		daysMap[td.Day] = append(daysMap[td.Day], td)
	}

	sort.Slice(days, func(i, j int) bool {
		return DaysOrder[days[i]] < DaysOrder[days[j]]
	})

	return days, daysMap
}
