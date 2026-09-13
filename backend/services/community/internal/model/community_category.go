package model

import (
	"time"

	"github.com/google/uuid"
)

type CommunityCategory struct {
	CategoryID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryName string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"notnull"`
}