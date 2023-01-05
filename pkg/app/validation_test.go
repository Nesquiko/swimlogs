package app

import "testing"

func Test_isTimeValid(t *testing.T) {
	testCases := []struct {
		desc      string
		startTime string
		expect    bool
	}{
		{
			desc:      "RandomString",
			startTime: "faoi2j5x",
			expect:    false,
		},
		{
			desc:      "TooManyColons",
			startTime: "A:B:C",
			expect:    false,
		},
		{
			desc:      "HighHour",
			startTime: "24:13",
			expect:    false,
		},
		{
			desc:      "LowHour",
			startTime: "-1:13",
			expect:    false,
		},
		{
			desc:      "HighMinutes",
			startTime: "0:60",
			expect:    false,
		},
		{
			desc:      "LowMinutes",
			startTime: "0:-1",
			expect:    false,
		},
		{
			desc:      "Midnight",
			startTime: "0:0",
			expect:    true,
		},
		{
			desc:      "BeforeMidnight",
			startTime: "23:59",
			expect:    true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			isValid := isTimeValid(tC.startTime)

			if isValid != tC.expect {
				t.Errorf("Time %q, expected %t, but was %t", tC.startTime, tC.expect, isValid)
			}
		})
	}
}
