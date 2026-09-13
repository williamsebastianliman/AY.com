package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadMedia struct {
    ThreadID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
    MediaID  uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`

    CreatedAt time.Time `gorm:"autoCreateTime"`
}