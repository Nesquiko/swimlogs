package app

import (
	"strings"
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_validateNewTraining(t *testing.T) {
	mondayDate := types.Date{Time: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)}
	testCases := []struct {
		desc    string
		newT    openapi.NewTraining
		wantErr *openapi.ErrorDetail
	}{
		{
			desc: "valid training",
			newT: openapi.NewTraining{
				Date:        mondayDate,
				StartTime:   "00:00",
				DurationMin: 1,
				Blocks:      []openapi.NewBlock{},
			},
			wantErr: nil,
		},
		{
			desc: "invalid start time",
			newT: openapi.NewTraining{
				Date:        mondayDate,
				StartTime:   "24:00",
				DurationMin: 1,
				Blocks:      []openapi.NewBlock{},
			},
			wantErr: &openapi.ErrorDetail{
				Title:  "Invalid training",
				Detail: "Training contains invalid values",
				Extensions: &map[string]any{
					"startTime": "Start time must be from 00:00 to 23:59, but was '24:00'",
				},
			},
		},
		{
			desc: "invalid duration",
			newT: openapi.NewTraining{
				Date:        mondayDate,
				StartTime:   "00:00",
				DurationMin: 0,
				Blocks:      []openapi.NewBlock{},
			},
			wantErr: &openapi.ErrorDetail{
				Title:  "Invalid training",
				Detail: "Training contains invalid values",
				Extensions: &map[string]any{
					"durationMin": "Duration must be greater than 0",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validateNewTraining(tC.newT)
			assert.Equal(t, tC.wantErr, err)
		})
	}
}

func Test_validateNewBlock(t *testing.T) {
	zero := 0
	testCases := []struct {
		desc string
		newB openapi.NewBlock
		want *openapi.InvalidBlock
	}{
		{
			desc: "valid block",
			newB: openapi.NewBlock{
				Num:    0,
				Name:   "test",
				Repeat: 1,
			},
			want: nil,
		},
		{
			desc: "long name",
			newB: openapi.NewBlock{
				Num:    0,
				Name:   strings.Repeat("a", 256),
				Repeat: 1,
			},
			want: &openapi.InvalidBlock{
				Num:  &zero,
				Name: &blockLongNameErr,
			},
		},
		{
			desc: "invalid repeat",
			newB: openapi.NewBlock{
				Num:    0,
				Name:   "test",
				Repeat: 0,
			},
			want: &openapi.InvalidBlock{
				Num:    &zero,
				Repeat: &repeatErr,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ib := validateNewBlock(tC.newB)
			assert.Equal(t, tC.want, ib)
		})
	}
}

func Test_validateNewSet(t *testing.T) {
	zero := 0
	setStartingRuleTypeErr := "Unknown starting rule name 'unknown'"
	setStartingRuleSecsErr := "Rule 'Interval' must have seconds set"
	testCases := []struct {
		desc string
		newS openapi.NewTrainingSet
		want *openapi.InvalidTrainingSet
	}{
		{
			desc: "valid set",
			newS: openapi.NewTrainingSet{
				Num:          0,
				Repeat:       1,
				Distance:     400,
				What:         "test",
				StartingRule: openapi.StartingRule{Type: openapi.None},
			},
			want: nil,
		},
		{
			desc: "unknown starting rule",
			newS: openapi.NewTrainingSet{
				Num:          0,
				Repeat:       1,
				Distance:     400,
				What:         "test",
				StartingRule: openapi.StartingRule{Type: openapi.StartingRuleType("unknown")},
			},
			want: &openapi.InvalidTrainingSet{
				Num: &zero,
				StartingRule: &struct {
					Seconds *string "json:\"seconds,omitempty\""
					Type    *string "json:\"type,omitempty\""
				}{Type: &setStartingRuleTypeErr},
			},
		},
		{
			desc: "starting rule seconds",
			newS: openapi.NewTrainingSet{
				Num:          0,
				Repeat:       1,
				Distance:     400,
				What:         "test",
				StartingRule: openapi.StartingRule{Type: openapi.Interval},
			},
			want: &openapi.InvalidTrainingSet{
				Num: &zero,
				StartingRule: &struct {
					Seconds *string "json:\"seconds,omitempty\""
					Type    *string "json:\"type,omitempty\""
				}{Seconds: &setStartingRuleSecsErr},
			},
		},
		{
			desc: "invalid distance",
			newS: openapi.NewTrainingSet{
				Num:          0,
				Repeat:       1,
				Distance:     0,
				What:         "test",
				StartingRule: openapi.StartingRule{Type: openapi.None},
			},
			want: &openapi.InvalidTrainingSet{
				Num:      &zero,
				Distance: &setDistanceErr,
			},
		},
		{
			desc: "invalid repeat",
			newS: openapi.NewTrainingSet{
				Num:          0,
				Repeat:       0,
				Distance:     400,
				What:         "test",
				StartingRule: openapi.StartingRule{Type: openapi.None},
			},
			want: &openapi.InvalidTrainingSet{
				Num:    &zero,
				Repeat: &repeatErr,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			is := validateNewSet(tC.newS)
			assert.Equal(t, tC.want, is)
		})
	}
}

func Test_validateNewSession(t *testing.T) {
	testCases := []struct {
		desc    string
		newSess openapi.CreateSessionJSONBody
		wantErr bool
	}{
		{
			desc: "valid session",
			newSess: openapi.CreateSessionJSONBody{
				StartTime:   "23:59",
				Day:         openapi.Friday,
				DurationMin: 60,
			},
			wantErr: false,
		},
		{
			desc: "invalid time",
			newSess: openapi.CreateSessionJSONBody{
				StartTime:   "24:00",
				Day:         openapi.Friday,
				DurationMin: 60,
			},
			wantErr: true,
		},
		{
			desc: "invalid day",
			newSess: openapi.CreateSessionJSONBody{
				StartTime:   "23:59",
				Day:         openapi.Day("invalid"),
				DurationMin: 60,
			},
			wantErr: true,
		},
		{
			desc: "invalid duration",
			newSess: openapi.CreateSessionJSONBody{
				StartTime:   "23:59",
				Day:         openapi.Friday,
				DurationMin: -1,
			},
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			gotErr := validateNewSession(tC.newSess)
			assert.Equal(t, tC.wantErr, gotErr != nil)
		})
	}
}

