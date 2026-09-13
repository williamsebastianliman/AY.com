package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	communitypb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/communitypb/proto/community"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCommunityServiceClient() (communitypb.CommunityServiceClient, *grpc.ClientConn, error){
	target := os.Getenv("COMMUNITY_SVC_ADDR")
		if target == "" {
		return nil, nil, errors.New("COMMUNITY_SVC_ADDR must be set")
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

	client := communitypb.NewCommunityServiceClient(conn)
	return client, conn, nil
}