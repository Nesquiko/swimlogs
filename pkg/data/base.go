package data

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	Id         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func createBase() *Base {
	return &Base{
		Id:         uuid.New(),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}
