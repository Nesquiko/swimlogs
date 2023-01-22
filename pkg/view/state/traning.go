package state

import (
	"errors"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

type TrainingStateStorage struct {
	t *oapiGen.Training
}

func (tss *TrainingStateStorage) IsInCache(id uuid.UUID) bool {
	if tss.t == nil {
		return false
	}
	return tss.t.Id == id
}

func (tss *TrainingStateStorage) SaveTraining(t oapiGen.Training) {
	tss.t = &t
}

func (tss *TrainingStateStorage) GetTraining() (oapiGen.Training, error) {
	if tss.t == nil {
		return oapiGen.Training{}, errors.New("no training stored")
	}
	return *tss.t, nil
}

func (tss *TrainingStateStorage) InvalidateTrainingCache() {
	tss.t = nil
}

type TrainingStateStorageSetter interface {
	TrainingStateStorageSet(tss *TrainingStateStorage)
}

// TrainingStateStorage ref for wiring
type TrainingStateStorageRef struct {
	*TrainingStateStorage
}

func (tssr *TrainingStateStorageRef) TrainingStateStorageSet(tss *TrainingStateStorage) {
	tssr.TrainingStateStorage = tss
}
