package model

import (
	"github.com/google/uuid"
)

type UserBlock struct{
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	BlockedUserID uuid.UUID `gorm:"type:uuid;primaryKey"`
}