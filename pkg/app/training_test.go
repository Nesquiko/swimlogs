package app

import (
	"testing"

	"github.com/Nesquiko/swimlogs/generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/data"
)

var zero = 0

func Test_updateTotalDist(t *testing.T) {
	dataT := data.Training{
		TotalDistance: 1500,
		Blocks: []data.Block{
			{
				Repeat: 3,
				Sets: []data.Set{
					{TotalDistance: 400},
					{TotalDistance: 100},
				},
				TotalDistance: 500,
			},
			{
				Repeat: 2,
				Sets: []data.Set{
					{TotalDistance: 400},
				},
				TotalDistance: 800,
			},
		},
	}

	training := oapiGen.Training{
		TotalDist: &zero,
		Blocks: []oapiGen.Block{
			{
				Sets: []oapiGen.Set{
					{TotalDist: &zero},
					{TotalDist: &zero},
				},
				TotalDist: &zero,
			},
			{
				Sets: []oapiGen.Set{
					{TotalDist: &zero},
				},
				TotalDist: &zero,
			},
		},
	}

	updateTotalDist(&training, dataT)

	if *training.TotalDist != dataT.TotalDistance {
		t.Errorf(
			"training total distance, expected %d, but was %d",
			dataT.TotalDistance,
			*training.TotalDist,
		)
	}
	if *training.Blocks[0].TotalDist != dataT.Blocks[0].TotalDistance {
		t.Errorf(
			"first block total distance, expected %d, but was %d",
			dataT.Blocks[0].TotalDistance,
			*training.Blocks[0].TotalDist,
		)
	}
	if *training.Blocks[1].TotalDist != dataT.Blocks[1].TotalDistance {
		t.Errorf(
			"second block total distance, expected %d, but was %d",
			dataT.Blocks[1].TotalDistance,
			*training.Blocks[1].TotalDist,
		)
	}
	if *training.Blocks[0].Sets[0].TotalDist != dataT.Blocks[0].Sets[0].TotalDistance {
		t.Errorf(
			"first set first block total distance, expected %d, but was %d",
			dataT.Blocks[0].Sets[0].TotalDistance,
			*training.Blocks[0].Sets[0].TotalDist,
		)
	}
	if *training.Blocks[0].Sets[1].TotalDist != dataT.Blocks[0].Sets[1].TotalDistance {
		t.Errorf(
			"second set of first block total distance, expected %d, but was %d",
			dataT.Blocks[0].Sets[1].TotalDistance,
			*training.Blocks[0].Sets[1].TotalDist,
		)
	}
	if *training.Blocks[1].Sets[0].TotalDist != dataT.Blocks[1].Sets[0].TotalDistance {
		t.Errorf(
			"first set second block total distance, expected %d, but was %d",
			dataT.Blocks[1].Sets[0].TotalDistance,
			*training.Blocks[1].Sets[0].TotalDist,
		)
	}
}
