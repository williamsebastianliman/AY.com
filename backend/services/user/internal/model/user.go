package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name             string         `gorm:"type:text;not null;check:char_length(name) > 4"`
    Username         string         `gorm:"type:text;not null;uniqueIndex"`
    Email            string         `gorm:"type:text;not null;uniqueIndex"`
    PasswordHash     string         `gorm:"type:text;not null"`
    Gender           string         `gorm:"type:text;not null"`
    DateOfBirth      time.Time      `gorm:"type:date;not null"`
    IsActivated      bool           `gorm:"not null;default:false"`
    ActivatedAt      *time.Time     `gorm:"type:timestamptz"`
    IsBanned         bool           `gorm:"not null;default:false"`
    IsPremium        bool           `gorm:"not null;default:false"`
    BannedAt         *time.Time     `gorm:"type:timestamptz"`
    BannedReason     *string        `gorm:"type:text"`
    SubscribedNews   bool           `gorm:"not null;default:false"`
    CreatedAt        time.Time      `gorm:"autoCreateTime"`
    UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
    DeletedAt        gorm.DeletedAt `gorm:"index"`
    Role string `gorm:"type:text"`

    ProfilePictureID *uuid.UUID `gorm:"type:uuid;index"`
    BannerMediaID    *uuid.UUID `gorm:"type:uuid;index"`
}