package grpc

import (
	"context"

	"github.com/google/uuid"
	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/service"
)

type MediaHandler struct {
	mediapb.UnimplementedMediaServiceServer
	svc service.MediaService
}

func NewMediaHandler(svc service.MediaService) *MediaHandler{
	return &MediaHandler{svc: svc}
}

func (h *MediaHandler) Upload(ctx context.Context, req *mediapb.MediaUploadRequest) (*mediapb.MediaUploadResponse, error){
	media, err := h.svc.Upload(ctx, req.File, req.Filename)
	if(err != nil){
		return nil, err
	}
	return &mediapb.MediaUploadResponse{Id: media.ID.String(), PublicUrl: media.PublicURL, CreatedAt: media.CreatedAt.String()}, nil
}

func (h *MediaHandler) GetMediaById(ctx context.Context, req *mediapb.GetMediaByIdRequest) (*mediapb.GetMediaByIdResponse, error){
	uid, err := uuid.Parse(req.Id)
	if err != nil{
		return nil, err
	}

	media, err := h.svc.GetMediaByID(ctx, uid)
	if err != nil{
		return nil, err
	}
	return &mediapb.GetMediaByIdResponse{Id: media.ID.String(), PublicUrl: media.PublicURL, Extension:media.Extension, CreatedAt: media.CreatedAt.String()}, nil
} 