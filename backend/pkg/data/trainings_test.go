package data

import (
	"testing"

	"github.com/Nesquiko/swimlogs/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func Test_setTotalDistances(t *testing.T) {
	training := openapi.NewTraining{
		Blocks: []openapi.NewBlock{
			{
				Repeat: 1,
				Sets: []openapi.NewTrainingSet{
					{Repeat: 1, Distance: 100},
					{Repeat: 1, Distance: 200},
					{Repeat: 2, Distance: 300},
				},
			},
			{
				Repeat: 4,
				Sets: []openapi.NewTrainingSet{
					{Repeat: 10, Distance: 100},
					{Repeat: 1, Distance: 200},
					{Repeat: 2, Distance: 300},
				},
			},
		},
	}

	setTotalDistances(&training)

	expectFirstBlock := 100 + 200 + 2*300
	expectSecondBlock := 4 * (10*100 + 200 + 2*300)

	assert.Equal(t, expectFirstBlock, training.Blocks[0].TotalDistance)
	assert.Equal(t, expectSecondBlock, training.Blocks[1].TotalDistance)
	assert.Equal(t, expectFirstBlock+expectSecondBlock, training.TotalDistance)
}
