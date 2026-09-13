package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	communitypb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/communitypb/proto/community"
	"google.golang.org/grpc/status"
)

type CommunityJSON struct {
	CommunityId          string   `json:"community_id"`
	CommunityName        string   `json:"community_name"`
	CommunityDescription string   `json:"community_description"`
	CommunityRules       string   `json:"community_rules"`
	Status               string   `json:"status"`
	CreatedAt            TimeJSON `json:"created_at"`
	CommunityLogo        string   `json:"community_logo"`
	CommunityBanner      string   `json:"community_banner"`
}

type CommunityCategoryJSON struct {
	CategoryId   string   `json:"category_id"`
	CategoryName string   `json:"category_name"`
	CreatedAt    TimeJSON `json:"created_at"`
}

type CommunityMemberJSON struct {
	CommunityId string   `json:"community_id"`
	MemberId    string   `json:"member_id"`
	Status      string   `json:"status"`
	CreatedAt   TimeJSON `json:"created_at"`
}

func toTimeJSON(pb *communitypb.Timestamp) TimeJSON {
	return TimeJSON{
		Year:     pb.GetYear(),
		Month:    pb.GetMonth(),
		Day:      pb.GetDay(),
		Hour:     pb.GetHour(),
		Minute:   pb.GetMinute(),
		Second:   pb.GetSecond(),
		Timezone: pb.GetTimezone(),
	}
}

func toCommunityJSON(pb *communitypb.Community) CommunityJSON {
	return CommunityJSON{
		CommunityId:          pb.GetCommunityId(),
		CommunityName:        pb.GetCommunityName(),
		CommunityDescription: pb.GetCommunityDescription(),
		CommunityRules:       pb.GetCommunityRules(),
		Status:               pb.GetStatus(),
		CreatedAt:            toTimeJSON(pb.GetCreatedAt()),
		CommunityLogo:        pb.GetCommunityLogo(),
		CommunityBanner:      pb.GetCommunityBanner(),
	}
}

func toCategoryJSON(pb *communitypb.CommunityCategory) CommunityCategoryJSON {
	return CommunityCategoryJSON{
		CategoryId:   pb.GetCategoryId(),
		CategoryName: pb.GetCategoryName(),
		CreatedAt:    toTimeJSON(pb.GetCreatedAt()),
	}
}

func toMemberJSON(pb *communitypb.CommunityMember) CommunityMemberJSON {
	return CommunityMemberJSON{
		CommunityId: pb.GetCommunityId(),
		MemberId:    pb.GetMemberId(),
		Status:      pb.GetStatus(),
		CreatedAt:   toTimeJSON(pb.GetCreatedAt()),
	}
}

func ListCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		req := &communitypb.ListCommunitiesRequest{
			Filter:    c.Query("query"),
			CategoryId: c.Query("category"),
			UserId: c.Query("user_id"),
			Page:      int32(page),
			Size:      int32(size),
		}
		resp, err := client.ListCommunities(c.Request.Context(), req)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, comm := range resp.Communities {
			communities = append(communities, toCommunityJSON(comm))
		}
		c.JSON(http.StatusOK, communities)
	}
}

func SearchCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		req := &communitypb.SearchCommunitiesRequest{
			Keyword: c.Query("keyword"),
			Page:    int32(page),
			Size:    int32(size),
		}
		resp, err := client.SearchCommunities(c.Request.Context(), req)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, comm := range resp.Communities {
			communities = append(communities, toCommunityJSON(comm))
		}
		c.JSON(http.StatusOK, communities)
	}
}

func GetCommunityByID(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		resp, err := client.GetCommunityByID(c.Request.Context(), &communitypb.GetCommunityByIDRequest{CommunityId: id})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, toCommunityJSON(resp.Community))
	}
}

