package model

import (
	"time"

	"github.com/google/uuid"
)

type Thread struct {
  ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
  UserID       uuid.UUID  `gorm:"type:uuid;not null;index"`
  CommunityID  *uuid.UUID `gorm:"type:uuid;index;default:null"`
  Content      string     `gorm:"type:text;not null"`
  IsPinned bool `gorm:"not null;default:false"`

  ParentID     *uuid.UUID  `gorm:"type:uuid;index;default:null"`
  RepostOfID *uuid.UUID `gorm:"type:uuid;index;default:null"`
  LikeCount    int64      `gorm:"not null;default:0"`
  ShareCount   int64      `gorm:"not null;default:0"`
  CommentCount int64      `gorm:"not null;default:0"`
  ViewCount    int64      `gorm:"not null;default:0"`

  CreatedAt    time.Time  `gorm:"autoCreateTime"`
  UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}