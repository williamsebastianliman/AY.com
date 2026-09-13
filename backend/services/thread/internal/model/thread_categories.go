package model

import (
	"time"

	"github.com/google/uuid"
)

type ThreadCategories struct {
    CategoryID uuid.UUID `gorm:"type:uuid;not null;index;primaryKey"`
 	CategoryName string `gorm:"type:text;not null;"`

    CreatedAt time.Time `gorm:"autoCreateTime"`
}