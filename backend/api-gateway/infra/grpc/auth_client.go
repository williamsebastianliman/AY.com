package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth            authpb.AuthServiceClient
	SecurityAnswer authpb.SecurityAnswerServiceClient
	Conn            *grpc.ClientConn
}

func NewAuthClients() (*Clients, error) {
	target := os.Getenv("AUTH_SVC_ADDR")
	if target == "" {
		return nil, errors.New("AUTH_SVC_ADDR must be set")
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
		return nil, err
	}

	return &Clients{
		Auth:           authpb.NewAuthServiceClient(conn),
		SecurityAnswer: authpb.NewSecurityAnswerServiceClient(conn),
		Conn:           conn,
	}, nil
}