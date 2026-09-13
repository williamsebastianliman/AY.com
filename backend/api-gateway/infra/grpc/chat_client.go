package grpc

import (
	"context"
	"errors"
	"os"
	"time"

	chatpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/chatpb/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewChatClients() (chatpb.ChatServiceClient, *grpc.ClientConn, error) {
	target := os.Getenv("CHAT_SVC_ADDR")
	if target == "" {
		return nil, nil, errors.New("CHAT_SVC_ADDR must be set")
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
		return nil,nil, err
	}
	client := chatpb.NewChatServiceClient(conn)
	return client, conn, nil
}