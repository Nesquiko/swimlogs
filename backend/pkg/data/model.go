package data

import (
	"time"

	"github.com/google/uuid"
)

type Training struct {
	Id          uuid.UUID
	Start       time.Time
	DurationMin int
	Sets        []TrainingSet

	TotalDistance int
	CreatedAt     time.Time
	ModifiedAt    time.Time
}

type TrainingSet struct {
	Id             uuid.UUID
	TrainingId     uuid.UUID
	SetOrder       int
	Repeat         int
	DistanceMeters int
	Description    *string
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	Group          *string
	IsMain         bool
	SetType        string
	StyleId        *uuid.UUID
	Style          *Style
	Intensity      *string
	Progression    *string
	Components     *[]SetComponent
}

type Style struct {
	Id                 uuid.UUID
	Name               string
	ExeciseName        *string
	ExeciseDescription *string
}

type SetComponent struct {
	Id             uuid.NullUUID
	SetId          uuid.UUID
	Orders         []int
	IterationOrder *int
	Repeat         int
	DistanceMeters int
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	StyleId        *uuid.UUID
	Style          *Style
	Group          *string
	Intensity      *string
	Progression    *string
	Description    *string
}
