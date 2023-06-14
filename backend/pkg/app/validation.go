package app

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/rs/zerolog/log"
)

const MaxNameLen = 255

var (
	blockLongNameErr = fmt.Sprintf("Name must have less than %d characters", MaxNameLen)
	setDistanceErr   = "Distance must be greater than 0"
	repeatErr        = "Repeat must be greater than 0"
	durationErr      = "Duration must be greater than 0"
)

func validateNewTraining(newTraining openapi.NewTraining) *openapi.ErrorDetail {
	errors := make(map[string]any)
	if !isTimeValid(string(newTraining.StartTime)) {
		errors["startTime"] = fmt.Sprintf(
			"Start time must be from 00:00 to 23:59, but was '%s'",
			newTraining.StartTime,
		)
	}
	if newTraining.DurationMin <= 0 {
		errors["durationMin"] = durationErr
	}

	invalidBlocks := make([]openapi.InvalidBlock, 0)
	for _, block := range newTraining.Blocks {
		invalidBlock := validateNewBlock(block)
		if invalidBlock != nil {
			invalidBlocks = append(invalidBlocks, *invalidBlock)
		}
	}
	if len(newTraining.Blocks) != 0 && len(invalidBlocks) > 0 {
		errors["blocks"] = invalidBlocks
	}

	if len(errors) == 0 {
		return nil
	}
	return &openapi.ErrorDetail{
		Title:      "Invalid training",
		Detail:     "Training contains invalid values",
		Extensions: &errors,
	}
}

func validateNewBlock(newBlock openapi.NewBlock) *openapi.InvalidBlock {
	invalid := openapi.InvalidBlock{}
	if len(newBlock.Name) > MaxNameLen {
		invalid.Name = &blockLongNameErr
	}

	if newBlock.Repeat <= 0 {
		invalid.Repeat = &repeatErr
	}

	invalidSets := make([]openapi.InvalidTrainingSet, 0)
	for _, set := range newBlock.Sets {
		invalidSet := validateNewSet(set)
		if invalidSet != nil {
			invalidSets = append(invalidSets, *invalidSet)
		}
	}
	if len(invalidSets) > 0 {
		invalid.Sets = &invalidSets
	}
	if allFieldsNill(invalid) {
		return nil
	}

	invalid.Num = &newBlock.Num
	return &invalid
}

func validateNewSet(set openapi.NewTrainingSet) *openapi.InvalidTrainingSet {
	invalid := openapi.InvalidTrainingSet{}

	if !openapi.StartingRulesTypes[set.StartingRule.Type] {
		errMsg := fmt.Sprintf("Unknown starting rule name '%s'", set.StartingRule.Type)
		invalid.StartingRule = &struct {
			Seconds *string `json:"seconds,omitempty"`
			Type    *string `json:"type,omitempty"`
		}{Type: &errMsg}
	} else if set.StartingRule.Type == openapi.Pause || set.StartingRule.Type == openapi.Interval {
		if set.StartingRule.Seconds == nil {
			errMsg := fmt.Sprintf("Rule '%s' must have seconds set", set.StartingRule.Type)
			invalid.StartingRule = &struct {
				Seconds *string `json:"seconds,omitempty"`
				Type    *string `json:"type,omitempty"`
			}{Seconds: &errMsg}
		}
	}

	if set.Distance <= 0 {
		invalid.Distance = &setDistanceErr
	}
	if set.Repeat <= 0 {
		invalid.Repeat = &repeatErr
	}

	if allFieldsNill(invalid) {
		return nil
	}

	invalid.Num = &set.Num
	return &invalid
}

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

func validateNewSession(newSess openapi.CreateSessionJSONBody) *openapi.ErrorDetail {
	errors := make(map[string]any)

	if _, ok := openapi.DaysOfWeek[newSess.Day]; !ok {
		errors["day"] = fmt.Sprintf("Unknown day name '%s'", newSess.Day)
	}
	if !isTimeValid(string(newSess.StartTime)) {
		errors["startTime"] = fmt.Sprintf(
			"Start time must be from 00:00 to 23:59, but was '%s'",
			newSess.StartTime,
		)
	}
	if newSess.DurationMin <= 0 {
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
