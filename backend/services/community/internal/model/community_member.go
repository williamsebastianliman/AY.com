package model

import (
	"time"

	"github.com/google/uuid"
)

type CommunityMember struct {
	CommunityID uuid.UUID `gorm:"type:uuid;primaryKey"`
	MemberID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status string `gorm:"type:text"` 
	CreatedAt time.Time `gorm:"notnull"`

	//role: pending accepted moderator owner
}