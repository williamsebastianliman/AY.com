package grpc

import (
	"context"
	"log"
	"sort"
	"strings"
	"time"

	communitypb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/communitypb/proto/community"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommunityHandler struct {
	communitypb.UnimplementedCommunityServiceServer
	svc service.CommunityService
	userClient userpb.UserServiceClient
}

func NewCommunityHandler(svc service.CommunityService, userClient userpb.UserServiceClient) *CommunityHandler {
	return &CommunityHandler{svc: svc, userClient: userClient}
}

func toTimestamp(t time.Time) *communitypb.Timestamp {
	return &communitypb.Timestamp{
		Year:      int32(t.Year()),
		Month:     int32(int(t.Month())),
		Day:       int32(t.Day()),
		Hour:      int32(t.Hour()),
		Minute:    int32(t.Minute()),
		Second:    int32(t.Second()),
		Timezone:  t.Format("MST"),
	}
}

func toPbCommunity(c *model.Community) *communitypb.Community {
	return &communitypb.Community{
		CommunityId:          c.CommunityID.String(),
		CommunityName:        c.CommunityName,
		CommunityDescription: c.CommunityDescription,
		CommunityRules:       c.CommunityRules,
		Status:               c.Status,
		CreatedAt:            toTimestamp(c.CreatedAt),
		CommunityLogo:        c.CommunityLogo,
		CommunityBanner:      c.CommunityBanner,
	}
}
func (h *CommunityHandler) ListCommunities(ctx context.Context, req *communitypb.ListCommunitiesRequest) (*communitypb.ListCommunitiesResponse, error) {
	communities, err := h.svc.ListCommunities(ctx, req.Filter, &req.CategoryId, int(req.Page), int(req.Size), req.UserId)
	if err != nil {
		return nil, err
	}
	var result []*communitypb.Community
	for _, c := range communities {
		result = append(result, toPbCommunity(&c))
	}
	return &communitypb.ListCommunitiesResponse{Communities: result}, nil
}

