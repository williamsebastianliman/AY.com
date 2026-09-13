package database

import (
	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/model"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	mediaEntries := []model.Media{
		{
			ID:        uuid.MustParse("9073750f-e4ab-48d3-8b5f-af6b03fd2fd6"),
			PublicURL: "https://itmuqguefzemfzsbryad.supabase.co/storage/v1/object/avatars/public/553d0605-a8b6-4456-8535-0b1c2e24c7af.png",
			Extension: "png",
		},
		{
			ID:        uuid.MustParse("1c2d3bd4-89f1-417c-9647-af99b4ac831f"),
			PublicURL: "https://itmuqguefzemfzsbryad.supabase.co/storage/v1/object/avatars/public/ecbb16de-8532-420d-98e1-88656539979b.png",
			Extension: "png",
		},
	}

	if err := db.Create(&mediaEntries).Error; err != nil {
		return err
	}

	return nil
}