package app

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

func Test_validateEditSessionRequestInvalidStart(t *testing.T) {
	editSession := apidef.EditSessionRequest{Start: &time.Time{}}

	err := validateEditSessionRequest(editSession)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidEditSessionRequestTitle, err.Title)
	assert.Equal(InvalidEditSessionRequestCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(fmt.Sprintf(StartErrorDetail, time.Time{}), err.Detail)
	assert.Nil(err.AdditionalProperties)
}

func Test_validateEditSessionRequestInvalidDuration(t *testing.T) {
	testCases := []struct {
		desc            string
		invalidDuration int
	}{
		{desc: "Lower bound", invalidDuration: 0},
		{desc: "Upper bound", invalidDuration: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			editSession := apidef.EditSessionRequest{DurationMin: &tC.invalidDuration}

			err := validateEditSessionRequest(editSession)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidEditSessionRequestTitle, err.Title)
			assert.Equal(InvalidEditSessionRequestCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(DurationErrorDetail, data.SmallIntMax, tC.invalidDuration),
				err.Detail,
			)
			assert.Nil(err.AdditionalProperties)
		})
	}
}

func Test_validateNewTrainingInvalidDuration(t *testing.T) {
	testCases := []struct {
		desc            string
		invalidDuration int
	}{
		{desc: "Lower bound", invalidDuration: 0},
		{desc: "Upper bound", invalidDuration: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newTraining := apidef.NewTraining{
				DurationMin: tC.invalidDuration,
			}

			err := validateNewTraining(newTraining)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidTrainingErrorTitle, err.Title)
			assert.Equal(InvalidTrainingErrorCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(DurationErrorDetail, data.SmallIntMax, tC.invalidDuration),
				err.Detail,
			)
			assert.Nil(err.AdditionalProperties)
		})
	}
}

func Test_validateNewTrainingInvalidStart(t *testing.T) {
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Time{},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidTrainingErrorTitle, err.Title)
	assert.Equal(InvalidTrainingErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(fmt.Sprintf(StartErrorDetail, time.Time{}), err.Detail)
	assert.Nil(err.AdditionalProperties)
}

func Test_validateNewTrainingNoSets(t *testing.T) {
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidTrainingErrorTitle, err.Title)
	assert.Equal(InvalidTrainingErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(SetsErrorDetail, err.Detail)
	assert.Nil(err.AdditionalProperties)
}

func Test_validateNewTrainingNonUniqueSetOrder(t *testing.T) {
	nonUniqueSetOrder := 420
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
		Sets: []apidef.NewTrainingSet{
			{SetOrder: nonUniqueSetOrder, Repeat: 1, DistanceMeters: 1},
			{SetOrder: nonUniqueSetOrder},
		},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidSetErrorTitle, err.Title)
	assert.Equal(NonUniqueSetOrderCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(fmt.Sprintf(NonUniqueSetOrderDetail, nonUniqueSetOrder), err.Detail)
	assert.Nil(err.AdditionalProperties)
}

func Test_validateNewTrainingInvalidSetOrder(t *testing.T) {
	testCases := []struct {
		desc            string
		invalidSetOrder int
	}{
		{desc: "Lower bound", invalidSetOrder: -1},
		{desc: "Upper bound", invalidSetOrder: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newTraining := apidef.NewTraining{
				DurationMin: 1,
				Start:       time.Now(),
				Sets: []apidef.NewTrainingSet{
					{SetOrder: tC.invalidSetOrder},
				},
			}

			err := validateNewTraining(newTraining)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidSetErrorTitle, err.Title)
			assert.Equal(InvalidSetErrorCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(SetOrderErrorDetail, data.SmallIntMax, tC.invalidSetOrder),
				err.Detail,
			)
			assert.NotNil(err.AdditionalProperties)
			assert.Equal(tC.invalidSetOrder, err.AdditionalProperties["setOrder"])
		})
	}
}

func Test_validateNewTrainingSetInvalidRepeat(t *testing.T) {
	testCases := []struct {
		desc          string
		invalidRepeat int
	}{
		{desc: "Lower bound", invalidRepeat: 0},
		{desc: "Upper bound", invalidRepeat: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newTraining := apidef.NewTraining{
				DurationMin: 1,
				Start:       time.Now(),
				Sets: []apidef.NewTrainingSet{
					{SetOrder: 0, Repeat: tC.invalidRepeat},
				},
			}

			err := validateNewTraining(newTraining)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidSetErrorTitle, err.Title)
			assert.Equal(InvalidSetErrorCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(RepeatErrorDetail, data.SmallIntMax, tC.invalidRepeat),
				err.Detail,
			)
			assert.NotNil(err.AdditionalProperties)
			assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
		})
	}
}

