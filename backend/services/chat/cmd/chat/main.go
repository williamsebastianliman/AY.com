package main

import (
	"log"
	"net"
	"os"

	chatpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/chatpb/proto/chat"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/database"
	chatgrpc "github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/delivery/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/service"
	"google.golang.org/grpc"
)

func main() {
	db := database.InitDB()
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	chatRepo := repository.NewChatRepository(db)
	chatSvc := service.NewChatService(chatRepo)
	communityHandler := chatgrpc.NewChatHandler(chatSvc)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, communityHandler)

	log.Printf("CommunityService gRPC listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}