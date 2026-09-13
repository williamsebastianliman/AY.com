package model

import (
	"github.com/google/uuid"
)

type PremiumRequest struct{
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CardNumber string `gorm:"type:text"`
	Reason string `gorm:"type:text"`
	FaceImage *uuid.UUID `gorm:"type:uuid;index"`
}