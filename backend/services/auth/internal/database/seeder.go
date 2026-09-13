package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("Pussy"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	answer := &model.UserSecurityAnswer{
		ID:         uuid.New(),
		UserID:     uuid.MustParse("ec724072-ddf8-4ebe-a6ec-ce041375a1e1"),
		Question:   "What was your first pet's name?",
		AnswerHash: string(hash),
		CreatedAt:  time.Now(),
	}

	return db.Create(answer).Error
}