func Test_isTimeValid(t *testing.T) {
	testCases := []struct {
		desc string
		time string
		want bool
	}{
		{
			desc: "valid time",
			time: "23:59",
			want: true,
		},
		{
			desc: "midnight",
			time: "00:00",
			want: true,
		},
		{
			desc: "more ':'",
			time: "23:59:00",
			want: false,
		},
		{
			desc: "big hour",
			time: "25:00",
			want: false,
		},
		{
			desc: "small hour",
			time: "-1:00",
			want: false,
		},
		{
			desc: "unparsable hour",
			time: "2a:00",
			want: false,
		},
		{
			desc: "big minute",
			time: "23:60",
			want: false,
		},
		{
			desc: "small minute",
			time: "23:-1",
			want: false,
		},
		{
			desc: "unparsable minute",
			time: "23:6a",
			want: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := isTimeValid(tC.time)
			assert.Equal(t, tC.want, got)
		})
	}
}

func Test_validateSessionUpdate(t *testing.T) {
	validStartTime := openapi.StartTime("23:59")
	invalidStartTime := openapi.StartTime("24:00")
	validDay := openapi.Friday
	invalidDay := openapi.Day("invalid")
	validDur := 60
	invalidDur := -1

	testCases := []struct {
		desc    string
		sess    openapi.UpdateSessionJSONBody
		wantErr bool
	}{
		{
			desc: "valid update",
			sess: openapi.UpdateSessionJSONBody{
				StartTime:   &validStartTime,
				Day:         &validDay,
				DurationMin: &validDur,
			},
			wantErr: false,
		},
		{
			desc: "invalid start time",
			sess: openapi.UpdateSessionJSONBody{
				StartTime:   &invalidStartTime,
				Day:         &validDay,
				DurationMin: &validDur,
			},
			wantErr: true,
		},
		{
			desc: "invalid day",
			sess: openapi.UpdateSessionJSONBody{
				StartTime:   nil,
				Day:         &invalidDay,
				DurationMin: &validDur,
			},
			wantErr: true,
		},
		{
			desc: "invalid duration",
			sess: openapi.UpdateSessionJSONBody{
				StartTime:   &validStartTime,
				Day:         nil,
				DurationMin: &invalidDur,
			},
			wantErr: true,
		},
		{
			desc: "all nil",
			sess: openapi.UpdateSessionJSONBody{
				StartTime:   nil,
				Day:         nil,
				DurationMin: nil,
			},
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			gotErr := validateSessionUpdate(tC.sess) != nil
			assert.Equal(t, tC.wantErr, gotErr)
		})
	}
}

func Test_areAllFieldsNil(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	nonNil := 1
	testCases := []struct {
		desc  string
		strct any
		want  bool
	}{
		{
			desc: "all nil",
			strct: struct {
				a *int
				b *string
			}{a: nil, b: nil},
			want: true,
		},
		{
			desc:  "no nils",
			strct: struct{ a *int }{a: &nonNil},
			want:  false,
		},
		{
			desc: "mixed",
			strct: struct {
				a *int
				b string
			}{a: nil, b: "b"},
			want: true,
		},
		{
			desc: "session update",
			strct: openapi.UpdateSessionJSONBody{
				StartTime:   nil,
				Day:         nil,
				DurationMin: nil,
			},
			want: true,
		},
		{
			desc:  "pointer",
			strct: &openapi.InvalidTrainingSet{},
			want:  false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := allFieldsNill(tC.strct)
			assert.Equal(t, tC.want, got)
		})
	}
}
