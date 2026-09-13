package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewMediaServiceClient() (mediapb.MediaServiceClient, *grpc.ClientConn, error){
	target := os.Getenv("MEDIA_SVC_ADDR")
		if target == "" {
		return nil, nil, errors.New("MEDIA_SVC_ADDR must be set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
		if err != nil {
		return nil, nil, err
	}

	client := mediapb.NewMediaServiceClient(conn)
	return client, conn, nil
}