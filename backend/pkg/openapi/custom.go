package openapi

import (
	"fmt"
	"time"
)

const FormatStartTime = "15:04"

type StartTime string

func (s StartTime) String() string {
	return string(s)
}

// Scans time.Time into a string with format "HH:MM"
func (s *StartTime) Scan(value any) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("invalid type: %T", value)
	}
	*s = StartTime(t.Format(FormatStartTime))
	return nil
}

var DaysOfWeek = map[Day]int{
	Sunday:    0,
	Monday:    1,
	Tuesday:   2,
	Wednesday: 3,
	Thursday:  4,
	Friday:    5,
	Saturday:  6,
}

var StartTypes = map[StartType]bool{
	None:     true,
	Pause:    true,
	Interval: true,
}
