package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
)

func SendOTPHandler(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id string `json:"id" binding:"required"`
			Email string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &authpb.SendOTPRequest{
			Email: req.Email,
			UserId: req.Id,
		}

		grpcResp, err := client.SendOTP(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"request_id": grpcResp.RequestId})
	}
}

func VerifyOTPHandler(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RequestID string `json:"request_id" binding:"required,uuid"`
			Code      string `json:"code"       binding:"required,len=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &authpb.VerifyOTPRequest{
			RequestId: req.RequestID,
			Code:      req.Code,
		}

		grpcResp, err := client.VerifyOTP(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": grpcResp.Success})
	}
}

func LoginHandler(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &authpb.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		}
		grpcResp, err := client.Login(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    grpcResp.AccessToken,
			Path:     "/",
			Domain:   "localhost",           
			MaxAge:   15 * 60,               
			Secure:   false,               
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    grpcResp.RefreshToken,
			Path:     "/",
			Domain:   "localhost",           
			MaxAge:   3 * 24 * 3600,   
			Secure:   false,                 
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func CreateSecurityAnswerHandler(client authpb.SecurityAnswerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID   string `json:"user_id"   binding:"required,uuid"`
			Question string `json:"question"  binding:"required"`
			Answer   string `json:"answer"    binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &authpb.CreateSecurityAnswerRequest{
			UserId:   req.UserID,
			Question: req.Question,
			Answer:   req.Answer,
		}

		grpcResp, err := client.CreateSecurityAnswer(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"request_id": grpcResp.RequestId})
	}
}

func GetSecurityAnswersHandler(client authpb.SecurityAnswerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id" binding:"required,uuid"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stream, err := client.GetSecurityAnswers(c.Request.Context(), &authpb.GetSecurityAnswerRequest{
			UserId: req.UserID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var answers []struct {
			ID         string `json:"id"`
			UserID     string `json:"user_id"`
			Question   string `json:"question"`
			AnswerHash string `json:"answer_hash"`
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			answers = append(answers, struct {
				ID         string `json:"id"`
				UserID     string `json:"user_id"`
				Question   string `json:"question"`
				AnswerHash string `json:"answer_hash"`
			}{
				ID:         resp.Id,
				UserID:     resp.UserId,
				Question:   resp.Question,
				AnswerHash: resp.AnswerHash,
			})
		}

		c.JSON(http.StatusOK, gin.H{"answers": answers})
	}
}

func RefreshHandler(client authpb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		rt, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
			return
		}

		grpcReq := &authpb.RefreshTokenRequest{RefreshToken: rt}
		grpcResp, err := client.RefreshToken(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "access_token",
			Value:    grpcResp.AccessToken,
			Path:     "/",
			Domain:   "localhost",       
			MaxAge:   15 * 60,          
			Secure:   false,            
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if grpcResp.RefreshToken != "" && grpcResp.RefreshToken != rt {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "refresh_token",
				Value:    grpcResp.RefreshToken,
				Path:     "/",
				Domain:   "localhost",
				MaxAge:   3 * 24 * 3600,
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func CheckSecurityAnswerHandler(client authpb.SecurityAnswerServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UserID      string `json:"user_id" binding:"required,uuid"`
			Question    string `json:"question" binding:"required"`
			Answer      string `json:"answer" binding:"required"`
			NewPassword string `json:"new_password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &authpb.CheckSecurityAnswerRequest{
			UserId:      req.UserID,
			Question:    req.Question,
			Answer:      req.Answer,
			NewPassword: req.NewPassword,
		}

		grpcResp, err := client.CheckSecurityAnswer(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":       grpcResp.Success,
			"error_message": grpcResp.ErrorMessage,
		})
	}
}