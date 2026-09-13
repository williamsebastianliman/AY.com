package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
)

type MyClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}


func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims := &MyClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil {
			var verr *jwt.ValidationError
			if ok := errors.As(err, &verr); ok {
				switch {
				case verr.Errors&jwt.ValidationErrorMalformed != 0:
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "malformed token"})
					return
				case verr.Errors&jwt.ValidationErrorExpired != 0:
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				default:
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			return
		}
		now := time.Now()
		if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not valid yet"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func ThreadOwnerOnly(threadClient threadpb.ThreadServiceClient, secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
		log.Printf("test! mnasuk");
        tokenStr, err := c.Cookie("access_token")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        claims := &MyClaims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(secret), nil
        })

        if err != nil || !token.Valid || claims.UserID == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        threadID := c.Param("id")
        if threadID == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "thread id required"})
            return
        }

        resp, err := threadClient.GetThreadByID(
            context.Background(),
            &threadpb.GetThreadByIDRequest{Id: threadID},
        )
        if err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "thread not found"})
            return
        }
        threadUserID := resp.Thread.GetUserId()
        if threadUserID != claims.UserID {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not the thread owner"})
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
	
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}