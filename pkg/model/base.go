package model

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
