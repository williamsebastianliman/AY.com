package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadBookmark struct {
    ThreadID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
    UserID   uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`

    CreatedAt time.Time `gorm:"autoCreateTime"`
}