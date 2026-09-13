package database

import (
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}

	return db.AutoMigrate(
        &model.UserSecurityAnswer{},
        &model.RefreshToken{},
    )
}
