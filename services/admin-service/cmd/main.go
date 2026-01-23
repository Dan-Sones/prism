package main

import (
	"admin-service/internal/api/http"
	"admin-service/internal/api/pb"
	"admin-service/internal/clients"
	"admin-service/internal/controller"
	"admin-service/internal/utils"
	"log/slog"
	"net"

	"admin-service/internal/repository"
	"admin-service/internal/service"
	"log"
	http2 "net/http"
	"os"

	"github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb2 "admin-service/server/api/proto/assignment/v1"
)

func main() {
	loadEnv()
	logger := initLogger()
	logger.Info("admin-service started")

	pgPool := clients.GetPostgresConnectionPool()
	defer pgPool.Close()

	// Global Values
	bucketCount := mustGetBucketCount()

	// Repositories
	experimentRepository := repository.NewExperimentRepository(pgPool, logger)

	// Services
	experimentService := service.NewExperimentService(experimentRepository, logger)
	assignmentService := service.NewAssignmentService(experimentRepository, bucketCount, logger)

	// Controllers
	experimentController := controller.NewExperimentController(experimentService)

	go startGrpcServer(logger, assignmentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		ExperimentController: experimentController,
	})

	http2.ListenAndServe(":8080", router)

}

func mustGetBucketCount() int32 {
	bucketCount, err := utils.GetBucketCount()
	if err != nil {
		log.Fatalf("Failed to get bucket count: %v", err)
	}
	return bucketCount
}

func initLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	switch env {
	case "development", "production":
		prismLog.InitLogger(env, "admin-service")
	default:
		log.Fatal("APP_ENV must be set to development or production")
	}
	return prismLog.GetLogger()
}

func loadEnv() {
	if err := godotenv.Load("../../infrastructure/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func startGrpcServer(logger *slog.Logger, assignmentService *service.AssignmentService) {
	grpcServer := grpc.NewServer()

	assignmentServer := pb.NewAssignmentServer(assignmentService)
	pb2.RegisterAssignmentServiceServer(grpcServer, assignmentServer)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Info("gRPC server started on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
