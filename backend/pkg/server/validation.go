package server

import (
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

const (
	InvalidTrainingErrorTitle = "Invalid training"
	InvalidTrainingErrorCode  = "training.invalid"
	DurationErrorDetail       = "Duration must be between 1 and %d, was %d"
	StartErrorDetail          = "Invalid value of start, %q"
	SetsErrorDetail           = "No sets provided"

	InvalidSetErrorTitle            = "Set is not valid"
	InvalidSetErrorCode             = "invalid.set"
	SetOrderErrorDetail             = "Set order must be between 0 and %d, was %d"
	RepeatErrorDetail               = "Repeat must be between 1 and %d, was %d"
	DistanceErrorDetail             = "Distance must be between 1 and %d, was %d"
	StartTypeUnknownErrorDetail     = "Start type must be either %s or %s, was %s"
	StartSecondsRequiredErrorDetail = "Start seconds must be provided"
	StartSecondsErrorDetail         = "Start seconds must be between 1 and %d, was %d"
	EquipmentErrorDetail            = "Unknown Equipment %s"
	GroupErrorDetail                = "Unknown Group %s"

	NonUniqueSetOrderCode   = "invalid.set.nonunique.setorder"
	NonUniqueSetOrderDetail = "Multiple sets had same set order %d"
)

var (
	StartTypesSet = map[apidef.StartTypeEnum]bool{apidef.Interval: true, apidef.Pause: true}
	EquipmentSet  = map[apidef.EquipmentEnum]bool{
		apidef.Board:   true,
		apidef.Fins:    true,
		apidef.Monofin: true,
		apidef.Paddles: true,
		apidef.Snorkel: true,
	}
	GroupSet = map[apidef.GroupEnum]bool{
		apidef.Bifi:   true,
		apidef.Long:   true,
		apidef.Middle: true,
		apidef.Mono:   true,
		apidef.Sprint: true,
	}
)

func validateNewTraining(nt apidef.NewTraining) *ApiError {
	if nt.DurationMin <= 0 || nt.DurationMin > data.SmallIntMax {
		return invalidTraining(fmt.Sprintf(DurationErrorDetail, data.SmallIntMax, nt.DurationMin))
	} else if nt.Start.IsZero() {
		return invalidTraining(fmt.Sprintf(StartErrorDetail, nt.Start))
	} else if len(nt.Sets) == 0 {
		return invalidTraining(SetsErrorDetail)
	}

	uniqueSetOrders := make(map[int]bool)
	for _, s := range nt.Sets {
		if uniqueSetOrders[s.SetOrder] {
			return nonUniqueSetOrder(s.SetOrder)
		}
		uniqueSetOrders[s.SetOrder] = true

		err := validateNewSet(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateNewSet(set apidef.NewTrainingSet) *ApiError {
	if set.SetOrder < 0 || set.SetOrder > data.SmallIntMax {
		return invalidSet(
			set.SetOrder,
			fmt.Sprintf(SetOrderErrorDetail, data.SmallIntMax, set.SetOrder),
		)
	} else if set.Repeat <= 0 || set.Repeat > data.SmallIntMax {
		return invalidSet(
			set.SetOrder,
			fmt.Sprintf(RepeatErrorDetail, data.SmallIntMax, set.Repeat),
		)
	} else if set.DistanceMeters <= 0 || set.DistanceMeters > data.SmallIntMax {
		return invalidSet(
			set.SetOrder,
			fmt.Sprintf(DistanceErrorDetail, data.SmallIntMax, set.DistanceMeters),
		)
	} else if set.StartType != nil && !StartTypesSet[*set.StartType] {
		return invalidSet(
			set.SetOrder,
			fmt.Sprintf(StartTypeUnknownErrorDetail, apidef.Interval, apidef.Pause, *set.StartType),
		)
	} else if set.StartType != nil && set.StartSeconds == nil {
		return invalidSet(set.SetOrder, StartSecondsRequiredErrorDetail)
	} else if set.StartType != nil && set.StartSeconds != nil && (*set.StartSeconds <= 0 || *set.StartSeconds > data.SmallIntMax) {
		return invalidSet(set.SetOrder, fmt.Sprintf(StartSecondsErrorDetail, data.SmallIntMax, *set.StartSeconds))
	} else if set.Equipment != nil && len(*set.Equipment) != 0 {
		for _, eq := range *set.Equipment {
			if !EquipmentSet[eq] {
				return invalidSet(set.SetOrder, fmt.Sprintf(EquipmentErrorDetail, eq))
			}
		}
	} else if set.Group != nil && !GroupSet[*set.Group] {
		return invalidSet(set.SetOrder, fmt.Sprintf(GroupErrorDetail, *set.Group))
	}

	return nil
}

func invalidTraining(detail string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Title:  InvalidTrainingErrorTitle,
			Code:   InvalidTrainingErrorCode,
			Detail: detail,
			Status: http.StatusBadRequest,
		},
	}
}

func invalidSet(setOrder int, detail string) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Title:                InvalidSetErrorTitle,
			Code:                 InvalidSetErrorCode,
			Detail:               detail,
			Status:               http.StatusBadRequest,
			AdditionalProperties: map[string]interface{}{"setOrder": setOrder},
		},
	}
}

func nonUniqueSetOrder(setOrder int) *ApiError {
	return &ApiError{
		ErrorDetail: apidef.ErrorDetail{
			Title:  InvalidSetErrorTitle,
			Code:   NonUniqueSetOrderCode,
			Detail: fmt.Sprintf(NonUniqueSetOrderDetail, setOrder),
			Status: http.StatusBadRequest,
		},
	}
}
