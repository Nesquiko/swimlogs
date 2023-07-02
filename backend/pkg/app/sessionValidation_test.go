package app

import (
	"testing"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

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
