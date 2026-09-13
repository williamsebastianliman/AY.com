package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	authpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/authpb/proto/auth"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/database"
	authgrpc "github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/delivery/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/auth/internal/service"
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
	
	var jwtSecret = "p9sD#7vZ!qE3rC@X1tL$zM4nA&bUoJ8w"
	var accessTTL = 15 * time.Minute

	redisURL := os.Getenv("REDIS_URL")
	cacheSvc, err := service.NewRedisCache(redisURL)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	mailer := service.NewSMTPMailer()

	mqURL := os.Getenv("MQ_URL")
	bus, err := service.NewRabbitEventBus(mqURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	userClient, userConn, err := createUserClient()
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	authRepo := repository.NewRefreshTokenRepository(db)
	authSvc := service.NewAuthService(cacheSvc, mailer, bus, userClient, authRepo, jwtSecret,accessTTL)
	authHandler := authgrpc.NewAuthHandler(authSvc)

	secRepo := repository.NewSecurityAnswerRepository(db)
	secSvc := service.NewSecurityAnswerService(secRepo)
	secHandler := authgrpc.NewSecurityAnswerHandler(secSvc, userClient)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)
	authpb.RegisterSecurityAnswerServiceServer(grpcServer, secHandler)

	log.Printf("AuthService gRPC listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}