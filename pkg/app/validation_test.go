package app

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

func Test_validateTraining(t *testing.T) {
	id := uuid.New()
	testCases := []struct {
		desc     string
		t        oapiGen.Training
		expected map[string]string
	}{
		{
			desc: "InvalidSessionData",
			t: oapiGen.Training{
				SessionId:   nil,
				Day:         &invalidSession.Day,
				StartTime:   &invalidSession.StartTime,
				DurationMin: &invalidSession.DurationMin,
				Blocks:      []oapiGen.Block{validBlock},
			},
			expected: map[string]string{
				"day":         "Unknown day name 'mOnDy'",
				"durationMin": "Duration can't be less than 1",
				"startTime":   "Start time must be from 00:00 to 23:59, but was '24:00'",
			},
		},
		{
			desc: "NoBlocks",
			t: oapiGen.Training{
				SessionId: &id,
				Blocks:    []oapiGen.Block{},
			},
			expected: map[string]string{"blocks": "No blocks in training"},
		},
		{
			desc: "Block#0-Set#0-Distance=0",
			t: oapiGen.Training{
				SessionId: &id,
				Blocks:    []oapiGen.Block{invalidBlock},
			},
			expected: map[string]string{
				"block#0-name":           "Name must have less than 255 characters",
				"block#0-set#0-distance": "Distance must be greater than 0",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			invalid := validateTraining(tC.t)

			if !reflect.DeepEqual(invalid, tC.expected) {
				t.Errorf("Invalid fields, expected %v, but was %v", tC.expected, invalid)
			}

		})
	}
}

func Test_validateBlock(t *testing.T) {
	testCases := []struct {
		desc     string
		b        oapiGen.Block
		expected map[string]string
	}{
		{
			desc: "LongName",
			b: oapiGen.Block{
				Name:   strings.Repeat("x", 256),
				Repeat: 1,
				Sets:   []oapiGen.Set{validSet},
			},
			expected: map[string]string{"name": "Name must have less than 255 characters"},
		},
		{
			desc: "Repeat=0",
			b: oapiGen.Block{
				Name:   strings.Repeat("x", 255),
				Repeat: 0,
				Sets:   []oapiGen.Set{validSet},
			},
			expected: map[string]string{"repeat": "Repeat must be greater than 0"},
		},
		{
			desc: "EmptySets",
			b: oapiGen.Block{
				Name:   strings.Repeat("x", 255),
				Repeat: 1,
				Sets:   []oapiGen.Set{},
			},
			expected: map[string]string{"sets": "No sets in block"},
		},
		{
			desc: "Set#0-Distance=0",
			b: oapiGen.Block{
				Name:   strings.Repeat("x", 255),
				Repeat: 1,
				Sets:   []oapiGen.Set{invalidSet},
			},
			expected: map[string]string{"set#0-distance": "Distance must be greater than 0"},
		},
		{
			desc: "Valid",
			b: oapiGen.Block{
				Name:   strings.Repeat("x", 255),
				Repeat: 1,
				Sets:   []oapiGen.Set{validSet},
			},
			expected: map[string]string{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			invalid := validateBlock(tC.b)

			if !reflect.DeepEqual(invalid, tC.expected) {
				t.Errorf("Invalid fields, expected %v, but was %v", tC.expected, invalid)
			}
		})
	}
}

