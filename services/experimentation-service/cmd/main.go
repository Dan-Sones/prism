package main

import (
	"experimentation-service/internal/api/http"
	"experimentation-service/internal/api/pb"
	"experimentation-service/internal/clients"
	"experimentation-service/internal/controller"
	"experimentation-service/internal/utils"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"experimentation-service/internal/repository"
	"experimentation-service/internal/service"
	"log"
	http2 "net/http"

	"github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbEventsCatalog "experimentation-service/internal/grpc/generated/events_catalog/v1"
	pbExperimentationAssignment "experimentation-service/internal/grpc/generated/experimentation_service_assignment/v1"
	pbExperimentationExperiments "experimentation-service/internal/grpc/generated/experimentation_service_experiments/v1"
)

func main() {
	loadEnv()
	logger := initLogger()
	utils.ValidateEnvVars(logger,
		"APP_ENV",
		"EXPERIMENTATION_SERVICE_HTTP_PORT",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT",
		"BUCKET_COUNT",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"CLICKHOUSE_HOST",
		"CLICKHOUSE_NATIVE_PORT",
		"CLICKHOUSE_DB",
		"CLICKHOUSE_USER",
		"CLICKHOUSE_PASSWORD",
		"STATS_ENGINE_GRPC_SERVER_ADDRESS",
		"STATS_ENGINE_GRPC_SERVER_PORT",
	)

	pgPool := clients.GetPostgresConnectionPool()
	defer pgPool.Close()

	clickhouseConn, err := clients.NewClickhouseConnection()
	if err != nil {
		logger.Error("Failed to connect to Clickhouse", "error", err)
		os.Exit(1)
	}
	defer clickhouseConn.Close()

	statsEngineAddress := fmt.Sprintf("%s:%s", os.Getenv("STATS_ENGINE_GRPC_SERVER_ADDRESS"), os.Getenv("STATS_ENGINE_GRPC_SERVER_PORT"))

	statsEngineClient, err := clients.NewStatsEngineClient(statsEngineAddress)
	if err != nil {
		logger.Error("Failed to connect to Stats Engine gRPC server", "error", err)
		os.Exit(1)
	}
	defer statsEngineClient.Close()

	kafkaProducerClient, err := clients.GetKafkaProducerClient()
	if err != nil {
		logger.Error("Failed to create Kafka producer client", "error", err)
		os.Exit(1)
	}

	// Global Values
	bucketCount, err := utils.GetBucketCount()
	if err != nil {
		logger.Error("Failed to get bucket count", "error", err)
		os.Exit(1)
	}

	// Repositories
	experimentRepository := repository.NewExperimentRepository(pgPool)
	eventsCatalogRepository := repository.NewEventsCatalogRepository(pgPool)
	eventsRepository := repository.NewClickHouseEventsRepository(clickhouseConn)
	metricsCatalogRepository := repository.NewMetricsCatalogRepository(pgPool)
	bucketAllocationRepository := repository.NewBucketAllocationRepository(pgPool)
	experimentPhaseRepository := repository.NewExperimentPhaseRepository(pgPool)

	// Services
	cacheInvalidationProducer := service.NewCacheInvalidationProducer(kafkaProducerClient)
	bucketAllocationService := service.NewBucketAllocationService(bucketAllocationRepository, logger)
	eventService := service.NewEventsService(eventsRepository, eventsCatalogRepository, logger)
	metricsCatalogService := service.NewMetricsCatalogService(metricsCatalogRepository, eventsCatalogRepository, logger)
	clickhouseQueryBuilder := service.NewClickhouseQueryBuilder()
	experimentService := service.NewExperimentService(experimentRepository, bucketAllocationService, clickhouseQueryBuilder, eventService, metricsCatalogService, statsEngineClient, experimentPhaseRepository, cacheInvalidationProducer, logger)
	experimentResultsRepository := repository.NewExperimentResultsRepository(pgPool)
	experimentResultsService := service.NewExperimentResultsService(experimentPhaseRepository, experimentResultsRepository, statsEngineClient, experimentService, metricsCatalogService, eventsRepository, clickhouseQueryBuilder, logger)
	assignmentService := service.NewAssignmentService(experimentRepository, bucketCount, logger)
	eventsCatalogService := service.NewEventsCatalogService(eventsCatalogRepository, logger)

	// Controllers
	experimentController := controller.NewExperimentController(experimentService)
	experimentResultsController := controller.NewExperimentResultsController(experimentResultsService)
	eventsCatalogController := controller.NewEventsCatalogController(eventsCatalogService)
	eventController := controller.NewEventController(eventService)
	metricsCatalogController := controller.NewMetricsCatalogController(metricsCatalogService)

	go startGrpcServer(logger, assignmentService, eventsCatalogService, experimentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		ExperimentController:        experimentController,
		ExperimentResultsController: experimentResultsController,
		EventsCatalogController:     eventsCatalogController,
		EventController:             eventController,
		MetricsCatalogController:    metricsCatalogController,
	})

	httpPort := fmt.Sprintf(":%s", os.Getenv("EXPERIMENTATION_SERVICE_HTTP_PORT"))

	logger.Info("experimentation-service started")

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

	prismLog.InitLogger(env, "experimentation-service")
	return prismLog.GetLogger()
}

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}

func startGrpcServer(logger *slog.Logger, assignmentService *service.AssignmentService, eventsCatalogService *service.EventsCatalogService, experimentService *service.ExperimentService) {
	grpcServer := grpc.NewServer()

	experimentationAssignmentServer := pb.NewAssignmentServer(assignmentService)
	eventsCatalogServer := pb.NewEventsCatalogServer(eventsCatalogService)
	experimentsServer := pb.NewExperimentsServer(experimentService)
	pbExperimentationAssignment.RegisterExperimentationServiceAssignmentServer(grpcServer, experimentationAssignmentServer)
	pbEventsCatalog.RegisterEventsCatalogServiceServer(grpcServer, eventsCatalogServer)
	pbExperimentationExperiments.RegisterExperimentationServiceExperimentsServer(grpcServer, experimentsServer)

	reflection.Register(grpcServer)

	grpcAddress := os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS")
	grpcPort := os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT")

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
