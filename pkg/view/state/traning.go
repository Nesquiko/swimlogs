package state

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
	"github.com/Nesquiko/swimlogs/pkg/view/api"
	"github.com/google/uuid"
)

const PageSize = 10

type TrainingStateStorage struct {
	t                  *oapiGen.Training
	currentWeekDetails []oapiGen.TrainingDetail

	details        []oapiGen.TrainingDetail
	lastPagination oapiGen.Pagination
}

func (tss *TrainingStateStorage) TrainingDetails() (*[]oapiGen.TrainingDetail, error) {
	if tss.details == nil || len(tss.details) == 0 {
		err := tss.LoadTrainingDetails()
		if err != nil {
			return nil, fmt.Errorf("TrainingDetails: %w", err)
		}
	}

	sort.Slice(tss.details, func(i, j int) bool {
		return tss.details[i].Date.After(tss.details[j].Date.Time)
	})

	return &tss.details, nil
}

func (tss *TrainingStateStorage) LoadTrainingDetails() error {
	fetched := 0
	for i := 0; i < 3; i++ {
		details, pagination, err := api.GetTrainingDetails(i, PageSize)
		if err != nil {
			return fmt.Errorf("LoadTrainingDetails: %w", err)
		}
		tss.details = append(tss.details, details...)
		fetched += len(details)
		tss.lastPagination = pagination
		if fetched >= pagination.Total {
			break
		}
	}

	return nil
}

func (tss *TrainingStateStorage) LoadNextDetailsPage() error {
	if tss.lastPagination.Page*PageSize >= tss.lastPagination.Total {
		return nil
	}
	details, pagination, err := api.GetTrainingDetails(tss.lastPagination.Page+1, PageSize)
	if err != nil {
		return fmt.Errorf("LoadNextDetailsPage: %w", err)
	}
	tss.details = append(tss.details, details...)
	sort.Slice(tss.details, func(i, j int) bool {
		return tss.details[i].Date.After(tss.details[j].Date.Time)
	})
	tss.lastPagination = pagination
	return nil
}

func (tss *TrainingStateStorage) GetDetails() ([]oapiGen.TrainingDetail, error) {
	if tss.currentWeekDetails == nil || len(tss.currentWeekDetails) == 0 {
		return nil, errors.New("no details stored")
	}
	return tss.currentWeekDetails, nil
}

func (tss *TrainingStateStorage) SaveDetails(details []oapiGen.TrainingDetail) {
	tss.currentWeekDetails = details
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
