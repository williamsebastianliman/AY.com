// @title My API
// @version 1.0
// @description This is my API.
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/williamsebastianliman/WEB-WS-242/api-gateway/docs"
	"github.com/williamsebastianliman/WEB-WS-242/api-gateway/infra/grpc"
	"github.com/williamsebastianliman/WEB-WS-242/api-gateway/internal/routes"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authClients, err := grpc.NewAuthClients()
	if err != nil {	
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authClients.Conn.Close()

	userClient, ucConn, err := grpc.NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
	}
	defer ucConn.Close()

	mediaClient, mdConn, err := grpc.NewMediaServiceClient()
	if err != nil {
		log.Fatalf("Failed to connect to MediaService: %v", err)
	}
	defer mdConn.Close()

	threadClient, tcConn, err := grpc.NewThreadServiceClient()
	if err != nil {
		log.Fatalf("Failed to connect to ThreadService: %v", err)
	}
	defer tcConn.Close()

	communityClient, tcConn, err := grpc.NewCommunityServiceClient()
	if err != nil {
		log.Fatalf("Failed to connect to Community Service: %v", err)
	}
	defer tcConn.Close()

	chatClient, tcConn, err := grpc.NewChatClients()
	if err != nil {
		log.Fatalf("Failed to connect to Community Service: %v", err)
	}
	defer tcConn.Close()

	routes.InitRoutes(
		router,
		userClient,
		authClients.Auth,
		authClients.SecurityAnswer,
		mediaClient,
		threadClient,
		communityClient,
		chatClient,
	)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API Gateway running on port %s", port)
	router.Run(":" + port)
}