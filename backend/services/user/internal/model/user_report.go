package model

import (
	"github.com/google/uuid"
)

type UserReport struct{

	ReportID uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReportedUser uuid.UUID `gorm:"type:uuid"`
	Reason string `gorm:"type:text"`
}