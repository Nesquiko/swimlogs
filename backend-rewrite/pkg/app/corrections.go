package app

import (
	"github.com/Nesquiko/swimlogs/apidef"
)

func recalcDistanceOnNewTraining(nt *apidef.NewTraining) {
	total := 0
	for i := 0; i < len(nt.Sets); i++ {
		ns := &nt.Sets[i]
		ns.TotalDistance = ns.Repeat * ns.DistanceMeters
		total += ns.TotalDistance
	}
	nt.TotalDistance = total
}

func recalcDistanceOnTraining(t *apidef.Training) {
	total := 0
	for i := 0; i < len(t.Sets); i++ {
		ns := &t.Sets[i]
		ns.TotalDistance = ns.Repeat * ns.DistanceMeters
		total += t.Sets[i].TotalDistance
	}
	t.TotalDistance = total
}
