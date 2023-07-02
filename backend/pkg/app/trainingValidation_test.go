package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func Test_validateNewTraining(t *testing.T) {
	testCases := []struct {
		desc     string
		tr       openapi.NewTraining
		expected TrainingValidation
	}{
		{
			desc: "valid training",
			tr: openapi.NewTraining{
				DurationMin: 60,
				Sets: []openapi.NewTrainingSet{
					{
						Repeat:    1,
						SetOrder:  asPtr(0),
						StartType: openapi.None,
					},
				},
			},
			expected: TrainingValidation{
				InvalidTraining: openapi.InvalidTraining{
					InvalidSets: &[]openapi.InvalidTrainingSet{{}},
				},
				IsValid: true,
			},
		},
		{
			desc: "zero duration",
			tr: openapi.NewTraining{
				DurationMin: 0,
				Sets: []openapi.NewTrainingSet{
					{
						Repeat:    1,
						SetOrder:  asPtr(0),
						StartType: openapi.None,
					},
				},
			},
			expected: TrainingValidation{
				InvalidTraining: openapi.InvalidTraining{
					DurationMin: &durationErr,
					InvalidSets: &[]openapi.InvalidTrainingSet{{}},
				},
				IsValid: false,
			},
		},
		{
			desc: "zero sets",
			tr: openapi.NewTraining{
				DurationMin: 60,
			},
			expected: TrainingValidation{
				InvalidTraining: openapi.InvalidTraining{
					Sets:        &setsErr,
					InvalidSets: &[]openapi.InvalidTrainingSet{},
				},
				IsValid: false,
			},
		},
		{
			desc: "invalid sets",
			tr: openapi.NewTraining{
				DurationMin: 60,
				Sets: []openapi.NewTrainingSet{
					{
						Repeat:    0,
						StartType: openapi.None,
						SetOrder:  asPtr(0),
					},
					{
						Repeat:    1,
						StartType: openapi.Interval,
						SetOrder:  asPtr(0),
					},
				},
			},
			expected: TrainingValidation{
				InvalidTraining: openapi.InvalidTraining{
					InvalidSets: &[]openapi.InvalidTrainingSet{
						{
							Repeat: &repeatErr,
						},
						{
							SetOrder: asPtr(fmt.Sprintf(setOrderDuplicateErrFormat, 0)),
							StartSeconds: asPtr(
								fmt.Sprintf(startSecondsErrFormat, openapi.Interval),
							),
						},
					},
				},
				IsValid: false,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := validateNewTraining(tC.tr)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func Test_validateSets(t *testing.T) {
	training := openapi.NewTraining{
		Sets: []openapi.NewTrainingSet{
			{
				Repeat:    0,
				StartType: openapi.None,
				SetOrder:  asPtr(0),
			},
			{
				Repeat:    1,
				StartType: openapi.None,
				SetOrder:  asPtr(-1),
			},
			{
				Repeat:    1,
				StartType: openapi.None,
				SetOrder:  asPtr(0),
				SubSets: &[]openapi.NewTrainingSet{
					{
						Repeat:      0,
						StartType:   openapi.None,
						SubSetOrder: asPtr(0),
					},
				},
			},
		},
	}
	expected := []openapi.InvalidTrainingSet{
		{Repeat: &repeatErr},
		{SetOrder: &setOrderNegativeErr},
		{
			SetOrder: asPtr(fmt.Sprintf(setOrderDuplicateErrFormat, 0)),
			SubSets: &[]openapi.InvalidTrainingSet{
				{Repeat: &repeatErr},
			},
		},
	}

	actual, isValid := validateSets(training)
	assert.Equal(t, expected, actual)
	assert.False(t, isValid)
}

func Test_validateNewSet(t *testing.T) {
	set := openapi.NewTrainingSet{
		Repeat:       0,
		SetOrder:     asPtr(0),
		StartType:    openapi.Interval,
		StartSeconds: asPtr(20),
		SubSets: &[]openapi.NewTrainingSet{
			{
				Repeat:      0,
				StartType:   openapi.None,
				SubSetOrder: asPtr(-1),
			},
		},
	}
	expected := openapi.InvalidTrainingSet{
		Repeat: &repeatErr,
		SubSets: &[]openapi.InvalidTrainingSet{
			{
				Repeat:      &repeatErr,
				SubSetOrder: &subSetOrderNegativeErr,
			},
		},
	}

	actual, isValid := validateNewSet(set)
	assert.Equal(t, expected, actual)
	assert.False(t, isValid)
}

func Test_validateSubSets(t *testing.T) {
	parent := openapi.NewTrainingSet{
		SubSets: &[]openapi.NewTrainingSet{
			{
				Repeat:      1,
				StartType:   openapi.None,
				SubSetOrder: asPtr(-1),
			},
			{
				Repeat:      0,
				StartType:   openapi.None,
				SubSetOrder: asPtr(0),
			},
			{
				Repeat:      1,
				StartType:   openapi.Pause,
				SubSetOrder: asPtr(0),
			},
		},
	}

	expected := []openapi.InvalidTrainingSet{
		{SubSetOrder: &subSetOrderNegativeErr},
		{Repeat: &repeatErr},
		{
			SubSetOrder:  asPtr(fmt.Sprintf(subSetOrderDuplicateErrFormat, 0)),
			StartSeconds: asPtr(fmt.Sprintf(startSecondsErrFormat, openapi.Pause)),
		},
	}
	actual, areSubSetsValid := validateSubSets(parent)

	assert.Equal(t, expected, actual)
	assert.False(t, areSubSetsValid)
}

func Test_validateSet(t *testing.T) {
	testCases := []struct {
		desc            string
		set             openapi.NewTrainingSet
		expected        openapi.InvalidTrainingSet
		expectedIsValid bool
	}{
		{
			desc: "valid set",
			set: openapi.NewTrainingSet{
				Repeat:    1,
				StartType: openapi.None,
				SetOrder:  asPtr(0),
			},
			expected:        openapi.InvalidTrainingSet{},
			expectedIsValid: true,
		},
		{
			desc: "order not set",
			set: openapi.NewTrainingSet{
				Repeat:    1,
				StartType: openapi.None,
			},
			expected: openapi.InvalidTrainingSet{
				SetOrder:    &orderNotSetErr,
				SubSetOrder: &orderNotSetErr,
			},
			expectedIsValid: false,
		},
		{
			desc: "negative set order",
			set: openapi.NewTrainingSet{
				Repeat:    1,
				StartType: openapi.None,
				SetOrder:  asPtr(-1),
			},
			expected: openapi.InvalidTrainingSet{
				SetOrder: &setOrderNegativeErr,
			},
			expectedIsValid: false,
		},
		{
			desc: "negative set sub order",
			set: openapi.NewTrainingSet{
				Repeat:      1,
				StartType:   openapi.None,
				SubSetOrder: asPtr(-1),
			},
			expected: openapi.InvalidTrainingSet{
				SubSetOrder: &subSetOrderNegativeErr,
			},
			expectedIsValid: false,
		},
		{
			desc: "zero distance",
			set: openapi.NewTrainingSet{
				Repeat:         1,
				StartType:      openapi.None,
				SetOrder:       asPtr(0),
				DistanceMeters: asPtr(0),
			},
			expected: openapi.InvalidTrainingSet{
				DistanceMeters: &distanceErr,
			},
			expectedIsValid: false,
		},
		{
			desc: "zero repeat",
			set: openapi.NewTrainingSet{
				Repeat:    0,
				StartType: openapi.None,
				SetOrder:  asPtr(0),
			},
			expected: openapi.InvalidTrainingSet{
				Repeat: &repeatErr,
			},
			expectedIsValid: false,
		},
		{
			desc: "invalid start type",
			set: openapi.NewTrainingSet{
				Repeat:    1,
				StartType: openapi.StartingRuleType("invalid"),
				SetOrder:  asPtr(0),
			},
			expected: openapi.InvalidTrainingSet{
				StartType: asPtr(fmt.Sprintf(startTypeErrFormat, "invalid")),
			},
			expectedIsValid: false,
		},
		{
			desc: "missing start seconds",
			set: openapi.NewTrainingSet{
				Repeat:    1,
				StartType: openapi.Pause,
				SetOrder:  asPtr(0),
			},
			expected: openapi.InvalidTrainingSet{
				StartSeconds: asPtr(fmt.Sprintf(startSecondsErrFormat, openapi.Pause)),
			},
			expectedIsValid: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, isSetValid := validateSet(tC.set)
			assert.Equal(t, tC.expected, actual)
			assert.Equal(t, tC.expectedIsValid, isSetValid)
		})
	}
}

func Test_validateNewTrainingShouldBeValid(t *testing.T) {
	training := openapi.NewTraining{
		Start:         time.Now(),
		DurationMin:   90,
		TotalDistance: 1900,
		Sets: []openapi.NewTrainingSet{
			{
				SetOrder:       asPtr(0),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  400,
				DistanceMeters: asPtr(400),
				Description:    asPtr("warm up freestyle"),
			},
			{
				SetOrder:       asPtr(1),
				Repeat:         3,
				StartType:      openapi.Pause,
				TotalDistance:  600,
				DistanceMeters: asPtr(200),
				Description:    asPtr("drills"),
				StartSeconds:   asPtr(20),
			},
			{
				SetOrder:      asPtr(2),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 300,
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  50,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  25,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(3),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  100,
				DistanceMeters: asPtr(100),
				Description:    asPtr("cool down"),
			},
			{
				SetOrder:      asPtr(4),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 300,
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  50,
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  25,
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(5),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  200,
				DistanceMeters: asPtr(200),
				Description:    asPtr("breastroke cool down"),
			},
		},
	}

	validation := validateNewTraining(training)
	assert.True(t, validation.IsValid)
}
