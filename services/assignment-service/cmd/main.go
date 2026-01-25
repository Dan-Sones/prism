package main

import (
	"assignment-service/internal/api/http"
	"assignment-service/internal/clients"
	"assignment-service/internal/controller"
	"assignment-service/internal/service"
	"assignment-service/internal/utils"
	"log"
	http2 "net/http"
	"os"

	prismLog "github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../infrastructure/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "development" {
		prismLog.InitLogger("development", "assignment-service")
	} else if env == "production" {
		prismLog.InitLogger("production", "assignment-service")
	} else {
		log.Fatal("APP_ENV must be set to development or production")
	}

	logger := prismLog.GetLogger()
	logger.Info("assignment-service started")

	pgPool := clients.GetPostgresConnectionPool()
	defer pgPool.Close()

	grpcClient := getGrpcClient()
	defer grpcClient.Close()

	salt, bucketCount := utils.GetBucketConfig()

	// Services
	bucketService := service.NewBucketService(salt, bucketCount)
	assignmentService := service.NewAssignmentService(logger, bucketService, *grpcClient)

	// Controllers
	assignmentController := controller.NewAssignmentController(assignmentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		AssignmentController: assignmentController,
	})

	http2.ListenAndServe(":8082", router)

}

func getGrpcClient() *clients.GrpcClient {
	client, err := clients.NewGrpcClient(mustGetGrpcAddr())
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func mustGetGrpcAddr() string {
	grpcAddr := os.Getenv("ADMIN_SERVICE_GRPC_ADDR")
	if grpcAddr == "" {
		log.Fatal("ADMIN_SERVICE_GRPC_ADDR is not set")
	}
	return grpcAddr
}
