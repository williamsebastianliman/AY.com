package grpc

import (
	"context"
	"log"
	"strconv"

	"github.com/google/uuid"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) ActivateAccount(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.RegisterUserResponse, error){
	user, err := h.svc.ActivateAccount(ctx, req.Id)
	if err != nil{
		return nil, err
	}
	return &userpb.RegisterUserResponse{Id: user.ID.String()}, nil
}

func (h *UserHandler) RegisterUser(
	ctx context.Context,
	req *userpb.RegisterUserRequest,
) (*userpb.RegisterUserResponse, error) {
	ppUID, err := uuid.Parse(req.ProfilePictureId)
	if err != nil {
		return nil, err
	}
	bmUID, err := uuid.Parse(req.BannerMediaId)
	if err != nil {
		return nil, err
	}

	id, err := h.svc.Register(
		ctx,
		req.Name,
		req.Username,
		req.Email,
		req.Password,
		req.Gender,
		req.BirthYear,
		req.BirthMonth,
		req.BirthDay,
		req.SubscribedNews,
		&ppUID,
		&bmUID,
	)
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterUserResponse{Id: id.String()}, nil
}

func (h *UserHandler) IsUsernameUnique(
	ctx context.Context,
	req *userpb.UsernameValidationRequest,
) (*userpb.UsernameValidationResponse, error) {
	isUnique, err := h.svc.IsUsernameUnique(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	return &userpb.UsernameValidationResponse{IsUnique: isUnique}, nil
}

func (h *UserHandler) GetUserByUsername(
	ctx context.Context,
	req *userpb.GetUserByUsernameRequest,
) (*userpb.GetUserByUsernameResponse, error) {
	user, err := h.svc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	year := strconv.Itoa(user.DateOfBirth.Year())
	month := strconv.Itoa(int(user.DateOfBirth.Month()))
	day := strconv.Itoa(user.DateOfBirth.Day())

	resp := &userpb.GetUserByUsernameResponse{
		Id:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Gender:         user.Gender,
		BirthYear:      year,
		BirthMonth:     month,
		BirthDay:       day,
		SubscribedNews: user.SubscribedNews,
		Password:       user.PasswordHash,
	}

	if user.ProfilePictureID != nil {
		resp.ProfilePictureId = user.ProfilePictureID.String()
	}
	if user.BannerMediaID != nil {
		resp.BannerMediaId = user.BannerMediaID.String()
	}

	return resp, nil
}

func (h *UserHandler) GetUserByEmail(
	ctx context.Context,
	req *userpb.GetUserByUsernameRequest,
) (*userpb.GetUserByUsernameResponse, error) {
	user, err := h.svc.GetUserByEmail(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	year := strconv.Itoa(user.DateOfBirth.Year())
	month := strconv.Itoa(int(user.DateOfBirth.Month()))
	day := strconv.Itoa(user.DateOfBirth.Day())

	resp := &userpb.GetUserByUsernameResponse{
		Id:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Gender:         user.Gender,
		BirthYear:      year,
		BirthMonth:     month,
		BirthDay:       day,
		SubscribedNews: user.SubscribedNews,
		Password:       user.PasswordHash,
		IsBanned: user.IsBanned,
	}

	if user.ProfilePictureID != nil {
		resp.ProfilePictureId = user.ProfilePictureID.String()
	}
	if user.BannerMediaID != nil {
		resp.BannerMediaId = user.BannerMediaID.String()
	}

	return resp, nil
}

func (h *UserHandler) GetUserById(
	ctx context.Context,
	req *userpb.GetUserByIdRequest,
) (*userpb.GetUserByIdResponse, error) {
	user, err := h.svc.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	year := strconv.Itoa(user.DateOfBirth.Year())
	month := strconv.Itoa(int(user.DateOfBirth.Month()))
	day := strconv.Itoa(user.DateOfBirth.Day())

	resp := &userpb.GetUserByIdResponse{
		Id:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Gender:         user.Gender,
		BirthYear:      year,
		BirthMonth:     month,
		BirthDay:       day,
		SubscribedNews: user.SubscribedNews,
		IsPremium: user.IsPremium,
	}

	if user.ProfilePictureID != nil {
		resp.ProfilePictureId = user.ProfilePictureID.String()
	}
	if user.BannerMediaID != nil {
		resp.BannerMediaId = user.BannerMediaID.String()
	}
	log.Printf("logzz: %t", resp.IsPremium);
	return resp, nil
}

func (h *UserHandler) GetMe(
	ctx context.Context,
	_ *userpb.GetMeRequest,
) (*userpb.GetMeResponse, error) {
	user, err := h.svc.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	year := strconv.Itoa(user.DateOfBirth.Year())
	month := strconv.Itoa(int(user.DateOfBirth.Month()))
	day := strconv.Itoa(user.DateOfBirth.Day())

	resp := &userpb.GetMeResponse{
		Id:             user.ID.String(),
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Gender:         user.Gender,
		BirthYear:      year,
		BirthMonth:     month,
		BirthDay:       day,
		SubscribedNews: user.SubscribedNews,
	}

	if user.ProfilePictureID != nil {
		resp.ProfilePictureId = user.ProfilePictureID.String()
	}
	if user.BannerMediaID != nil {
		resp.BannerMediaId = user.BannerMediaID.String()
	}

	return resp, nil
}


func (h *UserHandler) FollowUser(
	ctx context.Context,
	req *userpb.FollowUserRequest,
) (*userpb.FollowUserResponse, error) {
	err := h.svc.Follow(ctx, req.UserId, req.FollowerId)
	if err != nil {
		return &userpb.FollowUserResponse{Success: false}, err
	}
	return &userpb.FollowUserResponse{Success: true}, nil
}

func (h *UserHandler) UnfollowUser(
	ctx context.Context,
	req *userpb.UnfollowUserRequest,
) (*userpb.UnfollowUserResponse, error) {
	err := h.svc.UnfollowUser(ctx, req.UserId, req.FollowerId)
	if err != nil {
		return &userpb.UnfollowUserResponse{Success: false}, err
	}
	return &userpb.UnfollowUserResponse{Success: true}, nil
}

func (h *UserHandler) GetFollowers(
	ctx context.Context,
	req *userpb.GetFollowersRequest,
) (*userpb.GetFollowersResponse, error) {
	users, err := h.svc.GetFollowers(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var followers []*userpb.UserInfo
	for _, u := range users {
		followers = append(followers, &userpb.UserInfo{
			Id:               u.ID.String(),
			Username:         u.Username,
			Name:             u.Name,
			ProfilePictureId: u.ProfilePictureID.String(),
		})
	}
	return &userpb.GetFollowersResponse{Followers: followers}, nil
}

func (h *UserHandler) GetFollowing(
	ctx context.Context,
	req *userpb.GetFollowingRequest,
) (*userpb.GetFollowingResponse, error) {
	users, err := h.svc.GetFollowing(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var following []*userpb.UserInfo
	for _, u := range users {
		following = append(following, &userpb.UserInfo{
			Id:               u.ID.String(),
			Username:         u.Username,
			Name:             u.Name,
			ProfilePictureId: u.ProfilePictureID.String(),
		})
	}
	return &userpb.GetFollowingResponse{Following: following}, nil
}

func (h *UserHandler) FollowerCount(
    ctx context.Context,
    req *userpb.GetUserByIdRequest,
) (*userpb.CountResponse, error) {
    count, err := h.svc.CountFollowers(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not count followers: %v", err)
    }
    return &userpb.CountResponse{Count: count}, nil
}

func (h *UserHandler) FollowingCount(
    ctx context.Context,
    req *userpb.GetUserByIdRequest,
) (*userpb.CountResponse, error) {
    count, err := h.svc.CountFollowing(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not count following: %v", err)
    }
    return &userpb.CountResponse{Count: count}, nil
}

func (h *UserHandler) IsFollowed(
    ctx context.Context,
    req *userpb.GetFollowStatusRequest,
) (*userpb.IsFollowedResponse, error) {
    isFollowed, err := h.svc.IsFollowed(ctx, req.UserId, req.FollowerId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not check follow status: %v", err)
    }
    return &userpb.IsFollowedResponse{IsFollowed: isFollowed}, nil
}

func (h *UserHandler) RequestPremium(
    ctx context.Context,
    req *userpb.RequestPremiumRequest,
) (*userpb.RequestPremiumResponse, error) {
    if req.GetUserId() == "" || req.GetCardNumber() == "" || req.GetReason() == "" || req.GetFaceImage() == "" {
        return &userpb.RequestPremiumResponse{
            Success: false,
        }, status.Error(codes.InvalidArgument, "All fields are required")
    }

    err := h.svc.RequestPremium(
        ctx,
        req.GetUserId(),
        req.GetCardNumber(),
        req.GetReason(),
        req.GetFaceImage(),
    )
    if err != nil {
        return &userpb.RequestPremiumResponse{Success: false}, status.Error(codes.Internal, err.Error())
    }

    return &userpb.RequestPremiumResponse{
        Success: true,
    }, nil
}

func (h *UserHandler) IsPremium(
	ctx context.Context,
	req *userpb.IsPremiumRequest,
) (*userpb.RequestPremiumResponse, error) {
	isPremium, err := h.svc.IsPremium(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &userpb.RequestPremiumResponse{
		Success: isPremium,
	}, nil
}

func (h *UserHandler) GetAllPremiumRequests(
	ctx context.Context,
	req *userpb.GetAllPremiumRequestsRequest,
) (*userpb.GetAllPremiumRequestsResponse, error) {
	requests, err := h.svc.GetAllPremiumRequests(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get premium requests: %v", err)
	}

	var pbRequests []*userpb.PremiumRequest
	for _, r := range requests {
		pbReq := &userpb.PremiumRequest{
			UserId:     r.UserID.String(),
			CardNumber: r.CardNumber,
			Reason:     r.Reason,
		}
		if r.FaceImage != nil {
			pbReq.FaceImage = r.FaceImage.String()
		}
		pbRequests = append(pbRequests, pbReq)
	}

	return &userpb.GetAllPremiumRequestsResponse{
		Requests: pbRequests,
	}, nil
}

func (h *UserHandler) AcceptPremium(
	ctx context.Context,
	req *userpb.AcceptPremiumRequest,
) (*userpb.AcceptPremiumResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	err := h.svc.AcceptPremium(ctx, req.UserId)
	if err != nil {
		return &userpb.AcceptPremiumResponse{Success: false}, status.Errorf(codes.Internal, "failed to accept premium: %v", err)
	}
	return &userpb.AcceptPremiumResponse{
		Success: true,
	}, nil
}
func (h *UserHandler) RejectPremium(
	ctx context.Context,
	req *userpb.RejectPremiumRequest,
) (*userpb.RejectPremiumResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	err := h.svc.RejectPremium(ctx, req.UserId)
	if err != nil {
		return &userpb.RejectPremiumResponse{Success: false}, status.Errorf(codes.Internal, "failed to reject premium: %v", err)
	}
	return &userpb.RejectPremiumResponse{
		Success: true,
	}, nil
}

func (h *UserHandler) GetUserByName(
	ctx context.Context,
	req *userpb.GetUserByUsernameRequest,
) (*userpb.GetUserByNameResponse, error) {
	users, err := h.svc.GetUserByName(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users by name: %v", err)
	}

	var infos []*userpb.UserInfo
	for _, user := range users {
		info := &userpb.UserInfo{
			Id:               user.ID.String(),
			Username:         user.Username,
			Name:             user.Name,
		}
		if user.ProfilePictureID != nil {
			info.ProfilePictureId = user.ProfilePictureID.String()
		}
		infos = append(infos, info)
	}

	return &userpb.GetUserByNameResponse{
		Users: infos,
	}, nil
}

func (h *UserHandler) UpdateProfile(
    ctx context.Context,
    req *userpb.UpdateProfileRequest,
) (*userpb.UpdateProfileResponse, error) {
    var profilePictureID, bannerMediaID *uuid.UUID
    if req.ProfilePictureId != "" {
        id, err := uuid.Parse(req.ProfilePictureId)
        if err != nil {
            return nil, status.Errorf(codes.InvalidArgument, "invalid profile_picture_id: %v", err)
        }
        profilePictureID = &id
    }
    if req.BannerMediaId != "" {
        id, err := uuid.Parse(req.BannerMediaId)
        if err != nil {
            return nil, status.Errorf(codes.InvalidArgument, "invalid banner_media_id: %v", err)
        }
        bannerMediaID = &id
    }

    err := h.svc.UpdateProfile(
        ctx,
        req.UserId,
        req.Name,
        req.Username,
        req.Email,
        req.Password,
        req.Gender,
        req.BirthYear,
        req.BirthMonth,
        req.BirthDay,
        profilePictureID,
        bannerMediaID,
    )
    if err != nil {
        return &userpb.UpdateProfileResponse{Success: false}, status.Errorf(codes.Internal, "failed to update profile: %v", err)
    }

    return &userpb.UpdateProfileResponse{Success: true}, nil
}

func (h *UserHandler) BlockUser(
	ctx context.Context,
	req *userpb.BlockUserRequest,
) (*userpb.BlockUserResponse, error) {
	if req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and target_id are required")
	}

	err := h.svc.BlockUser(ctx, req.UserId, req.TargetId)
	if err != nil {
		return &userpb.BlockUserResponse{Success: false}, status.Errorf(codes.Internal, "failed to block user: %v", err)
	}

	return &userpb.BlockUserResponse{Success: true}, nil
}

func (h *UserHandler) UnblockUser(
	ctx context.Context,
	req *userpb.UnblockUserRequest,
) (*userpb.UnblockUserResponse, error) {
	if req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and target_id are required")
	}

	err := h.svc.UnblockUser(ctx, req.UserId, req.TargetId)
	if err != nil {
		return &userpb.UnblockUserResponse{Success: false}, status.Errorf(codes.Internal, "failed to unblock user: %v", err)
	}

	return &userpb.UnblockUserResponse{Success: true}, nil
}

func (h *UserHandler) IsUserBlocked(
	ctx context.Context,
	req *userpb.IsUserBlockedRequest,
) (*userpb.IsUserBlockedResponse, error) {
	if req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and target_id are required")
	}

	isBlocked, err := h.svc.IsUserBlocked(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check block status: %v", err)
	}

	return &userpb.IsUserBlockedResponse{IsBlocked: isBlocked}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
    users, err := h.svc.GetAllUsers(ctx)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "could not fetch users: %v", err)
    }

    var resp userpb.GetAllUsersResponse
    for _, user := range users {
        resp.Users = append(resp.Users, &userpb.UserInfo{
            Id:               user.ID.String(),
            Username:         user.Username,
            Name:             user.Name,
            ProfilePictureId: user.ProfilePictureID.String(),
        })
    }
    return &resp, nil
}

