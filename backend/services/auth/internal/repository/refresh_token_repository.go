package repository

import (
	"errors"
	"time"

	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(rt *model.RefreshToken) error
	FindByToken(token string) (*model.RefreshToken, error)
	UpdateByToken(token string, newExpiresAt time.Time) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(rt *model.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.
		Where("token = ?", token).
		First(&rt).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}
func (r *refreshTokenRepository) UpdateByToken(token string, newExpiresAt time.Time) error {
	result := r.db.
		Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("expires_at", newExpiresAt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}