func Test_validateSet(t *testing.T) {
	testCases := []struct {
		desc     string
		s        oapiGen.Set
		expected map[string]string
	}{
		{
			desc: "UnknownStartingRule",
			s: oapiGen.Set{
				StartingRule: oapiGen.StartingRule{Rule: oapiGen.StartingRuleRule("heh")},
			},
			expected: map[string]string{
				"startingRule": "Unkwnown starting rule name 'heh'",
				"distance":     "Distance must be greater than 0",
				"repeat":       "Repeat must be greater than 0",
			},
		},
		{
			desc: "NotSetSecondsPause",
			s: oapiGen.Set{
				StartingRule: oapiGen.StartingRule{Rule: oapiGen.Pause},
			},
			expected: map[string]string{
				"startingRule": "Rule 'pause' must have seconds set",
				"distance":     "Distance must be greater than 0",
				"repeat":       "Repeat must be greater than 0",
			},
		},
		{
			desc: "NotSetSecondsInterval",
			s: oapiGen.Set{
				StartingRule: oapiGen.StartingRule{Rule: oapiGen.Interval},
			},
			expected: map[string]string{
				"startingRule": "Rule 'interval' must have seconds set",
				"distance":     "Distance must be greater than 0",
				"repeat":       "Repeat must be greater than 0",
			},
		},
		{
			desc: "Valid",
			s: oapiGen.Set{
				StartingRule: oapiGen.StartingRule{Rule: oapiGen.None},
				Distance:     100,
				Repeat:       2,
				What:         "Freestyle",
			},
			expected: map[string]string{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			invalid := validateSet(tC.s)

			if !reflect.DeepEqual(invalid, tC.expected) {
				t.Errorf("Invalid fields, expected %v, but was %v", tC.expected, invalid)
			}
		})
	}
}

func Test_validateSession(t *testing.T) {
	testCases := []struct {
		desc     string
		s        oapiGen.Session
		expected map[string]string
	}{
		{
			desc: "Invalid",
			s:    oapiGen.Session{Day: oapiGen.Day("mOnDy"), StartTime: "24:00", DurationMin: 0},
			expected: map[string]string{
				"day":         "Unknown day name 'mOnDy'",
				"durationMin": "Duration can't be less than 1",
				"startTime":   "Start time must be from 00:00 to 23:59, but was '24:00'",
			},
		},
		{
			desc:     "Valid",
			s:        oapiGen.Session{Day: oapiGen.Monday, StartTime: "23:59", DurationMin: 60},
			expected: map[string]string{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			invalidFields := validateSession(tC.s)

			if !reflect.DeepEqual(invalidFields, tC.expected) {
				t.Errorf("Invalid fields, expected %v, but was %v", tC.expected, invalidFields)
			}
		})
	}
}

func Test_isTimeValid(t *testing.T) {
	testCases := []struct {
		desc      string
		startTime string
		expected  bool
	}{
		{
			desc:      "RandomString",
			startTime: "faoi2j5x",
			expected:  false,
		},
		{
			desc:      "TooManyColons",
			startTime: "A:B:C",
			expected:  false,
		},
		{
			desc:      "HighHour",
			startTime: "24:13",
			expected:  false,
		},
		{
			desc:      "LowHour",
			startTime: "-1:13",
			expected:  false,
		},
		{
			desc:      "HighMinutes",
			startTime: "0:60",
			expected:  false,
		},
		{
			desc:      "LowMinutes",
			startTime: "0:-1",
			expected:  false,
		},
		{
			desc:      "InvalidHours",
			startTime: "AA:10",
			expected:  false,
		},
		{
			desc:      "InvalidMinutes",
			startTime: "00:BB",
			expected:  false,
		},
		{
			desc:      "Midnight",
			startTime: "0:0",
			expected:  true,
		},
		{
			desc:      "BeforeMidnight",
			startTime: "23:59",
			expected:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			isValid := isTimeValid(tC.startTime)

			if isValid != tC.expected {
				t.Errorf("Time %q, expected %t, but was %t", tC.startTime, tC.expected, isValid)
			}
		})
	}
}

var (
	invalidSession = oapiGen.Session{Day: oapiGen.Day("mOnDy"), StartTime: "24:00", DurationMin: 0}
	validSet       = oapiGen.Set{
		Distance:     400,
		Repeat:       1,
		What:         "Freestyle",
		StartingRule: oapiGen.StartingRule{Rule: oapiGen.None},
	}
	invalidSet = oapiGen.Set{
		Distance:     0,
		Repeat:       1,
		What:         "Freestyle",
		StartingRule: oapiGen.StartingRule{Rule: oapiGen.None},
	}
	validBlock = oapiGen.Block{
		Name:   strings.Repeat("x", 255),
		Repeat: 1,
		Sets:   []oapiGen.Set{validSet},
	}
	invalidBlock = oapiGen.Block{
		Name:   strings.Repeat("x", 256),
		Repeat: 1,
		Sets:   []oapiGen.Set{invalidSet},
	}
)
