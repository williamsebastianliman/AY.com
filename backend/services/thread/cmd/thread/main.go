package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	threadpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/threadpb/proto/thread"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/database"
	grpcSvr "github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/delivery/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/thread/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createMediaClient() (mediapb.MediaServiceClient, *grpc.ClientConn, error) {
	target := os.Getenv("MEDIA_SVC_ADDR")
	if target == "" {
		target = "media:50051"
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
		log.Fatalf("database migration failed: %v", err)
	}



	threadRepo := repository.NewThreadRepository(db)
	threadSvc  := service.NewThreadService(threadRepo)


	port := os.Getenv("THREAD_PORT")
	if port == "" {
		port = "50051"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on :%s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	mediaClient, mediaConn, err := createMediaClient()
	if err != nil{
		log.Fatalf("failed to connect to media service: %v", err)
	}
	defer mediaConn.Close()

	userClient, userConn, err := createUserClient()
	if err != nil{
		log.Fatalf("failed to connect to media service: %v", err)
	}
	defer userConn.Close()

	threadpb.RegisterThreadServiceServer(grpcServer, grpcSvr.NewThreadHandler(threadSvc, mediaClient, userClient))

	log.Printf("Thread gRPC server listening on :%s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}