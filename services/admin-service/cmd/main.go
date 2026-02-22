package main

import (
	"admin-service/internal/api/http"
	"admin-service/internal/api/pb"
	"admin-service/internal/clients"
	"admin-service/internal/controller"
	"admin-service/internal/utils"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"admin-service/internal/repository"
	"admin-service/internal/service"
	"log"
	http2 "net/http"

	"github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbAssignment "admin-service/internal/grpc/generated/assignment/v1"
	pbEventsCatalog "admin-service/internal/grpc/generated/events_catalog/v1"
)

func main() {
	loadEnv()
	logger := initLogger()
	utils.ValidateEnvVars(logger,
		"APP_ENV",
		"ADMIN_SERVICE_HTTP_PORT",
		"ADMIN_SERVICE_GRPC_SERVER_ADDRESS",
		"ADMIN_SERVICE_GRPC_SERVER_PORT",
		"BUCKET_COUNT",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	)
	logger.Info("admin-service started")

	pgPool := clients.GetPostgresConnectionPool()
	defer pgPool.Close()

	// Global Values
	bucketCount, err := utils.GetBucketCount()
	if err != nil {
		logger.Error("Failed to get bucket count", "error", err)
		os.Exit(1)
	}

	// Repositories
	experimentRepository := repository.NewExperimentRepository(pgPool)
	_ = repository.NewEventsCatalogRepository(pgPool)

	// Services
	experimentService := service.NewExperimentService(experimentRepository, logger)
	assignmentService := service.NewAssignmentService(experimentRepository, bucketCount, logger)
	eventsCatalogService := service.NewEventsCatalogService(repository.NewEventsCatalogRepository(pgPool), logger)

	// Controllers
	experimentController := controller.NewExperimentController(experimentService)
	eventsCatalogController := controller.NewEventsCatalogController(eventsCatalogService)

	go startGrpcServer(logger, assignmentService, eventsCatalogService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		ExperimentController:    experimentController,
		EventsCatalogController: eventsCatalogController,
	})

	httpPort := fmt.Sprintf(":%s", os.Getenv("ADMIN_SERVICE_HTTP_PORT"))

	err = http2.ListenAndServe(httpPort, router)
	if err != nil {
		logger.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}

func initLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "production" {
		log.Fatal("APP_ENV must be set to development or production")
	}

	prismLog.InitLogger(env, "admin-service")
	return prismLog.GetLogger()
}

func loadEnv() {
	if err := godotenv.Load("../../infrastructure/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func startGrpcServer(logger *slog.Logger, assignmentService *service.AssignmentService, eventsCatalogService *service.EventsCatalogService) {
	grpcServer := grpc.NewServer()

	assignmentServer := pb.NewAssignmentServer(assignmentService)
	eventsCatalogServer := pb.NewEventsCatalogServer(eventsCatalogService)
	pbAssignment.RegisterAssignmentServiceServer(grpcServer, assignmentServer)
	pbEventsCatalog.RegisterEventsCatalogServiceServer(grpcServer, eventsCatalogServer)

	reflection.Register(grpcServer)

	grpcAddress := os.Getenv("ADMIN_SERVICE_GRPC_SERVER_ADDRESS")
	grpcPort := os.Getenv("ADMIN_SERVICE_GRPC_SERVER_PORT")

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