func (h *CommunityHandler) SearchCommunities(ctx context.Context, req *communitypb.SearchCommunitiesRequest) (*communitypb.SearchCommunitiesResponse, error) {
	communities, err := h.svc.SearchCommunities(ctx, req.Keyword, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var result []*communitypb.Community
	for _, c := range communities {
		result = append(result, toPbCommunity(&c))
	}
	return &communitypb.SearchCommunitiesResponse{Communities: result}, nil
}

func (h *CommunityHandler) GetCommunityByID(ctx context.Context, req *communitypb.GetCommunityByIDRequest) (*communitypb.GetCommunityByIDResponse, error) {
	community, err := h.svc.GetCommunityByID(ctx, req.CommunityId)
	if err != nil {
		return nil, err
	}
	if community == nil {
		return nil, status.Error(codes.NotFound, "community not found")
	}
	return &communitypb.GetCommunityByIDResponse{
		Community: toPbCommunity(community),
	}, nil
}

func (h *CommunityHandler) CreateCommunity(ctx context.Context, req *communitypb.CreateCommunityRequest) (*communitypb.CreateCommunityResponse, error) {
	communityID, err := h.svc.CreateCommunity(
		ctx,
		req.CommunityName,
		req.CommunityDescription,
		req.CommunityRules,
		req.CommunityLogo,
		req.CommunityBanner,
		req.Status,
		req.CategoryIds,
		req.OwnerId,
	)
	if err != nil {
		return nil, err
	}
	return &communitypb.CreateCommunityResponse{CommunityId: communityID}, nil
}

func (h *CommunityHandler) RequestJoinCommunity(ctx context.Context, req *communitypb.RequestJoinCommunityRequest) (*communitypb.Empty, error) {
	err := h.svc.RequestJoinCommunity(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) ListJoinedCommunities(ctx context.Context, req *communitypb.ListJoinedCommunitiesRequest) (*communitypb.ListJoinedCommunitiesResponse, error) {
	log.Printf("uidz: %s ",req.UserId)
	communities, err := h.svc.ListJoinedCommunities(ctx, req.UserId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var result []*communitypb.Community
	for _, c := range communities {
		result = append(result, toPbCommunity(&c))
	}
	return &communitypb.ListJoinedCommunitiesResponse{Communities: result}, nil
}

func (h *CommunityHandler) ListRequestedCommunities(ctx context.Context, req *communitypb.ListRequestedCommunitiesRequest) (*communitypb.ListRequestedCommunitiesResponse, error) {
	communities, err := h.svc.ListRequestedCommunities(ctx, req.UserId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	var result []*communitypb.Community
	for _, c := range communities {
		result = append(result, toPbCommunity(&c))
	}
	return &communitypb.ListRequestedCommunitiesResponse{Communities: result}, nil
}

func (h *CommunityHandler) ListModeratedCommunities(ctx context.Context, req *communitypb.ListModeratedCommunitiesRequest) (*communitypb.ListModeratedCommunitiesResponse, error) {
	communities, err := h.svc.ListModeratedCommunities(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var result []*communitypb.Community
	for _, c := range communities {
		result = append(result, toPbCommunity(&c))
	}
	return &communitypb.ListModeratedCommunitiesResponse{Communities: result}, nil
}

func (h *CommunityHandler) ListCommunityMembers(ctx context.Context, req *communitypb.ListCommunityMembersRequest) (*communitypb.ListCommunityMembersResponse, error) {
	var memberIds []string

	userResp, err := h.userClient.GetUserByName(ctx, &userpb.GetUserByUsernameRequest{Username: req.Name})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users by name: %v", err)
	}
	for _, u := range userResp.Users {
		memberIds = append(memberIds, u.Id)
	}
	
	log.Printf("member ids: %s", strings.Join(memberIds, ", "))

	members, err := h.svc.ListCommunityMembers(ctx, req.CommunityId, req.Role, int(req.Page), int(req.Size), memberIds)
	if err != nil {
		return nil, err
	}
	var result []*communitypb.CommunityMember
	for _, m := range members {
		result = append(result, &communitypb.CommunityMember{
			CommunityId: m.CommunityID.String(),
			MemberId:    m.MemberID.String(),
			Status:      m.Status,
			CreatedAt:   toTimestamp(m.CreatedAt),
		})
	}
	return &communitypb.ListCommunityMembersResponse{Members: result}, nil
}

func (h *CommunityHandler) PromoteMember(ctx context.Context, req *communitypb.PromoteMemberRequest) (*communitypb.Empty, error) {
	err := h.svc.PromoteMember(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) DemoteMember(ctx context.Context, req *communitypb.DemoteMemberRequest) (*communitypb.Empty, error) {
	err := h.svc.DemoteMember(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) AcceptJoinRequest(ctx context.Context, req *communitypb.AcceptJoinRequestRequest) (*communitypb.Empty, error) {
	err := h.svc.AcceptJoinRequest(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) RejectJoinRequest(ctx context.Context, req *communitypb.RejectJoinRequestRequest) (*communitypb.Empty, error) {
	err := h.svc.RejectJoinRequest(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) CommunityAbout(ctx context.Context, req *communitypb.CommunityAboutRequest) (*communitypb.CommunityAboutResponse, error) {
	community, err := h.svc.GetCommunityByID(ctx, req.CommunityId)
	if err != nil {
		return nil, err
	}
	if community == nil {
		return nil, status.Error(codes.NotFound, "community not found")
	}
	return &communitypb.CommunityAboutResponse{
		Community: toPbCommunity(community),
	}, nil
}

func (h *CommunityHandler) ListCommunityCategories(ctx context.Context, req *communitypb.ListCommunityCategoriesRequest) (*communitypb.ListCommunityCategoriesResponse, error) {
	cats, err := h.svc.ListCommunityCategories(ctx, req.CommunityId)
	if err != nil {
		return nil, err
	}
	var result []*communitypb.CommunityCategory
	for _, cat := range cats {
		result = append(result, &communitypb.CommunityCategory{
			CategoryId:   cat.CategoryID.String(),
			CategoryName: cat.CategoryName,
			CreatedAt:    toTimestamp(cat.CreatedAt),
		})
	}
	return &communitypb.ListCommunityCategoriesResponse{Categories: result}, nil
}

func (h *CommunityHandler) ListCategories(ctx context.Context, _ *communitypb.Empty) (*communitypb.ListCategoriesResponse, error) {
	cats, err := h.svc.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	var result []*communitypb.CommunityCategory
	for _, cat := range cats {
		result = append(result, &communitypb.CommunityCategory{
			CategoryId:   cat.CategoryID.String(),
			CategoryName: cat.CategoryName,
			CreatedAt:    toTimestamp(cat.CreatedAt),
		})
	}
	return &communitypb.ListCategoriesResponse{Categories: result}, nil
}

func (h *CommunityHandler) CreateCategory(ctx context.Context, req *communitypb.CreateCategoryRequest) (*communitypb.CreateCategoryResponse, error) {
	id, err := h.svc.CreateCategory(ctx, req.CategoryName)
	if err != nil {
		return nil, err
	}
	return &communitypb.CreateCategoryResponse{CategoryId: id}, nil
}

func (h *CommunityHandler) UpdateCategory(ctx context.Context, req *communitypb.UpdateCategoryRequest) (*communitypb.Empty, error) {
	err := h.svc.UpdateCategory(ctx, req.CategoryId, req.CategoryName)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) DeleteCategory(ctx context.Context, req *communitypb.DeleteCategoryRequest) (*communitypb.Empty, error) {
	err := h.svc.DeleteCategory(ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) AddCommunityCategory(ctx context.Context, req *communitypb.AddCommunityCategoryRequest) (*communitypb.Empty, error) {
	err := h.svc.AddCommunityCategory(ctx, req.CommunityId, req.CategoryId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) RemoveCommunityCategory(ctx context.Context, req *communitypb.RemoveCommunityCategoryRequest) (*communitypb.Empty, error) {
	err := h.svc.RemoveCommunityCategory(ctx, req.CommunityId, req.CategoryId)
	if err != nil {
		return nil, err
	}
	return &communitypb.Empty{}, nil
}

func (h *CommunityHandler) ListPendingCommunities(
    ctx context.Context,
    req *communitypb.Empty,
) (*communitypb.ListRequestedCommunitiesResponse, error) {
    comms, err := h.svc.ListPendingCommunities()
    if err != nil {
        return nil, err
    }

    var pbComms []*communitypb.Community
    for _, c := range comms {
        pbComms = append(pbComms, toPbCommunity(c))
    }

    return &communitypb.ListRequestedCommunitiesResponse{
        Communities: pbComms,
    }, nil
}

func (h *CommunityHandler) AcceptCommunityCreation(
    ctx context.Context,
    req *communitypb.AcceptRejectCommunityRequest,
) (*communitypb.CommunityOperationResponse, error) {
    err := h.svc.AcceptCommunityCreation(ctx, req.CommunityId)
    if err != nil {
        return &communitypb.CommunityOperationResponse{Message: err.Error()}, err
    }
    return &communitypb.CommunityOperationResponse{Message: "success"}, nil
}

func (h *CommunityHandler) RejectCommunityCreation(
    ctx context.Context,
    req *communitypb.AcceptRejectCommunityRequest,
) (*communitypb.CommunityOperationResponse, error) {
    err := h.svc.RejectCommunityCreation(ctx, req.CommunityId)
    if err != nil {
        return &communitypb.CommunityOperationResponse{Message: err.Error()}, err
    }
    return &communitypb.CommunityOperationResponse{Message: "success"}, nil
}

type MemberFollower struct {
    MemberID     string
    FollowerCount int
}

func (h *CommunityHandler) GetTopCommunityMembers(
    ctx context.Context,
    req *communitypb.GetTopCommunityMembersRequest,
) (*communitypb.GetTopCommunityMembersResponse, error) {
    memberIDs, err := h.svc.GetTopCommunityMembers(ctx, req.CommunityId)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to get community members: %v", err)
    }

    type mf struct {
        MemberID      string
        FollowerCount int
    }
    members := make([]mf, 0, len(memberIDs))

    for _, memberId := range memberIDs {
        followersResp, err := h.userClient.GetFollowers(ctx, &userpb.GetFollowersRequest{UserId: memberId})
        if err != nil {
            continue
        }
        members = append(members, mf{
            MemberID:      memberId,
            FollowerCount: len(followersResp.Followers),
        })
    }

    sort.Slice(members, func(i, j int) bool {
        return members[i].FollowerCount > members[j].FollowerCount
    })

    top := 3
    if len(members) < 3 {
        top = len(members)
    }
    topIds := make([]string, top)
    for i := 0; i < top; i++ {
        topIds[i] = members[i].MemberID
    }

    return &communitypb.GetTopCommunityMembersResponse{
        MemberIds: topIds,
    }, nil
}

func (h *CommunityHandler) GetUserRoleInCommunity(
	ctx context.Context,
	req *communitypb.GetUserRoleInCommunityRequest,
) (*communitypb.GetUserRoleInCommunityResponse, error) {
	role, err := h.svc.GetUserRoleInCommunity(ctx, req.CommunityId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user role: %v", err)
	}
	return &communitypb.GetUserRoleInCommunityResponse{Role: role}, nil
}

func (h *CommunityHandler) GetCommunityMemberCount(
	ctx context.Context,
	req *communitypb.GetCommunityByIDRequest,
) (*communitypb.CountResponse, error) {
	count, err := h.svc.GetCommunityMemberCount(ctx, req.CommunityId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count community members: %v", err)
	}
	return &communitypb.CountResponse{Count: int32(count)}, nil
}

func (h *CommunityHandler) ExploreCommunityByName(
    ctx context.Context, 
    req *communitypb.ExploreCommunityByNameRequest,
) (*communitypb.ExploreCommunityByNameResponse, error) {
    communities, total, err := h.svc.ExploreCommunityByName(
        req.Query, int(req.Threshold), int(req.Page), int(req.Size),
    )
    if err != nil {
        return nil, status.Errorf(codes.Internal, "explore failed: %v", err)
    }

    var pbCommunities []*communitypb.Community
    for _, c := range communities {
        pbCommunities = append(pbCommunities, &communitypb.Community{
            CommunityId:          c.CommunityID.String(),
            CommunityName:        c.CommunityName,
            CommunityDescription: c.CommunityDescription,
            CommunityRules:       c.CommunityRules,
            Status:               c.Status,
            CommunityLogo:        c.CommunityLogo,
            CommunityBanner:      c.CommunityBanner,
        })
    }
    return &communitypb.ExploreCommunityByNameResponse{
        Communities: pbCommunities,
        Total:       int32(total),
    }, nil
}