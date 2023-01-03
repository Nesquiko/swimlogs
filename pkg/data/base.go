package data

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