func CreateCommunity(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CommunityName        string   `json:"community_name" binding:"required"`
			CommunityDescription string   `json:"community_description"`
			CommunityRules       string   `json:"community_rules"`
			Status               string   `json:"status"`
			CommunityLogo        string   `json:"community_logo"`
			CommunityBanner      string   `json:"community_banner"`
			CategoryIds          []string `json:"category_ids"`
			OwnerId string `json:"owner_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.CreateCommunity(c.Request.Context(), &communitypb.CreateCommunityRequest{
			CommunityName:        req.CommunityName,
			CommunityDescription: req.CommunityDescription,
			CommunityRules:       req.CommunityRules,
			Status:               req.Status,
			CommunityLogo:        req.CommunityLogo,
			CommunityBanner:      req.CommunityBanner,
			CategoryIds:          req.CategoryIds,
			OwnerId: req.OwnerId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"community_id": resp.CommunityId})
	}
}

func GetTopCommunityMembers(client communitypb.CommunityServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        communityId := c.Param("community_id")
        resp, err := client.GetTopCommunityMembers(
            c.Request.Context(),
            &communitypb.GetTopCommunityMembersRequest{
                CommunityId: communityId,
            },
        )
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"member_ids": resp.MemberIds})
    }
}

func RequestJoinCommunity(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserId string `json:"user_id" binding:"required"`
		}
		communityId := c.Param("community_id")
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.RequestJoinCommunity(c.Request.Context(), &communitypb.RequestJoinCommunityRequest{
			CommunityId: communityId,
			UserId:      req.UserId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func ListJoinedCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		pageStr := c.Query("page")
		sizeStr := c.Query("size")

		page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
            return
        }
        size, err := strconv.Atoi(sizeStr)
        if err != nil || size < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
            return
        }
		resp, err := client.ListJoinedCommunities(c.Request.Context(), &communitypb.ListJoinedCommunitiesRequest{UserId: userId, Page: int32(page), Size: int32(size)})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, comm := range resp.Communities {
			communities = append(communities, toCommunityJSON(comm))
		}
		c.JSON(http.StatusOK, communities)
	}
}

func ListRequestedCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		pageStr := c.DefaultQuery("page","0")
		sizeStr := c.DefaultQuery("size", "10")

		page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
            return
        }
        size, err := strconv.Atoi(sizeStr)
        if err != nil || size < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
            return
        }

		resp, err := client.ListRequestedCommunities(c.Request.Context(), &communitypb.ListRequestedCommunitiesRequest{UserId: userId, Page: int32(page), Size: int32(size)})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, comm := range resp.Communities {
			communities = append(communities, toCommunityJSON(comm))
		}
		c.JSON(http.StatusOK, communities)
	}
}

func ListModeratedCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("user_id")
		resp, err := client.ListModeratedCommunities(c.Request.Context(), &communitypb.ListModeratedCommunitiesRequest{UserId: userId})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, comm := range resp.Communities {
			communities = append(communities, toCommunityJSON(comm))
		}
		c.JSON(http.StatusOK, communities)
	}
}

func ListCommunityMembers(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		role := c.Query("role")
		pageStr := c.DefaultQuery("page","1")
		sizeStr := c.DefaultQuery("size","2000")
		search := c.DefaultQuery("search","")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			size = 2000
		}
		resp, err := client.ListCommunityMembers(c.Request.Context(), &communitypb.ListCommunityMembersRequest{CommunityId: communityId, Role: role, Page: int32(page), Size: int32(size), Name: search})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		members := make([]CommunityMemberJSON, 0, len(resp.Members))
		for _, m := range resp.Members {
			members = append(members, toMemberJSON(m))
		}
		c.JSON(http.StatusOK, members)
	}
}

func PromoteMember(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		userId := c.Param("user_id")
		_, err := client.PromoteMember(c.Request.Context(), &communitypb.PromoteMemberRequest{
			CommunityId: communityId,
			UserId:      userId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DemoteMember(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		userId := c.Param("user_id")
		_, err := client.DemoteMember(c.Request.Context(), &communitypb.DemoteMemberRequest{
			CommunityId: communityId,
			UserId:      userId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func AcceptJoinRequest(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		userId := c.Param("user_id")
		_, err := client.AcceptJoinRequest(c.Request.Context(), &communitypb.AcceptJoinRequestRequest{
			CommunityId: communityId,
			UserId:      userId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RejectJoinRequest(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		userId := c.Param("user_id")
		_, err := client.RejectJoinRequest(c.Request.Context(), &communitypb.RejectJoinRequestRequest{
			CommunityId: communityId,
			UserId:      userId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CommunityAbout(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		resp, err := client.CommunityAbout(c.Request.Context(), &communitypb.CommunityAboutRequest{CommunityId: communityId})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, toCommunityJSON(resp.Community))
	}
}

func ListCommunityCategories(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		resp, err := client.ListCommunityCategories(c.Request.Context(), &communitypb.ListCommunityCategoriesRequest{CommunityId: communityId})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		cats := make([]CommunityCategoryJSON, 0, len(resp.Categories))
		for _, cat := range resp.Categories {
			cats = append(cats, toCategoryJSON(cat))
		}
		c.JSON(http.StatusOK, cats)
	}
}

func AddCommunityCategory(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		categoryId := c.Param("category_id")
		_, err := client.AddCommunityCategory(c.Request.Context(), &communitypb.AddCommunityCategoryRequest{
			CommunityId: communityId,
			CategoryId:  categoryId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func RemoveCommunityCategory(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		categoryId := c.Param("category_id")
		_, err := client.RemoveCommunityCategory(c.Request.Context(), &communitypb.RemoveCommunityCategoryRequest{
			CommunityId: communityId,
			CategoryId:  categoryId,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func ListCategories(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.ListCategories(c.Request.Context(), &communitypb.Empty{})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		cats := make([]CommunityCategoryJSON, 0, len(resp.Categories))
		for _, cat := range resp.Categories {
			cats = append(cats, toCategoryJSON(cat))
		}
		c.JSON(http.StatusOK, cats)
	}
}

func CreateCategory(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			CategoryName string `json:"category_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.CreateCategory(c.Request.Context(), &communitypb.CreateCategoryRequest{CategoryName: req.CategoryName})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"category_id": resp.CategoryId})
	}
}