func Test_validateNewTrainingSetInvalidDistance(t *testing.T) {
	testCases := []struct {
		desc            string
		invalidDistance int
	}{
		{desc: "Lower bound", invalidDistance: 0},
		{desc: "Upper bound", invalidDistance: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newTraining := apidef.NewTraining{
				DurationMin: 1,
				Start:       time.Now(),
				Sets: []apidef.NewTrainingSet{
					{SetOrder: 0, Repeat: 1, DistanceMeters: tC.invalidDistance},
				},
			}

			err := validateNewTraining(newTraining)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidSetErrorTitle, err.Title)
			assert.Equal(InvalidSetErrorCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(DistanceErrorDetail, data.SmallIntMax, tC.invalidDistance),
				err.Detail,
			)
			assert.NotNil(err.AdditionalProperties)
			assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
		})
	}
}

func Test_validateNewTrainingSetInvalidStartType(t *testing.T) {
	invalidStartType := apidef.StartTypeEnum("invalid")
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
		Sets: []apidef.NewTrainingSet{
			{SetOrder: 0, Repeat: 1, DistanceMeters: 1, StartType: &invalidStartType},
		},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidSetErrorTitle, err.Title)
	assert.Equal(InvalidSetErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(
		fmt.Sprintf(StartTypeUnknownErrorDetail, apidef.Interval, apidef.Pause, invalidStartType),
		err.Detail,
	)
	assert.NotNil(err.AdditionalProperties)
	assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
}

func Test_validateNewTrainingSetNoStartSeconds(t *testing.T) {
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
		Sets: []apidef.NewTrainingSet{
			{SetOrder: 0, Repeat: 1, DistanceMeters: 1, StartType: asPtr(apidef.Interval)},
		},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidSetErrorTitle, err.Title)
	assert.Equal(InvalidSetErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(StartSecondsRequiredErrorDetail, err.Detail)
	assert.NotNil(err.AdditionalProperties)
	assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
}

func Test_validateNewTrainingSetInvalidStartSeconds(t *testing.T) {
	testCases := []struct {
		desc                string
		invalidStartSeconds int
	}{
		{desc: "Lower bound", invalidStartSeconds: 0},
		{desc: "Upper bound", invalidStartSeconds: data.SmallIntMax + 1},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			newTraining := apidef.NewTraining{
				DurationMin: 1,
				Start:       time.Now(),
				Sets: []apidef.NewTrainingSet{
					{SetOrder: 0, Repeat: 1, DistanceMeters: 1},
					{
						SetOrder:       1,
						Repeat:         1,
						DistanceMeters: 1,
						StartType:      asPtr(apidef.Interval),
						StartSeconds:   &tC.invalidStartSeconds,
					},
				},
			}

			err := validateNewTraining(newTraining)

			assert := assert.New(t)
			assert.NotNil(err)
			assert.Equal(InvalidSetErrorTitle, err.Title)
			assert.Equal(InvalidSetErrorCode, err.Code)
			assert.Equal(http.StatusBadRequest, err.Status)
			assert.Equal(
				fmt.Sprintf(StartSecondsErrorDetail, data.SmallIntMax, tC.invalidStartSeconds),
				err.Detail,
			)
			assert.NotNil(err.AdditionalProperties)
			assert.Equal(newTraining.Sets[1].SetOrder, err.AdditionalProperties["setOrder"])
		})
	}
}

func Test_validateNewTrainingUnknownEquipment(t *testing.T) {
	unknownEquipment := apidef.EquipmentEnum("unknownEquipment")
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
		Sets: []apidef.NewTrainingSet{
			{
				SetOrder:       0,
				Repeat:         1,
				DistanceMeters: 1,
				Equipment: &[]apidef.EquipmentEnum{
					apidef.Fins,
					apidef.Board,
					unknownEquipment,
				},
			},
		},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidSetErrorTitle, err.Title)
	assert.Equal(InvalidSetErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(fmt.Sprintf(EquipmentErrorDetail, unknownEquipment), err.Detail)
	assert.NotNil(err.AdditionalProperties)
	assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
}

func Test_validateNewTrainingUnknownGroup(t *testing.T) {
	unknownGroup := apidef.GroupEnum("unknownGroup")
	newTraining := apidef.NewTraining{
		DurationMin: 1,
		Start:       time.Now(),
		Sets: []apidef.NewTrainingSet{
			{
				SetOrder:       0,
				Repeat:         1,
				DistanceMeters: 1,
				Group:          &unknownGroup,
			},
		},
	}

	err := validateNewTraining(newTraining)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.Equal(InvalidSetErrorTitle, err.Title)
	assert.Equal(InvalidSetErrorCode, err.Code)
	assert.Equal(http.StatusBadRequest, err.Status)
	assert.Equal(fmt.Sprintf(GroupErrorDetail, unknownGroup), err.Detail)
	assert.NotNil(err.AdditionalProperties)
	assert.Equal(newTraining.Sets[0].SetOrder, err.AdditionalProperties["setOrder"])
}

func asPtr[T any](v T) *T {
	return &v
}
