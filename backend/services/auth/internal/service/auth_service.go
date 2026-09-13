package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"golang.org/x/crypto/bcrypt"

	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/repository"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type AuthService interface {
	SendOTP(ctx context.Context, email string, id string) (requestID string, err error)
	VerifyOTP(ctx context.Context, requestID, code string) (success bool, err error)
	Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, err error)
}

type authService struct {
	cache      CacheService
	mailer     MailService
	eventBus   EventBus
	userClient userpb.UserServiceClient
	rtRepo     repository.RefreshTokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
}

func NewAuthService(
	cache CacheService,
	mailer MailService,
	bus EventBus,
	userClient userpb.UserServiceClient,
	rtRepo repository.RefreshTokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
) AuthService {
	return &authService{
		cache:      cache,
		mailer:     mailer,
		eventBus:   bus,
		userClient: userClient,
		rtRepo:     rtRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
	}
}

func (s *authService) SendOTP(ctx context.Context, email string, id string) (string, error) {
	rid := uuid.New().String()
	code := fmt.Sprintf("%06d", uuid.New().ID()%1_000_000)

	entry := struct {
		Code string
		Email string
		Id string
	}{Code: code, Email: email, Id:id}

	buf, _ := json.Marshal(entry)
	if err := s.cache.Set(ctx, "otp:"+rid, string(buf), 5*time.Minute); err != nil {
		return "", err
	}
	go s.mailer.SendCode(email, code)
	return rid, nil
}

func (s *authService) VerifyOTP(ctx context.Context, requestID, code string) (bool, error) {
	key := "otp:" + requestID
	raw, err := s.cache.Get(ctx, key)
	if err != nil {
		return false, nil
	}
	var entry struct {
		Code string
		Email string
		Id string
	}
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return false, nil
	}
	if entry.Code != code {
		return false, nil
	}
	payload := map[string]interface{}{
		"id" : entry.Id,
		"email":         entry.Email,
	}
	_ = s.eventBus.Publish(ctx, "user.activate", "user_events", payload)
	_ = s.cache.Delete(ctx, key)
	return true, nil
}

func (s *authService) Login(
	ctx context.Context,
	username, password string,
) (string, string, error) {
	resp, err := s.userClient.GetUserByUsername(ctx, &userpb.GetUserByUsernameRequest{Username: username})
	if err != nil {
		return "", "", ErrUserNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(resp.Password), []byte(password)) != nil {
		return "", "", ErrInvalidCredential
	}

	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id": resp.Id,
		"iat":     now.Unix(),
		"exp":     now.Add(s.accessTTL).Unix(),
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessJWT.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()
	expiresAt := now.Add(72 * time.Hour)
	rtRecord := &model.RefreshToken{
		UserID:    uuid.MustParse(resp.Id),
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.rtRepo.Create(rtRecord); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

func (s *authService) Refresh(
    ctx context.Context,
    refreshToken string,
) (newAccessToken string, err error) {
    
    rec, err := s.rtRepo.FindByToken(refreshToken)
    if err != nil {
        return "", err
    }
    if rec == nil || rec.Revoked || rec.ExpiresAt.Before(time.Now()) {
        return "", ErrInvalidRefreshToken
    }

    newExpiry := time.Now().Add(72 * time.Hour)
    if err := s.rtRepo.UpdateByToken(refreshToken, newExpiry); err != nil {
        return "", err
    }

    now := time.Now()
    claims := jwt.MapClaims{
        "user_id": rec.UserID.String(),
        "iat":     now.Unix(),
        "exp":     now.Add(s.accessTTL).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(s.jwtSecret)
    if err != nil {
        return "", err
    }
    return signed, nil
}