func UpdateCategory(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId := c.Param("category_id")
		var req struct {
			CategoryName string `json:"category_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := client.UpdateCategory(c.Request.Context(), &communitypb.UpdateCategoryRequest{
			CategoryId:   categoryId,
			CategoryName: req.CategoryName,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DeleteCategory(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId := c.Param("category_id")
		_, err := client.DeleteCategory(c.Request.Context(), &communitypb.DeleteCategoryRequest{CategoryId: categoryId})
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func ListPendingCommunities(client communitypb.CommunityServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        resp, err := client.ListPendingCommunities(c.Request.Context(), &communitypb.Empty{})
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        communities := make([]CommunityJSON, 0, len(resp.Communities))
        for _, comm := range resp.Communities {
            communities = append(communities, toCommunityJSON(comm))
        }
        c.JSON(http.StatusOK, communities)
    }
}
func AcceptCommunityRequest(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		resp, err := client.AcceptCommunityCreation(
			c.Request.Context(),
			&communitypb.AcceptRejectCommunityRequest{CommunityId: communityId},
		)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": resp.Message})
	}
}

func RejectCommunityRequest(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		resp, err := client.RejectCommunityCreation(
			c.Request.Context(),
			&communitypb.AcceptRejectCommunityRequest{CommunityId: communityId},
		)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": resp.Message})
	}
	
}

func GetUserRoleInCommunity(client communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		userId := c.Param("user_id")

		resp, err := client.GetUserRoleInCommunity(
			c.Request.Context(),
			&communitypb.GetUserRoleInCommunityRequest{
				CommunityId: communityId,
				UserId:      userId,
			},
		)
		if err != nil {
			st := status.Convert(err)
			c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"community_id": communityId,
			"user_id":      userId,
			"role":         resp.Role,
		})
	}
}

func GetCommunityMemberCountHandler(communityClient communitypb.CommunityServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		communityId := c.Param("community_id")
		resp, err := communityClient.GetCommunityMemberCount(
			c,
			&communitypb.GetCommunityByIDRequest{CommunityId: communityId},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get community member count",
				"detail": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"community_id": communityId,
			"count":        resp.Count,
		})
	}
}

func ExploreCommunityByNameHandler(client communitypb.CommunityServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        query := c.Query("query")
        threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "2"))
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

        req := &communitypb.ExploreCommunityByNameRequest{
            Query:     query,
            Threshold: int32(threshold),
            Page:      int32(page),
            Size:      int32(size),
        }

        resp, err := client.ExploreCommunityByName(context.Background(), req)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }

        communities := make([]CommunityJSON, 0, len(resp.Communities))
		for _, com := range resp.Communities {
			communities = append(communities, toCommunityJSON(com))
		}
		c.JSON(http.StatusOK, communities)
    }
}