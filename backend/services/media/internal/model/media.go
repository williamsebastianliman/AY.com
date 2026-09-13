package model

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    PublicURL string    `gorm:"type:text;not null"`
    Extension string `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}