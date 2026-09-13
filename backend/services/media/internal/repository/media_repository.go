package repository

import (
	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/model"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(publicURL string, extension string) (*model.Media, error)

	GetByID(id uuid.UUID) (*model.Media, error)
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(publicURL string, extension string) (*model.Media, error) {
	m := &model.Media{
		PublicURL: publicURL,
		Extension: extension,
	}

	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *mediaRepository) GetByID(id uuid.UUID) (*model.Media, error) {
	var m model.Media
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}