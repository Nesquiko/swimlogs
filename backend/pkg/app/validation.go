package app

import (
	"fmt"
	"net/http"
	"reflect"

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

type ValidationError struct {
	apidef.ErrorDetail
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("error %q, status %d, detail %q", e.Title, e.Status, e.Detail)
}

const (
	InvalidNewOrderSetCode   = "invalid.setorder"
	InvalidNewOrderSetTitle  = "New set order is invalid"
	InvalidNewOrderSetDetail = "Set order must be between 0 and %d, was %d"
)

func validateNewSetOrder(newSetOrder int) *ValidationError {
	if newSetOrder < 0 || newSetOrder > data.SmallIntMax {
		return &ValidationError{
			ErrorDetail: apidef.ErrorDetail{
				Title:  InvalidNewOrderSetTitle,
				Code:   InvalidNewOrderSetCode,
				Detail: fmt.Sprintf(InvalidNewOrderSetDetail, data.SmallIntMax, newSetOrder),
				Status: http.StatusBadRequest,
			},
		}
	}
	return nil
}

const (
	InvalidEditSessionRequestCode  = "invalid.session"
	InvalidEditSessionRequestTitle = "Invalid edit session"
	NoSessionChangesDetail         = "Request contained no changes start nor duration changes."
)

func validateEditSessionRequest(req apidef.EditSessionRequest) *ValidationError {
	if allNilFields(req) {
		return invalidSession(NoSessionChangesDetail)
	}

	if req.DurationMin != nil && (*req.DurationMin <= 0 || *req.DurationMin > data.SmallIntMax) {
		return invalidSession(fmt.Sprintf(DurationErrorDetail, data.SmallIntMax, *req.DurationMin))
	}

	if req.Start != nil && req.Start.IsZero() {
		return invalidSession(fmt.Sprintf(StartErrorDetail, req.Start))
	}

	return nil
}

func validateNewTraining(nt apidef.NewTraining) *ValidationError {
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

func validateNewSet(set apidef.NewTrainingSet) *ValidationError {
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

const NoSetChangesDetail = "Request contained no changes to set."

func validateEditSetRequest(set apidef.EditSetRequest) *ValidationError {
	if allNilFields(set) {
		return invalidEditSet(NoSetChangesDetail)
	}

	if set.Repeat != nil && (*set.Repeat < 1 || *set.Repeat > data.SmallIntMax) {
		return invalidEditSet(fmt.Sprintf(RepeatErrorDetail, data.SmallIntMax, *set.Repeat))
	} else if set.DistanceMeters != nil && (*set.DistanceMeters <= 0 || *set.DistanceMeters > data.SmallIntMax) {
		return invalidEditSet(fmt.Sprintf(DistanceErrorDetail, data.SmallIntMax, *set.DistanceMeters))
	} else if set.StartType != nil && !StartTypesSet[*set.StartType] {
		return invalidEditSet(fmt.Sprintf(StartTypeUnknownErrorDetail, apidef.Interval, apidef.Pause, *set.StartType))
	} else if set.StartType != nil && set.StartSeconds == nil {
		return invalidEditSet(StartSecondsRequiredErrorDetail)
	} else if set.StartType != nil && set.StartSeconds != nil && (*set.StartSeconds <= 0 || *set.StartSeconds > data.SmallIntMax) {
		return invalidEditSet(fmt.Sprintf(StartSecondsErrorDetail, data.SmallIntMax, *set.StartSeconds))
	} else if set.Equipment != nil && len(*set.Equipment) != 0 {
		for _, eq := range *set.Equipment {
			if !EquipmentSet[eq] {
				return invalidEditSet(fmt.Sprintf(EquipmentErrorDetail, eq))
			}
		}
	} else if set.Group != nil && !GroupSet[*set.Group] {
		return invalidEditSet(fmt.Sprintf(GroupErrorDetail, *set.Group))
	}

	return nil
}

func invalidSession(detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: apidef.ErrorDetail{
			Title:  InvalidEditSessionRequestTitle,
			Code:   InvalidEditSessionRequestCode,
			Detail: detail,
			Status: http.StatusBadRequest,
		},
	}
}

func invalidTraining(detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: apidef.ErrorDetail{
			Title:  InvalidTrainingErrorTitle,
			Code:   InvalidTrainingErrorCode,
			Detail: detail,
			Status: http.StatusBadRequest,
		},
	}
}

func invalidEditSet(detail string) *ValidationError {
	err := invalidSet(-1, detail)
	err.AdditionalProperties = nil
	return err
}

func invalidSet(setOrder int, detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: apidef.ErrorDetail{
			Title:                InvalidSetErrorTitle,
			Code:                 InvalidSetErrorCode,
			Detail:               detail,
			Status:               http.StatusBadRequest,
			AdditionalProperties: map[string]interface{}{"setOrder": setOrder},
		},
	}
}

func nonUniqueSetOrder(setOrder int) *ValidationError {
	return &ValidationError{
		ErrorDetail: apidef.ErrorDetail{
			Title:  InvalidSetErrorTitle,
			Code:   NonUniqueSetOrderCode,
			Detail: fmt.Sprintf(NonUniqueSetOrderDetail, setOrder),
			Status: http.StatusBadRequest,
		},
	}
}

func allNilFields(s any) bool {
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return false
	}

	structVal := reflect.ValueOf(s)
	for i := 0; i < structVal.NumField(); i++ {
		if f := structVal.Field(i); f.IsValid() && !f.IsNil() {
			return false
		}
	}
	return true
}
