package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ActivateAccount(client userpb.UserServiceClient) gin.HandlerFunc{
	return func(c *gin.Context){
		var req struct{
			Id string `json:"id" binding:"required"`
		}
		if err:= c.ShouldBindJSON(&req); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.GetUserByIdRequest{
			Id: req.Id,
		}
		grpcResp, err:= client.ActivateAccount(c, grpcReq)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"id":grpcResp.Id})

	}
} 
func RegisterUser(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name             string `json:"name" binding:"required"`
			Username         string `json:"username" binding:"required"`
			Email            string `json:"email" binding:"required,email"`
			Password         string `json:"password" binding:"required"`
			Gender           string `json:"gender" binding:"required"`
			BirthYear        string `json:"birth_year" binding:"required"`
			BirthMonth       string `json:"birth_month" binding:"required"`
			BirthDay         string `json:"birth_day" binding:"required"`
			SubscribedNews   bool   `json:"subscribed_news" binding:"required"`
			ProfilePictureId string `json:"profile_picture_id" binding:"required"`
			BannerMediaId    string `json:"banner_media_id" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &userpb.RegisterUserRequest{
			Name:             req.Name,
			Username:         req.Username,
			Email:            req.Email,
			Password:         req.Password,
			Gender:           req.Gender,
			BirthYear:        req.BirthYear,
			BirthMonth:       req.BirthMonth,
			BirthDay:         req.BirthDay,
			SubscribedNews:   req.SubscribedNews,
			ProfilePictureId: req.ProfilePictureId,
			BannerMediaId:    req.BannerMediaId,
		}
		grpcResp, err := client.RegisterUser(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": grpcResp.Id})
	}
}

func IsUsernameUnique(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.UsernameValidationRequest{Username: req.Username}
		grpcResp, err := client.IsUsernameUnique(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"is_unique": grpcResp.IsUnique})
	}
}

func GetUserByUsernameHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &userpb.GetUserByUsernameRequest{
			Username: req.Username,
		}
		grpcResp, err := client.GetUserByUsername(c.Request.Context(), grpcReq)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                 grpcResp.Id,
			"name":               grpcResp.Name,
			"username":           grpcResp.Username,
			"email":              grpcResp.Email,
			"gender":             grpcResp.Gender,
			"birth_year":         grpcResp.BirthYear,
			"birth_month":        grpcResp.BirthMonth,
			"birth_day":          grpcResp.BirthDay,
			"subscribed_news":    grpcResp.SubscribedNews,
			"profile_picture_id": grpcResp.ProfilePictureId,
			"banner_media_id":    grpcResp.BannerMediaId,
		})
	}
}

func GetUserByEmail(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &userpb.GetUserByUsernameRequest{
			Username: req.Username,
		}
		grpcResp, err := client.GetUserByEmail(c.Request.Context(), grpcReq)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if grpcResp.IsBanned{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User Is Banned"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                 grpcResp.Id,
			"name":               grpcResp.Name,
			"username":           grpcResp.Username,
			"email":              grpcResp.Email,
			"gender":             grpcResp.Gender,
			"birth_year":         grpcResp.BirthYear,
			"birth_month":        grpcResp.BirthMonth,
			"birth_day":          grpcResp.BirthDay,
			"subscribed_news":    grpcResp.SubscribedNews,
			"profile_picture_id": grpcResp.ProfilePictureId,
			"banner_media_id":    grpcResp.BannerMediaId,
		})
	}
}

func GetMeHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		md := metadata.New(map[string]string{"user_id": userID.(string)})
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		grpcResp, err := client.GetMe(ctx, &userpb.GetMeRequest{})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":                 grpcResp.Id,
			"name":               grpcResp.Name,
			"username":           grpcResp.Username,
			"email":              grpcResp.Email,
			"gender":             grpcResp.Gender,
			"birth_year":         grpcResp.BirthYear,
			"birth_month":        grpcResp.BirthMonth,
			"birth_day":          grpcResp.BirthDay,
			"subscribed_news":    grpcResp.SubscribedNews,
			"profile_picture_id": grpcResp.ProfilePictureId,
			"banner_media_id":    grpcResp.BannerMediaId,
		})
	}
}

func GetUserByIDHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
            return
        }

        grpcResp, err := client.GetUserById(c.Request.Context(), &userpb.GetUserByIdRequest{Id: id})
        if err != nil {
            st := status.Convert(err)
            switch st.Code() {
            case codes.NotFound:
                c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            }
            return
        }
		log.Printf("logzz gateway: %t", grpcResp.IsPremium);
        c.JSON(http.StatusOK, gin.H{
            "id":                  grpcResp.Id,
			"is_premium":          grpcResp.IsPremium,
            "name":                grpcResp.Name,
            "username":            grpcResp.Username,
            "email":               grpcResp.Email,
            "gender":              grpcResp.Gender,
            "birth_year":          grpcResp.BirthYear,
            "birth_month":         grpcResp.BirthMonth,
            "birth_day":           grpcResp.BirthDay,
            "subscribed_news":     grpcResp.SubscribedNews,
            "profile_picture_id":  grpcResp.ProfilePictureId,
            "banner_media_id":     grpcResp.BannerMediaId,
			"password" : grpcResp.Password,
        })
    }
}

func FollowUserHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UserID     string `json:"user_id" binding:"required"`
            FollowerID string `json:"follower_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &userpb.FollowUserRequest{
            UserId:     req.UserID,
            FollowerId: req.FollowerID,
        }
        grpcResp, err := client.FollowUser(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
    }
}
func UnfollowUserHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UserID     string `json:"user_id" binding:"required"`
            FollowerID string `json:"follower_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &userpb.UnfollowUserRequest{
            UserId:     req.UserID,
            FollowerId: req.FollowerID,
        }
        grpcResp, err := client.UnfollowUser(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
    }
}
func GetFollowersHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }
        grpcReq := &userpb.GetFollowersRequest{
            UserId: userID,
        }
        grpcResp, err := client.GetFollowers(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"followers": grpcResp.Followers})
    }
}

func GetFollowingHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }
        grpcReq := &userpb.GetFollowingRequest{
            UserId: userID,
        }
        grpcResp, err := client.GetFollowing(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"following": grpcResp.Following})
    }
}

func FollowerCountHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }
        grpcReq := &userpb.GetUserByIdRequest{
            Id: userID,
        }
        grpcResp, err := client.FollowerCount(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"count": grpcResp.Count})
    }
}
func FollowingCountHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Param("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
            return
        }
        grpcReq := &userpb.GetUserByIdRequest{
            Id: userID,
        }
        grpcResp, err := client.FollowingCount(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"count": grpcResp.Count})
    }
}
func IsFollowedHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UserID     string `json:"user_id" binding:"required"`
            FollowerID string `json:"follower_id" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        grpcReq := &userpb.GetFollowStatusRequest{
            UserId:     req.UserID,
            FollowerId: req.FollowerID,
        }
        grpcResp, err := client.IsFollowed(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"is_followed": grpcResp.IsFollowed})
    }
}

func BlockUserHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID   string `json:"user_id" binding:"required"`
			TargetID string `json:"target_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.BlockUserRequest{
			UserId:   req.UserID,
			TargetId: req.TargetID,
		}
		grpcResp, err := client.BlockUser(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
	}
}

func UnblockUserHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID   string `json:"user_id" binding:"required"`
			TargetID string `json:"target_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.UnblockUserRequest{
			UserId:   req.UserID,
			TargetId: req.TargetID,
		}
		grpcResp, err := client.UnblockUser(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
	}
}

func IsUserBlockedHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID   string `json:"user_id" binding:"required"`
			TargetID string `json:"target_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.IsUserBlockedRequest{
			UserId:   req.UserID,
			TargetId: req.TargetID,
		}
		grpcResp, err := client.IsUserBlocked(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"is_blocked": grpcResp.IsBlocked})
	}
}

func RequestPremiumHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID     string `json:"user_id" binding:"required"`
			CardNumber string `json:"card_number" binding:"required"`
			Reason     string `json:"reason" binding:"required"`
			FaceImage  string `json:"face_image" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.RequestPremiumRequest{
			UserId:     req.UserID,
			CardNumber: req.CardNumber,
			Reason:     req.Reason,
			FaceImage:  req.FaceImage,
		}
		grpcResp, err := client.RequestPremium(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
	}
}




func IsPremium(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID   string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.IsPremiumRequest{
			UserId:   req.UserID,
		}
		grpcResp, err := client.IsPremium(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"is_premium": grpcResp.Success})
	}
}

func GetAllPremiumRequestsHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcResp, err := client.GetAllPremiumRequests(c.Request.Context(), &userpb.GetAllPremiumRequestsRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"requests": grpcResp.Requests})
	}
}

func AcceptPremiumHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.AcceptPremiumRequest{UserId: req.UserID}
		grpcResp, err := client.AcceptPremium(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": grpcResp.Success,
		})
	}
}

func RejectPremiumHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		grpcReq := &userpb.RejectPremiumRequest{UserId: req.UserID}
		grpcResp, err := client.RejectPremium(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": grpcResp.Success,
		})
	}
}

func UpdateProfileHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            UserID           string `json:"user_id" binding:"required"`
            Name             string `json:"name" binding:"required"`
            Username         string `json:"username" binding:"required"`
            Email            string `json:"email" binding:"required,email"`
            Password         string `json:"password"`
            Gender           string `json:"gender" binding:"required"`
            BirthYear        string `json:"birth_year" binding:"required"`
            BirthMonth       string `json:"birth_month" binding:"required"`
            BirthDay         string `json:"birth_day" binding:"required"`
            ProfilePictureId string `json:"profile_picture_id"`
            BannerMediaId    string `json:"banner_media_id"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        grpcReq := &userpb.UpdateProfileRequest{
            UserId:           req.UserID,
            Name:             req.Name,
            Username:         req.Username,
            Email:            req.Email,
            Password:         req.Password,
            Gender:           req.Gender,
            BirthYear:        req.BirthYear,
            BirthMonth:       req.BirthMonth,
            BirthDay:         req.BirthDay,
            ProfilePictureId: req.ProfilePictureId,
            BannerMediaId:    req.BannerMediaId,
        }

        grpcResp, err := client.UpdateProfile(c.Request.Context(), grpcReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
    }
}

func GetAllUsers(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        resp, err := client.GetAllUsers(c.Request.Context(), &userpb.GetAllUsersRequest{})
        if err != nil {
            st := status.Convert(err)
            c.JSON(codesToHTTP(st.Code()), gin.H{"error": st.Message()})
            return
        }
        type UserResp struct {
            Id       string `json:"id"`
            Username string `json:"username"`
            Name     string `json:"name"`
            ProfilePictureId string `json:"profile_picture_id"`
        }
        users := make([]UserResp, 0, len(resp.Users))
        for _, u := range resp.Users {
            users = append(users, UserResp{
                Id:       u.Id,
                Username: u.Username,
                Name:     u.Name,
                ProfilePictureId: u.ProfilePictureId,
            })
        }
        c.JSON(http.StatusOK, users)
    }
}

func CreateUserReportHandler(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ReportedUser string `json:"reported_user" binding:"required"`
			Reason       string `json:"reason" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.CreateUserReport(
			c.Request.Context(),
			&userpb.CreateUserReportRequest{
				ReportedUser: req.ReportedUser,
				Reason:       req.Reason,
			},
		)
		if err != nil {
			st := status.Convert(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"report_id": resp.ReportId})
	}
}

func GetAllUserReportsHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        resp, err := client.GetAllUserReports(
            c.Request.Context(),
            &userpb.GetAllUserReportsRequest{},
        )
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }
        reports := make([]gin.H, 0, len(resp.Reports))
        for _, r := range resp.Reports {
            reports = append(reports, gin.H{
                "report_id":     r.ReportId,
                "reported_user": r.ReportedUser,
                "reason":        r.Reason,
            })
        }
        c.JSON(http.StatusOK, reports)
    }
}

func ExploreUserByNameHandler(client userpb.UserServiceClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        query := c.Query("query")
        threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "2"))
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

        req := &userpb.ExploreUserByNameRequest{
            Query:     query,
            Threshold: int32(threshold),
            Page:      int32(page),
            Size:      int32(size),
        }

        resp, err := client.ExploreUserByName(c.Request.Context(), req)
        if err != nil {
            st := status.Convert(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
            return
        }

        users := make([]gin.H, 0, len(resp.Users))
        for _, u := range resp.Users {
            users = append(users, gin.H{
                "id":                 u.Id,
                "username":           u.Username,
                "name":               u.Name,
                "profile_picture_id": u.ProfilePictureId,
            })
        }
        c.JSON(http.StatusOK, users)
    }
}
