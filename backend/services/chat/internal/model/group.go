package model

import "github.com/google/uuid"

type Group struct {
	GroupID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsPrivate bool `gorm:"not null;default:false"`
}