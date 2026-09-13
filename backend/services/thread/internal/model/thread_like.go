// thread_like.go
package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadLike struct {
    ThreadID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
    UserID   uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`

    CreatedAt time.Time `gorm:"autoCreateTime"`
}