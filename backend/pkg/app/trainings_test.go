package app

import (
	"testing"
	"time"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func Test_recalculateTotalDistances(t *testing.T) {
	training := openapi.NewTraining{
		Start:         time.Now(),
		DurationMin:   90,
		TotalDistance: 0, // should be 1900
		Sets: []openapi.NewTrainingSet{
			{
				SetOrder:       asPtr(0),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  0, // should be 400
				DistanceMeters: asPtr(400),
				Description:    asPtr("warm up freestyle"),
			},
			{
				SetOrder:       asPtr(1),
				Repeat:         3,
				StartType:      openapi.Pause,
				TotalDistance:  0, // should be 600
				DistanceMeters: asPtr(200),
				Description:    asPtr("drills"),
				StartSeconds:   asPtr(20),
			},
			{
				SetOrder:      asPtr(2),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 0, // should be 300
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  0, // should be 50
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  0, // should be 25
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(3),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  0, // should be 100
				DistanceMeters: asPtr(100),
				Description:    asPtr("cool down"),
			},
			{
				SetOrder:      asPtr(4),
				Repeat:        4,
				StartType:     openapi.None,
				TotalDistance: 0, // should be 300
				SubSets: &[]openapi.NewTrainingSet{
					{
						SubSetOrder:    asPtr(0),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  0, // should be 50
						DistanceMeters: asPtr(50),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(60),
					},
					{
						SubSetOrder:    asPtr(1),
						Repeat:         1,
						StartType:      openapi.Pause,
						TotalDistance:  0, // should be 25
						DistanceMeters: asPtr(25),
						Description:    asPtr("max speed"),
						StartSeconds:   asPtr(45),
					},
				},
			},
			{
				SetOrder:       asPtr(5),
				Repeat:         1,
				StartType:      openapi.None,
				TotalDistance:  0, // should be 200
				DistanceMeters: asPtr(200),
				Description:    asPtr("breastroke cool down"),
			},
		},
	}

	recalculateTotalDistances(&training)

	assert := assert.New(t)
	assert.Equal(1900, training.TotalDistance)
	assert.Equal(400, training.Sets[0].TotalDistance)
	assert.Equal(600, training.Sets[1].TotalDistance)
	assert.Equal(300, training.Sets[2].TotalDistance)
	assert.Equal(50, (*training.Sets[2].SubSets)[0].TotalDistance)
	assert.Equal(25, (*training.Sets[2].SubSets)[1].TotalDistance)
	assert.Equal(100, training.Sets[3].TotalDistance)
	assert.Equal(300, training.Sets[4].TotalDistance)
	assert.Equal(50, (*training.Sets[4].SubSets)[0].TotalDistance)
	assert.Equal(25, (*training.Sets[4].SubSets)[1].TotalDistance)
	assert.Equal(200, training.Sets[5].TotalDistance)
}
