package model

import (
	"github.com/google/uuid"
)
type CommunityDetail struct {
	CommunityID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryID uuid.UUID `gorm:"type:uuid;primaryKey"`
}