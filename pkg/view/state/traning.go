package state

import (
	"errors"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/google/uuid"
)

type TrainingStateStorage struct {
	t       *oapiGen.Training
	details []oapiGen.TrainingDetail
}

func (tss *TrainingStateStorage) GetDetails() ([]oapiGen.TrainingDetail, error) {
	if tss.details == nil || len(tss.details) == 0 {
		return nil, errors.New("no details stored")
	}
	return tss.details, nil
}

func (tss *TrainingStateStorage) SaveDetails(details []oapiGen.TrainingDetail) {
	tss.details = details
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
