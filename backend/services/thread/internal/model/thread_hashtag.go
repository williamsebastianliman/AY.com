package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadHashtag struct {
    ThreadID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
    HashtagName string `gorm:"type:text;not null;primaryKey"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}