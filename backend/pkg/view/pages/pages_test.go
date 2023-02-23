package pages

import (
	"reflect"
	"testing"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

func Test_splitIntoDays(t *testing.T) {
	details := []oapiGen.TrainingDetail{
		{Day: oapiGen.Monday},
		{Day: oapiGen.Sunday},
		{Day: oapiGen.Friday},
		{Day: oapiGen.Sunday},
		{Day: oapiGen.Wednesday},
	}

	expectedDaysOrdererd := []oapiGen.Day{
		oapiGen.Monday,
		oapiGen.Wednesday,
		oapiGen.Friday,
		oapiGen.Sunday,
	}
	expectedWeekMap := map[oapiGen.Day][]oapiGen.TrainingDetail{
		oapiGen.Monday:    {{Day: oapiGen.Monday}},
		oapiGen.Wednesday: {{Day: oapiGen.Wednesday}},
		oapiGen.Sunday:    {{Day: oapiGen.Sunday}, {Day: oapiGen.Sunday}},
		oapiGen.Friday:    {{Day: oapiGen.Friday}},
	}

	daysOrdered, weekMap := splitIntoDays(details)

	if !reflect.DeepEqual(expectedDaysOrdererd, daysOrdered) {
		t.Fatalf(
			"days ordered not equal, expected %v, but was %v",
			expectedDaysOrdererd,
			daysOrdered,
		)
	}
	if !reflect.DeepEqual(expectedWeekMap, weekMap) {
		t.Fatalf("weeks not equal, expected %v, but was %v", expectedWeekMap, weekMap)
	}
}
