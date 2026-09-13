package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
)

func UploadHandler(client mediapb.MediaServiceClient) gin.HandlerFunc {
  return func(c *gin.Context) {
    fh, err := c.FormFile("file")
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
      return
    }
    f, err := fh.Open()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    filename := c.PostForm("filename")
    if filename == "" {
      filename = fh.Filename
    }

    grpcReq := &mediapb.MediaUploadRequest{
      Filename: filename,
      File:     data,
    }
    grpcResp, err := client.Upload(c.Request.Context(), grpcReq)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, gin.H{
      "id":         grpcResp.Id,
      "public_url": grpcResp.PublicUrl,
      "extension" : grpcResp.Extension,
      "created_at": grpcResp.CreatedAt,
    })
  }
}

func GetMediaByIDHandler(client mediapb.MediaServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Id string `json:"id" binding:"required,uuid"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		grpcReq := &mediapb.GetMediaByIdRequest{Id: req.Id}
		grpcResp, err := client.GetMediaById(c.Request.Context(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          grpcResp.Id,
			"public_url":  grpcResp.PublicUrl,
      "extension" : grpcResp.Extension,
			"created_at":  grpcResp.CreatedAt,
		})
	}
}