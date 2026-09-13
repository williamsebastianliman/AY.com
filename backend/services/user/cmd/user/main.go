package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/streadway/amqp"
	userpb "github.com/williamsebastianliman/WEB-WS-242/proto/gen/go/userpb/proto/user"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/database"
	event "github.com/williamsebastianliman/WEB-WS-242/services/user/internal/delivery/event"
	usergrpc "github.com/williamsebastianliman/WEB-WS-242/services/user/internal/delivery/grpc"
	eventBus "github.com/williamsebastianliman/WEB-WS-242/services/user/internal/event"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/repository"
	"github.com/williamsebastianliman/WEB-WS-242/services/user/internal/service"
	"google.golang.org/grpc"
)

func main() {
	amqpURL := os.Getenv("MQ_URL")
	if amqpURL == "" {
		log.Fatal("MQ_URL environment variable is required")
	}

	db := database.InitDB()
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	if err := database.Seeder(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ at %q: %v", amqpURL, err)
	}
	defer conn.Close()

	bus := eventBus.NewRabbitMQEventBus(conn)
	if err := event.StartEventConsumer(context.Background(), bus, userSvc); err != nil {
		log.Fatalf("failed to start event consumer: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, usergrpc.NewUserHandler(userSvc))

	log.Println("User gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}