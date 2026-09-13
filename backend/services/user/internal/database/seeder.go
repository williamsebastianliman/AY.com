package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/model"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
    dob1 := time.Date(2004, time.November, 21, 0, 0, 0, 0, time.UTC)
    activatedAt1 := time.Now()

    user1 := &model.User{
        ID:               uuid.MustParse("ec724072-ddf8-4ebe-a6ec-ce041375a1e1"),
        Name:             "William Sebastian Liman",
        Username:         "DocCom",
        Email:            "william.sebastian.liman@gmail.com",
        PasswordHash:     "$2a$10$8jP4TTqNVHllbe9G372GGeg7VVIqvJzF63IWJlDdvtE3F90Bcb88q",
        Gender:           "male",
        DateOfBirth:      dob1,
        IsActivated:      true,
        ActivatedAt:      &activatedAt1,
        IsBanned:         false,
        SubscribedNews:   false,
        ProfilePictureID: uuidPtr("9073750f-e4ab-48d3-8b5f-af6b03fd2fd6"),
        BannerMediaID:    uuidPtr("1c2d3bd4-89f1-417c-9647-af99b4ac831f"),
    }

    dob2 := time.Date(2003, time.March, 12, 0, 0, 0, 0, time.UTC)
    activatedAt2 := time.Now()

    user2 := &model.User{
        ID:               uuid.MustParse("8af807c2-0b1c-4d83-bf18-fc7e6315c5f2"),
        Name:             "Gloria",
        Username:         "Gloria",
        Email:            "gloria@gmail.com",
        PasswordHash:     "$2a$10$8jP4TTqNVHllbe9G372GGeg7VVIqvJzF63IWJlDdvtE3F90Bcb88q",
        Gender:           "female",
        DateOfBirth:      dob2,
        IsActivated:      true,
        ActivatedAt:      &activatedAt2,
        IsBanned:         false,
        SubscribedNews:   true,
        ProfilePictureID: uuidPtr("1c2d3bd4-89f1-417c-9647-af99b4ac831f"),
        BannerMediaID:    uuidPtr("9073750f-e4ab-48d3-8b5f-af6b03fd2fd6"),
    }
    admin := &model.User{
        ID:               uuid.MustParse("5f08da0e-01a6-4a0b-8f3d-5a8fbb0f8764"),
        Name:             "Admin",
        Username:         "Admin",
        Email:            "admin@gmail.com",
        PasswordHash:     "$2a$10$8jP4TTqNVHllbe9G372GGeg7VVIqvJzF63IWJlDdvtE3F90Bcb88q",
        Gender:           "male",
        DateOfBirth:      dob2,
        IsActivated:      true,
        ActivatedAt:      &activatedAt2,
        IsBanned:         false,
        SubscribedNews:   true,
        ProfilePictureID: uuidPtr("1c2d3bd4-89f1-417c-9647-af99b4ac831f"),
        BannerMediaID:    uuidPtr("9073750f-e4ab-48d3-8b5f-af6b03fd2fd6"),
    }
    return db.Create([]*model.User{user1, user2, admin}).Error
}

func uuidPtr(s string) *uuid.UUID {
	u := uuid.MustParse(s)
	return &u
}