package grpc

import (
	"context"

	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/service"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) SendOTP(
	ctx context.Context,
	req *authpb.SendOTPRequest,
) (*authpb.SendOTPResponse, error) {

	requestID, err := h.svc.SendOTP(ctx, req.Email, req.UserId)
	if err != nil {
		return nil, err
	}
	return &authpb.SendOTPResponse{RequestId: requestID}, nil
}

func (h *AuthHandler) VerifyOTP(
	ctx context.Context,
	req *authpb.VerifyOTPRequest,
) (*authpb.VerifyOTPResponse, error) {
	success, err := h.svc.VerifyOTP(ctx, req.RequestId, req.Code)
	if err != nil {
		return nil, err
	}
	return &authpb.VerifyOTPResponse{Success: success}, nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	accessToken, refreshToken, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AuthHandler) RefreshToken(
	ctx context.Context,
	req *authpb.RefreshTokenRequest,
) (*authpb.RefreshTokenResponse, error) {
	newAccess, err := h.svc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &authpb.RefreshTokenResponse{
		AccessToken: newAccess,
	}, nil
}