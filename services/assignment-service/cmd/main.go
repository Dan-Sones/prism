package main

import (
	"assignment-service/internal/api/http"
	"assignment-service/internal/clients"
	"assignment-service/internal/controller"
	"assignment-service/internal/service"
	"assignment-service/internal/utils"
	"fmt"
	"log"
	"log/slog"
	http2 "net/http"
	"os"

	prismLog "github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger := initLogger()

	utils.ValidateEnvVars(logger,
		"APP_ENV",
		"ADMIN_SERVICE_GRPC_SERVER_ADDRESS",
		"ADMIN_SERVICE_GRPC_SERVER_PORT",
		"BUCKET_COUNT",
		"SALT_VALUE",
		"REDIS_HOST",
		"REDIS_PORT",
		"ASSIGNMENT_SERVICE_KAFKA_CONSUMER_GROUP_ID",
		"KAFKA_BOOTSTRAP_SERVER_HOST",
		"KAFKA_BOOTSTRAP_SERVER_PORT",
		"KAFKA_CACHE_INVALIDATIONS_TOPIC",
	)

	grpcClient := getGrpcAssignmentClient()
	defer grpcClient.Close()

	redisClient := clients.NewRedisClient()
	defer redisClient.Close()

	kafkaClient, err := clients.GetKafkaClient()
	if err != nil {
		log.Fatal("Failed to create Kafka client: ", err)
	}
	defer kafkaClient.Close()

	salt, bucketCount := utils.GetBucketConfig()

	// Services
	bucketService := service.NewBucketService(salt, bucketCount)
	assignmentCacheService := service.NewAssignmentCache(redisClient, logger)
	assignmentService := service.NewAssignmentService(logger, bucketService, grpcClient, assignmentCacheService)
	kafkaConsumer := service.NewKafkaConsumerImp(kafkaClient, logger)
	assignmentCacheInvalidationService := service.NewCacheInvalidationServiceKafka(kafkaConsumer, logger, assignmentCacheService)

	// Controllers
	assignmentController := controller.NewAssignmentController(assignmentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		AssignmentController: assignmentController,
	})

	go assignmentCacheInvalidationService.ListenForInvalidations()

	logger.Info("assignment-service started")
	http2.ListenAndServe(":8082", router)

}

func loadEnv() {
	if err := godotenv.Load("../../infrastructure/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "production" {
		log.Fatal("APP_ENV must be set to development or production")
	}

	prismLog.InitLogger(env, "assignment-service")
	return prismLog.GetLogger()
}

func getGrpcAssignmentClient() clients.AssignmentClient {
	address := fmt.Sprintf("%s:%s", os.Getenv("ADMIN_SERVICE_GRPC_SERVER_ADDRESS"), os.Getenv("ADMIN_SERVICE_GRPC_SERVER_PORT"))
	client, err := clients.NewGrpcClient(address)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