func (h *UserHandler) CreateUserReport(ctx context.Context, req *userpb.CreateUserReportRequest) (*userpb.CreateUserReportResponse, error) {
	reportID, err := h.svc.CreateUserReport(req.ReportedUser, req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user report: %v", err)
	}
	return &userpb.CreateUserReportResponse{ReportId: reportID}, nil
}

func (h *UserHandler) GetAllUserReports(ctx context.Context, req *userpb.GetAllUserReportsRequest) (*userpb.GetAllUserReportsResponse, error) {
    reports, err := h.svc.GetAllUserReports()
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch user reports: %v", err)
    }
    var pbReports []*userpb.UserReport
    for _, r := range reports {
        pbReports = append(pbReports, &userpb.UserReport{
            ReportId:     r.ReportID.String(),
            ReportedUser: r.ReportedUser.String(),
            Reason:       r.Reason,
        })
    }
    return &userpb.GetAllUserReportsResponse{
        Reports: pbReports,
    }, nil
}

func (h *UserHandler) ChangePassword(
    ctx context.Context,
    req *userpb.ChangePasswordRequest,
) (*userpb.ChangePasswordResponse, error) {
    if req.UserId == "" || req.NewPassword == "" {
        return &userpb.ChangePasswordResponse{
            Success:      false,
            ErrorMessage: "user_id and new_password are required",
        }, status.Error(codes.InvalidArgument, "user_id and new_password are required")
    }

    err := h.svc.ChangePassword(ctx, req.UserId, req.NewPassword)
    if err != nil {
        return &userpb.ChangePasswordResponse{
            Success:      false,
            ErrorMessage: err.Error(),
        }, status.Errorf(codes.Internal, "failed to change password: %v", err)
    }

    return &userpb.ChangePasswordResponse{Success: true}, nil
}

func (s *UserHandler) ExploreUserByName(ctx context.Context, req *userpb.ExploreUserByNameRequest) (*userpb.ExploreUserByNameResponse, error) {
    users, total, err := s.svc.ExploreUserByName(req.Query, int(req.Threshold), int(req.Page), int(req.Size))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "explore failed: %v", err)
    }

    var pbUsers []*userpb.UserInfo
    for _, u := range users {
        pbUsers = append(pbUsers, &userpb.UserInfo{
            Id: u.ID.String(),
            Username: u.Username,
            Name: u.Name,
            ProfilePictureId: u.ProfilePictureID.String(),
        })
    }
    return &userpb.ExploreUserByNameResponse{
        Users: pbUsers,
        Total: int32(total),
    }, nil
}
