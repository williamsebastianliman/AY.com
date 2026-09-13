package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	communitypb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/communitypb/proto/community"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/database"
	communitygrpc "github.com/williamsebastianliman/WEB-WS-242/services/community/internal/delivery/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/community/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createUserClient() (userpb.UserServiceClient, *grpc.ClientConn, error) {
	target := os.Getenv("USER_SVC_ADDR")
	if target == "" {
		target = "user:50051"
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

func main() {
	db := database.InitDB()
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := database.Seeder(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}
	userClient, userConn, err := createUserClient()
	if err != nil{
		log.Fatalf("failed to connect to media service: %v", err)
	}
	defer userConn.Close()
	communityRepo := repository.NewCommunityRepository(db)
	communitySvc := service.NewCommunityService(communityRepo)
	communityHandler := communitygrpc.NewCommunityHandler(communitySvc, userClient)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	communitypb.RegisterCommunityServiceServer(grpcServer, communityHandler)

	log.Printf("CommunityService gRPC listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}