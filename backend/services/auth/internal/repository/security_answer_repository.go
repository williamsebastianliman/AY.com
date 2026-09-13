package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SecurityAnswerRepository interface {
	Create(answer *model.UserSecurityAnswer) error
	GetByUser(userID uuid.UUID) ([]model.UserSecurityAnswer, error)
	CheckAnswer(userID uuid.UUID, question, providedAnswer string) (bool, error)
}

type securityAnswerRepository struct {
	db *gorm.DB
}

func NewSecurityAnswerRepository(db *gorm.DB) SecurityAnswerRepository {
	return &securityAnswerRepository{db: db}
}

func (r *securityAnswerRepository) Create(answer *model.UserSecurityAnswer) error {
	return r.db.Create(answer).Error
}

func (r *securityAnswerRepository) GetByUser(userID uuid.UUID) ([]model.UserSecurityAnswer, error) {
	var answers []model.UserSecurityAnswer
	err := r.db.
		Where(&model.UserSecurityAnswer{UserID: userID}).
		Find(&answers).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return answers, err
}

func (r *securityAnswerRepository) CheckAnswer(userID uuid.UUID, question, providedAnswer string) (bool, error) {
	var answer model.UserSecurityAnswer
	err := r.db.
		Where("user_id = ? AND question = ?", userID, question).
		First(&answer).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(answer.AnswerHash), []byte(providedAnswer)); err != nil {
		return false, nil
	}
	return true, nil
}