package model

import (
	"time"

	"github.com/google/uuid"
)


type UserSecurityAnswer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Question   string    `gorm:"not null;index"`
	AnswerHash string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}