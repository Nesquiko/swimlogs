package app

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nesquiko/swimlogs/pkg/api"
)

func Test_validateNewTrainingNonUniqueSetOrder(t *testing.T) {
	nonUniqueSetOrder := 420
	newTraining := api.NewTraining{
		DurationMinutes: 1,
		Start:           time.Now(),
		Sets: []api.NewTrainingSet{
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
