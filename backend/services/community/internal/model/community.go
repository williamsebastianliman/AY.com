package model

import (
	"time"

	"github.com/google/uuid"
)
type Community struct{
	CommunityID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CommunityName string `gorm:"type:text"`
	CommunityRules string `gorm:"type:text"`
	CommunityDescription string `gorm:"type:text"`
	Status string `gorm:"type:text"`
	CreatedAt time.Time `gorm:"notnull"`

	CommunityLogo string `gorm:"type:text"`
	CommunityBanner string `gorm:"type:text"`

	//pending, activated
}