package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadDetail struct {
    ThreadID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
    CategoryID   uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`

    CreatedAt time.Time `gorm:"autoCreateTime"`
}