package app

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

const (
	InvalidTrainingErrorTitle = "Invalid training"
	InvalidTrainingErrorCode  = "invalid.training"
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
	UnknownEnumErrorDetail          = "Unknown %s %q"
	NoComponentOrdersErrorDetail    = "No order in set's component"
	StyleIdDoesntExistErrorDetail   = "Style id %q does not exist"

	InvalidSetComponentErrorTitle = "Set component is not valid"
	InvalidSetComponentErrorCode  = "invalid.set.component"
	NoComponentsErrorDetail       = "Set type %q must have components"

	NonUniqueSetOrderCode   = "invalid.set.nonunique.setorder"
	NonUniqueSetOrderDetail = "Multiple sets had same set order %d"

	Equipment   = "Equipment"
	Group       = "Group"
	Type        = "Type"
	Intensity   = "Intensity"
	Progression = "Progression"
)

var (
	EquipmentSet = map[api.EquipmentEnum]bool{
		api.Board:   true,
		api.Fins:    true,
		api.Monofin: true,
		api.Paddles: true,
		api.Snorkel: true,
	}
	GroupSet = map[api.GroupEnum]bool{
		api.Bifi:   true,
		api.Long:   true,
		api.Middle: true,
		api.Mono:   true,
		api.Sprint: true,
	}
	TypeSet = map[api.TypeEnum]bool{
		api.Normal:   true,
		api.Compound: true,
		api.Pyramid:  true,
	}
	ProgressionSet = map[api.ProgressionEnum]bool{
		api.Asc:  true,
		api.Desc: true,
	}
	IntensitySet = map[api.IntensityEnum]bool{
		api.Rec: true,
		api.En1: true,
		api.En2: true,
		api.En3: true,
		api.Sp1: true,
		api.Sp2: true,
		api.Sp3: true,
	}
)

type ValidationError struct {
	api.ErrorDetail
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
			ErrorDetail: api.ErrorDetail{
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

func validateEditSessionRequest(req api.EditSessionRequest) *ValidationError {
	if allNilFields(req) {
		return invalidSession(NoSessionChangesDetail)
	}

	if req.DurationMinutes != nil &&
		(*req.DurationMinutes <= 0 || *req.DurationMinutes > data.SmallIntMax) {
		return invalidSession(
			fmt.Sprintf(DurationErrorDetail, data.SmallIntMax, *req.DurationMinutes),
		)
	}

	if req.Start != nil && req.Start.IsZero() {
		return invalidSession(fmt.Sprintf(StartErrorDetail, req.Start))
	}

	return nil
}

func validateNewTraining(nt api.NewTraining) *ValidationError {
	if nt.Start.IsZero() {
		return invalidTraining(fmt.Sprintf(StartErrorDetail, nt.Start))
	}

	uniqueSetOrders := make(map[int]bool)
	for _, s := range nt.Sets {
		if uniqueSetOrders[s.SetOrder] {
			return nonUniqueSetOrder(s.SetOrder)
		}
		uniqueSetOrders[s.SetOrder] = true
	}

	return nil
}

func checkStyleIdsExist(idChecks []data.IdCheck) *ValidationError {
	for _, check := range idChecks {
		if check.Exists {
			continue
		}

		err := nonExistendStyleId(check.Id.String())
		err.AdditionalProperties = map[string]interface{}{"styleId": check.Id.String()}
		return err
	}
	return nil
}

const NoSetChangesDetail = "Request contained no changes to set."

// func validateEditSetRequest(set api.EditSetRequest) *ValidationError {
// 	if allNilFields(set) {
// 		return invalidEditSet(NoSetChangesDetail)
// 	}
//
// 	if set.Repeat != nil && (*set.Repeat < 1 || *set.Repeat > data.SmallIntMax) {
// 		return invalidEditSet(fmt.Sprintf(RepeatErrorDetail, data.SmallIntMax, *set.Repeat))
// 	} else if set.DistanceMeters != nil && (*set.DistanceMeters <= 0 || *set.DistanceMeters > data.SmallIntMax) {
// 		return invalidEditSet(fmt.Sprintf(DistanceErrorDetail, data.SmallIntMax, *set.DistanceMeters))
// 	} else if set.StartType != nil && !StartTypesSet[*set.StartType] {
// 		return invalidEditSet(fmt.Sprintf(StartTypeUnknownErrorDetail, api.Interval, api.Pause, *set.StartType))
// 	} else if set.StartType != nil && set.StartSeconds == nil {
// 		return invalidEditSet(StartSecondsRequiredErrorDetail)
// 	} else if set.StartType != nil && set.StartSeconds != nil && (*set.StartSeconds <= 0 || *set.StartSeconds > data.SmallIntMax) {
// 		return invalidEditSet(fmt.Sprintf(StartSecondsErrorDetail, data.SmallIntMax, *set.StartSeconds))
// 	} else if set.Equipment != nil && len(*set.Equipment) != 0 {
// 		for _, eq := range *set.Equipment {
// 			if !EquipmentSet[eq] {
// 				return invalidEditSet(fmt.Sprintf(UnknownEnumErrorDetail, Equipment, eq))
// 			}
// 		}
// 	} else if set.Group != nil && !GroupSet[*set.Group] {
// 		return invalidEditSet(fmt.Sprintf(UnknownEnumErrorDetail, Group, *set.Group))
// 	}
//
// 	return nil
// }

func invalidSession(detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: api.ErrorDetail{
			Title:  InvalidEditSessionRequestTitle,
			Code:   InvalidEditSessionRequestCode,
			Detail: detail,
			Status: http.StatusBadRequest,
		},
	}
}

func invalidTraining(detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: api.ErrorDetail{
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
		ErrorDetail: api.ErrorDetail{
			Title:                InvalidSetErrorTitle,
			Code:                 InvalidSetErrorCode,
			Detail:               detail,
			Status:               http.StatusBadRequest,
			AdditionalProperties: map[string]interface{}{"setOrder": setOrder},
		},
	}
}

func invalidSetComponent(detail string) *ValidationError {
	return &ValidationError{
		ErrorDetail: api.ErrorDetail{
			Title:  InvalidSetComponentErrorTitle,
			Code:   InvalidSetComponentErrorCode,
			Detail: detail,
			Status: http.StatusBadRequest,
		},
	}
}

const (
	NonExistentStyleIdCode   = "non.existent.styleid"
	NonExistentStyleIdTitle  = "Style with provided id doesn't exits"
	NonExistentStyleIdDetail = "Style with id %q doesn't exits"
)

func nonExistendStyleId(styleId string) *ValidationError {
	return &ValidationError{
		ErrorDetail: api.ErrorDetail{
			Title:  NonExistentStyleIdTitle,
			Code:   NonExistentStyleIdCode,
			Detail: fmt.Sprintf(NonExistentStyleIdDetail, styleId),
			Status: http.StatusBadRequest,
		},
	}
}

func nonUniqueSetOrder(setOrder int) *ValidationError {
	return &ValidationError{
		ErrorDetail: api.ErrorDetail{
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
