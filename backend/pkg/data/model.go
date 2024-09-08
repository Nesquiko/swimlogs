package data

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Training struct {
	// DB rows
	Id          uuid.UUID
	Start       time.Time
	DurationMin int
	CreatedAt   time.Time
	ModifiedAt  time.Time

	Sets []TrainingSet
	// Not persisted in DB, calculated during query or in mapping from sets
	TotalDistance int
}

type TrainingSet struct {
	// DB rows
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
	Intensity      *string
	Progression    *string

	Style      *Style
	Components *[]SetComponent
}

type emptyTrainingSet struct {
	Id             *uuid.UUID
	TrainingId     *uuid.UUID
	SetOrder       *int
	Repeat         *int
	DistanceMeters *int
	Description    *string
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	Group          *string
	IsMain         *bool
	SetType        *string
	StyleId        *uuid.UUID
	Intensity      *string
	Progression    *string

	Style emptyStyle
}

func (s emptyTrainingSet) intoTrainingSet() TrainingSet {
	if s.Id == nil {
		slog.Error("emptyTrainingSet.intoTrainingSet: called without checking if Id is nil")
		return TrainingSet{}
	}

	set := TrainingSet{
		Id:             *s.Id,
		TrainingId:     *s.TrainingId,
		SetOrder:       *s.SetOrder,
		Repeat:         *s.Repeat,
		DistanceMeters: *s.DistanceMeters,
		Description:    s.Description,
		Equipment:      s.Equipment,
		StartType:      s.StartType,
		StartSeconds:   s.StartSeconds,
		Group:          s.Group,
		IsMain:         *s.IsMain,
		SetType:        *s.SetType,
		Intensity:      s.Intensity,
		Progression:    s.Progression,
	}

	if s.StyleId != nil {
		set.Style = &Style{
			Id:                 *s.Style.Id,
			Name:               *s.Style.Name,
			ExeciseName:        s.Style.ExeciseName,
			ExeciseDescription: s.Style.ExeciseDescription,
		}
		set.StyleId = s.Style.Id
	}

	return set
}

type Style struct {
	Id                 uuid.UUID
	Name               string
	ExeciseName        *string
	ExeciseDescription *string
}

type emptyStyle struct {
	Id                 *uuid.UUID
	Name               *string
	ExeciseName        *string
	ExeciseDescription *string
}

type SetComponent struct {
	// DB rows
	Id             uuid.UUID
	SetId          uuid.UUID
	Orders         []int
	IterationOrder *int
	Repeat         int
	DistanceMeters int
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	StyleId        *uuid.UUID
	Group          *string
	Intensity      *string
	Progression    *string
	Description    *string

	Style *Style
}

type emptySetComponent struct {
	Id             *uuid.UUID
	SetId          *uuid.UUID
	Orders         []int
	IterationOrder *int
	Repeat         *int
	DistanceMeters *int
	Equipment      *[]string
	StartType      *string
	StartSeconds   *int
	StyleId        *uuid.UUID
	Group          *string
	Intensity      *string
	Progression    *string
	Description    *string
}
