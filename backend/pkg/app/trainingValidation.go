package app

import (
	"fmt"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
)

var (
	distanceErr                   = "Distance must be greater than 0"
	repeatErr                     = "Repeat must be greater than 0"
	durationErr                   = "Duration must be greater than 0"
	setsErr                       = "Training must have at least one set"
	orderNotSetErr                = "Either set order or sub set order must be set"
	setOrderNegativeErr           = "Set order cna't be less than 0"
	subSetOrderNegativeErr        = "Sub set order can't be less than 0"
	setOrderDuplicateErrFormat    = "Duplicate set order '%d'"
	subSetOrderDuplicateErrFormat = "Duplicate sub set order '%d'"
	startTypeErrFormat            = "Unknown starting type name '%s'"
	startSecondsErrFormat         = "Type '%s' must have seconds set"
)

type TrainingValidation struct {
	InvalidTraining openapi.InvalidTraining
	IsValid         bool
}

func validateNewTraining(newTraining openapi.NewTraining) TrainingValidation {
	tv := TrainingValidation{IsValid: true}
	invalidTraining := openapi.InvalidTraining{}

	if newTraining.DurationMin <= 0 {
		invalidTraining.DurationMin = &durationErr
		tv.IsValid = false
	}

	if len(newTraining.Sets) == 0 {
		invalidTraining.Sets = &setsErr
		tv.IsValid = false
	}

	invalidSets, isValid := validateSets(newTraining)
	invalidTraining.InvalidSets = &invalidSets
	tv.IsValid = tv.IsValid && isValid

	tv.InvalidTraining = invalidTraining
	return tv
}

func validateSets(t openapi.NewTraining) ([]openapi.InvalidTrainingSet, bool) {
	invalidSets := make([]openapi.InvalidTrainingSet, 0)
	if t.Sets == nil {
		return invalidSets, true
	}

	setOrders := make(map[int]bool)
	isValid := true

	for _, set := range t.Sets {
		is, isSetValid := validateNewSet(set)
		isValid = isValid && isSetValid

		if is.SetOrder == nil { // set order isn't negative
			if setOrders[*set.SetOrder] {
				is.SetOrder = asPtr(fmt.Sprintf(setOrderDuplicateErrFormat, *set.SetOrder))
			} else {
				setOrders[*set.SetOrder] = true
			}
		}

		invalidSets = append(invalidSets, is)
	}

	return invalidSets, isValid
}

func validateNewSet(set openapi.NewTrainingSet) (openapi.InvalidTrainingSet, bool) {
	invalid, isSetValid := validateSet(set)

	if invalidSubsets, areSubsetsValid := validateSubSets(set); len(invalidSubsets) > 0 {
		invalid.SubSets = &invalidSubsets
		isSetValid = isSetValid && areSubsetsValid
	}

	return invalid, isSetValid
}

func validateSubSets(set openapi.NewTrainingSet) ([]openapi.InvalidTrainingSet, bool) {
	if set.SubSets == nil {
		return nil, true
	}

	isValid := true
	invalidSubsets := make([]openapi.InvalidTrainingSet, 0)
	subSetOrders := make(map[int]bool)

	for _, subSet := range *set.SubSets {
		its, isSubSetValid := validateSet(subSet)
		isValid = isValid && isSubSetValid

		if its.SubSetOrder == nil { // sub set order isn't negative
			if subSetOrders[*subSet.SubSetOrder] {
				its.SubSetOrder = asPtr(
					fmt.Sprintf(subSetOrderDuplicateErrFormat, *subSet.SubSetOrder),
				)
			} else {
				subSetOrders[*subSet.SubSetOrder] = true
			}
		}

		invalidSubsets = append(invalidSubsets, its)
	}

	return invalidSubsets, isValid
}

func validateSet(set openapi.NewTrainingSet) (openapi.InvalidTrainingSet, bool) {
	invalid := openapi.InvalidTrainingSet{}
	isValid := true

	if set.SetOrder == nil && set.SubSetOrder == nil {
		invalid.SetOrder = &orderNotSetErr
		invalid.SubSetOrder = &orderNotSetErr
		isValid = false
	}
	if set.SetOrder != nil && *set.SetOrder < 0 {
		invalid.SetOrder = &setOrderNegativeErr
		isValid = false
	}
	if set.SubSetOrder != nil && *set.SubSetOrder < 0 {
		invalid.SubSetOrder = &subSetOrderNegativeErr
		isValid = false
	}
	if set.DistanceMeters != nil && *set.DistanceMeters <= 0 {
		invalid.DistanceMeters = &distanceErr
		isValid = false
	}
	if set.Repeat <= 0 {
		invalid.Repeat = &repeatErr
		isValid = false
	}

	if !openapi.StartTypes[set.StartType] {
		invalid.StartType = asPtr(fmt.Sprintf(startTypeErrFormat, set.StartType))
		isValid = false
	} else if set.StartType != openapi.None && set.StartSeconds == nil {
		invalid.StartSeconds = asPtr(fmt.Sprintf(startSecondsErrFormat, set.StartType))
		isValid = false
	}

	return invalid, isValid
}
