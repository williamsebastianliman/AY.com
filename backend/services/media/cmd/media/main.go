package main

import (
	"log"
	"net"

	mediapb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/mediapb/proto/media"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/database"
	mediagrpc "github.com/williamsebastianliman/WEB-WS-242/services/media/internal/delivery/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/media/internal/service"
	"google.golang.org/grpc"
)

func main() {
	db := database.InitDB()
	
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := database.Seeder(db); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	repo := repository.NewMediaRepository(db)

	supabaseURL := "https://itmuqguefzemfzsbryad.supabase.co/storage/v1"
	serviceKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Iml0bXVxZ3VlZnplbWZ6c2JyeWFkIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc0NTU2MjM0NCwiZXhwIjoyMDYxMTM4MzQ0fQ.f6THbFc5nUUp-ofyuWG8ga506t58Y6saVDo8QWDlku8"
	bucket := "avatars"

	mediaSvc := service.NewMediaService(repo, supabaseURL, serviceKey, bucket)
	handler := mediagrpc.NewMediaHandler(mediaSvc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	mediapb.RegisterMediaServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}