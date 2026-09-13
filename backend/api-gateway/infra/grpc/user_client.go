package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserServiceClient() (userpb.UserServiceClient, *grpc.ClientConn, error){
	target := os.Getenv("USER_SVC_ADDR")
		if target == "" {
		return nil, nil, errors.New("USER_SVC_ADDR must be set")
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

	client := userpb.NewUserServiceClient(conn)
	return client, conn, nil
}