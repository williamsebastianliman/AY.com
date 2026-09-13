package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewThreadServiceClient() (threadpb.ThreadServiceClient, *grpc.ClientConn, error) {
	target := os.Getenv("THREAD_SVC_ADDR")
	if target == "" {
		return nil, nil, errors.New("THREAD_SVC_ADDR must be set")
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

	client := threadpb.NewThreadServiceClient(conn)
	return client, conn, nil
}