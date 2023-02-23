package pages

import (
	"sort"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/view/util"
)

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

	for d := range daysMap {
		sort.Slice(daysMap[d], func(i, j int) bool {
			return daysMap[d][i].StartTime < daysMap[d][j].StartTime
		})
	}

	sort.Slice(days, func(i, j int) bool {
		return util.DaysOrder[days[i]] < util.DaysOrder[days[j]]
	})

	return days, daysMap
}
