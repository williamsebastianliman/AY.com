package grpc

import (
	"context"

	"github.com/google/uuid"
	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/service"
)

type SecurityAnswerHandler struct {
	authpb.UnimplementedSecurityAnswerServiceServer
	svc service.SecurityAnswer
	userClient userpb.UserServiceClient
}

func NewSecurityAnswerHandler(svc service.SecurityAnswer, userClient userpb.UserServiceClient) *SecurityAnswerHandler {
	return &SecurityAnswerHandler{svc: svc, userClient: userClient}
}

func (h *SecurityAnswerHandler) CreateSecurityAnswer(
	ctx context.Context,
	req *authpb.CreateSecurityAnswerRequest,
) (*authpb.CreateSecurityAnswerResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	requestID, err := h.svc.Create(ctx, userID, req.Question, req.Answer)
	if err != nil {
		return nil, err
	}

	return &authpb.CreateSecurityAnswerResponse{
		RequestId: requestID,
	}, nil
}

func (h *SecurityAnswerHandler) GetSecurityAnswers(
	req *authpb.GetSecurityAnswerRequest,
	stream authpb.SecurityAnswerService_GetSecurityAnswersServer,
) error {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return err
	}

	answers, err := h.svc.GetByUser(stream.Context(), userID)
	if err != nil {
		return err
	}

	for _, a := range answers {
		if err := stream.Send(&authpb.GetSecurityAnswerResponse{
			Id:         a.ID.String(),
			UserId:     a.UserID.String(),
			Question:   a.Question,
			AnswerHash: a.AnswerHash,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (h *SecurityAnswerHandler) CheckSecurityAnswer(
    ctx context.Context,
    req *authpb.CheckSecurityAnswerRequest,
) (*authpb.CheckSecurityAnswerResponse, error) {
    userID, err := uuid.Parse(req.UserId)
    if err != nil {
        return &authpb.CheckSecurityAnswerResponse{
            Success:      false,
            ErrorMessage: "invalid user_id",
        }, nil
    }

    ok, err := h.svc.CheckAnswer(ctx, userID, req.Question, req.Answer)
    if err != nil {
        return &authpb.CheckSecurityAnswerResponse{
            Success:      false,
            ErrorMessage: err.Error(),
        }, nil
    }
    if !ok {
        return &authpb.CheckSecurityAnswerResponse{
            Success:      false,
            ErrorMessage: "security answer is incorrect",
        }, nil
    }

    _, err = h.userClient.ChangePassword(ctx, &userpb.ChangePasswordRequest{
        UserId:     req.UserId,
        NewPassword: req.NewPassword,
    })
    if err != nil {
        return &authpb.CheckSecurityAnswerResponse{
            Success:      false,
            ErrorMessage: "failed to change password: " + err.Error(),
        }, nil
    }

    return &authpb.CheckSecurityAnswerResponse{
        Success: true,
    }, nil
}


