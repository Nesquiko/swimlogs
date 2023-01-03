package data

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;not null;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.Id = uuid.New()
	return nil
}
