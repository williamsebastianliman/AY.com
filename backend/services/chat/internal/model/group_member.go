package model

import "github.com/google/uuid"

type GroupMember struct {
	GroupID uuid.UUID `gorm:"type:uuid;primaryKey"`
	MemberID uuid.UUID `gorm:"type:uuid;primaryKey"`
}