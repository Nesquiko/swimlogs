package app

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/rs/zerolog/log"
)

const (
	StartTimeErrFormat = "Start time must be from 00:00 to 23:59, but was '%s'"
	DayErrFormat       = "Unknown day name '%s'"
)

func validateSessionUpdate(newSess openapi.UpdateSessionJSONBody) *openapi.ErrorDetail {
	if allFieldsNill(newSess) {
		return &openapi.ErrorDetail{
			Title:  "Invalid session",
			Detail: "Session contains no values",
		}
	}

	errors := make(map[string]any)
	if newSess.Day != nil {
		if _, ok := openapi.DaysOfWeek[*newSess.Day]; !ok {
			errors["day"] = fmt.Sprintf("Unknown day name '%s'", *newSess.Day)
		}
	}
	if newSess.StartTime != nil && !isTimeValid(string(*newSess.StartTime)) {
		errors["startTime"] = fmt.Sprintf(
			"Start time must be from 00:00 to 23:59, but was '%s'",
			*newSess.StartTime,
		)
	}
	if newSess.DurationMin != nil && *newSess.DurationMin <= 0 {
		errors["durationMin"] = "Duration can't be 0 or less"
	}

	if len(errors) == 0 {
		return nil
	}

	return &openapi.ErrorDetail{
		Title:      "Invalid session",
		Detail:     "Session contains invalid values",
		Extensions: &errors,
	}
}

type SessionValidation struct {
	InvalidSession openapi.InvalidSession
	IsValid        bool
}

func validateNewSession(newSess openapi.CreateSessionJSONBody) SessionValidation {
	invalid := openapi.InvalidSession{}

	if _, ok := openapi.DaysOfWeek[newSess.Day]; !ok {
		invalid.Day = asPtr(fmt.Sprintf(DayErrFormat, newSess.Day))
	}
	if !isTimeValid(newSess.StartTime) {
		invalid.StartTime = asPtr(fmt.Sprintf(StartTimeErrFormat, newSess.StartTime))
	}
	if newSess.DurationMin <= 0 {
		invalid.DurationMin = &durationErr
	}

	sv := SessionValidation{InvalidSession: invalid, IsValid: true}
	if !allFieldsNill(invalid) {
		sv.IsValid = false
	}
	return sv
}

func isTimeValid(startTime string) bool {
	time := strings.Split(startTime, ":")
	if len(time) != 2 {
		return false
	}
	hours, err := strconv.Atoi(time[0])
	if err != nil {
		return false
	} else if hours < 0 || hours > 23 {
		return false
	}

	minutes, err := strconv.Atoi(time[1])
	if err != nil {
		return false
	} else if minutes < 0 || minutes > 59 {
		return false
	}

	return true
}

func allFieldsNill(v any) bool {
	structType := reflect.TypeOf(v)
	if structType.Kind() != reflect.Struct {
		log.Warn().Msgf("allFieldsNill: %v is not a struct", v)
		return false
	}

	structVal := reflect.ValueOf(v)
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		if field.Kind() != reflect.Pointer {
			continue
		}
		if !field.IsNil() {
			return false
		}
	}

	return true
}
