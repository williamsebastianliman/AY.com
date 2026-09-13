package model

import (
	"github.com/google/uuid"
)

type UserFollower struct{
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	FollowerID uuid.UUID `gorm:"type:uuid;primaryKey"`
}