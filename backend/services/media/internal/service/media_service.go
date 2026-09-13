package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	storage "github.com/supabase-community/storage-go"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/repository"
	"gorm.io/gorm"
)

type MediaService interface {
	Upload(ctx context.Context, fileBytes []byte, filename string) (*model.Media, error)
	GetMediaByID(ctx context.Context, id uuid.UUID) (*model.Media, error)
}

type mediaService struct {
	repo          repository.MediaRepository
	storageClient *storage.Client
	supabaseURL   string
	bucket        string
}

func NewMediaService(
	repo repository.MediaRepository,
	storageURL, serviceKey, bucket string,
) MediaService {
	client := storage.NewClient(storageURL, serviceKey, nil)
	return &mediaService{
		repo:          repo,
		storageClient: client,
		supabaseURL:   storageURL,
		bucket:        bucket,
	}
}

func (s *mediaService) Upload(ctx context.Context, fileBytes []byte, filename string) (*model.Media, error) {
	ext := filepath.Ext(filename)

	objectKey := fmt.Sprintf("public/%s%s", uuid.New().String(), ext)

	reader := bytes.NewReader(fileBytes)
	if _, err := s.storageClient.UploadFile(s.bucket, objectKey, reader); err != nil {
		return nil, fmt.Errorf("storage upload failed: %w", err)
	}

	publicURL := fmt.Sprintf("%s/object/%s/%s", s.supabaseURL, s.bucket, objectKey)

	media, err := s.repo.Create(publicURL, ext)
	if err != nil {
		return nil, fmt.Errorf("db insert failed: %w", err)
	}

	return media, nil
}

func (s *mediaService) GetMediaByID(ctx context.Context, id uuid.UUID) (*model.Media, error) {
	media, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("media %s not found", id)
		}
		return nil, err
	}
	return media, nil
}