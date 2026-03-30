package main

import (
	"assignment-service/internal/api/http"
	"assignment-service/internal/api/pb"
	"assignment-service/internal/clients"
	"assignment-service/internal/controller"
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	"assignment-service/internal/service"
	"assignment-service/internal/utils"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	http2 "net/http"
	"os"
	"strings"

	prismLog "github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	loadEnv()

	logger := initLogger()

	utils.ValidateEnvVars(logger,
		"APP_ENV",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT",
		"BUCKET_COUNT",
		"SALT_VALUE",
		"REDIS_HOST",
		"REDIS_PORT",
		"ASSIGNMENT_SERVICE_KAFKA_CONSUMER_GROUP_ID",
		"KAFKA_BOOTSTRAP_SERVER_HOST",
		"KAFKA_BOOTSTRAP_SERVER_PORT",
		"KAFKA_CACHE_INVALIDATIONS_TOPIC",
		"ASSIGNMENT_SERVICE_GRPC_SERVER_ADDRESS",
		"ASSIGNMENT_SERVICE_GRPC_SERVER_PORT",
	)

	grpcClient := getGrpcExperimentClient()
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
	experimentCache := service.NewExperimentConfigCache(redisClient, logger)
	assignmentService := service.NewAssignmentService(logger, bucketService, grpcClient, experimentCache)
	kafkaConsumer := service.NewKafkaConsumerImp(kafkaClient, logger)
	assignmentCacheInvalidationService := service.NewCacheInvalidationServiceKafka(kafkaConsumer, logger, experimentCache)

	// Controllers
	assignmentController := controller.NewAssignmentController(assignmentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		AssignmentController: assignmentController,
	})

	go assignmentCacheInvalidationService.ListenForInvalidations(context.Background())
	go startGrpcServer(logger, assignmentService)

	logger.Info("assignment-service started")
	http2.ListenAndServe(":8082", router)

}

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}

func initLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "production" {
		log.Fatal("APP_ENV must be set to development or production")
	}

	prismLog.InitLogger(env, "assignment-service")
	return prismLog.GetLogger()
}

func getGrpcExperimentClient() clients.ExperimentClient {
	address := fmt.Sprintf("%s:%s", os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS"), os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT"))
	client, err := clients.NewGrpcClient(address)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func startGrpcServer(logger *slog.Logger, assignmentService *service.AssignmentService) {
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(256*1024*1024),
		grpc.MaxSendMsgSize(256*1024*1024),
	)

	assignmentServer := pb.NewAssignmentServer(assignmentService, logger)
	assignment_service.RegisterAssignmentServiceServer(grpcServer, assignmentServer)

	reflection.Register(grpcServer)

	grpcAddress := os.Getenv("ASSIGNMENT_SERVICE_GRPC_SERVER_ADDRESS")
	grpcPort := os.Getenv("ASSIGNMENT_SERVICE_GRPC_SERVER_PORT")

	lis, err := net.Listen("tcp", strings.Join([]string{grpcAddress, grpcPort}, ":"))
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	logger.Info("gRPC server started on " + grpcAddress)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve grpc", "error", err)
		os.Exit(1)
	}
}
