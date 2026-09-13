package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type SecurityAnswer interface {
	Create(ctx context.Context, userID uuid.UUID, question string, answer string) (requestID string, err error)
	GetByUser(ctx context.Context,userID uuid.UUID) ([]model.UserSecurityAnswer, error)
	CheckAnswer(ctx context.Context, userID uuid.UUID, question, providedAnswer string) (bool, error)
}

type securityAnswerService struct {
	repo repository.SecurityAnswerRepository
}

func NewSecurityAnswerService(repo repository.SecurityAnswerRepository) SecurityAnswer {
	return &securityAnswerService{repo: repo}
}

func (s *securityAnswerService) Create(ctx context.Context, userID uuid.UUID, question string, answer string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	sa := &model.UserSecurityAnswer{
		ID:         uuid.New(),
		UserID:     userID,
		Question: question,
		AnswerHash: string(hash),
	}
	if err := s.repo.Create(sa); err != nil {
		return "", err
	}
	return sa.ID.String(), nil
}

func (s *securityAnswerService) GetByUser(
    ctx context.Context,
    userID uuid.UUID,
) ([]model.UserSecurityAnswer, error) {
    return s.repo.GetByUser(userID)
}

func (s *securityAnswerService) CheckAnswer(ctx context.Context, userID uuid.UUID, question, providedAnswer string) (bool, error){
	flag, err := s.repo.CheckAnswer(userID, question, providedAnswer)
	if err != nil{
		return false, err
	}
	return flag, nil
}

