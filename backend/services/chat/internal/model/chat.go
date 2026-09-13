package model

import (
	"time"

	"github.com/google/uuid"
)
type Chat struct{
	ChatID uuid.UUID `gorm:"type:uuid;primaryKey"`
	GroupID uuid.UUID `gorm:"type:uuid"`
	SenderID uuid.UUID `gorm:"type:uuid"`
	Message string `gorm:"type:text"`
	IsMasked bool `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"notnull"`
}