package main

import (
	"context"
	"data-cooking-service/internal/clients"
	"data-cooking-service/internal/repository"
	"data-cooking-service/internal/services"
	"data-cooking-service/internal/utils"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Dan-Sones/prismhash"
	prismLog "github.com/Dan-Sones/prismlogger"
	microbatcher "github.com/Dan-Sones/prismmicrobatcher"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger := initLogger()

	utils.ValidateEnvVars(logger,
		"APP_ENV",
		"KAFKA_BOOTSTRAP_SERVER_HOST",
		"KAFKA_BOOTSTRAP_SERVER_PORT",
		"DATA_COOKING_SERVICE_KAFKA_CONSUMER_GROUP_ID",
		"KAFKA_EVENTS_TOPIC",
		"CLICKHOUSE_HOST",
		"CLICKHOUSE_USER",
		"CLICKHOUSE_PASSWORD",
		"CLICKHOUSE_DB",
		"CLICKHOUSE_NATIVE_PORT",
		"DATA_COOKING_SERVICE_MICROBATCH_SIZE",
		"DATA_COOKING_SERVICE_MICROBATCH_FLUSH_TIMEOUT_SECONDS",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS",
		"EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT",
		"SALT_VALUE",
		"BUCKET_COUNT",
	)

	experimentationAssignmentGrpcClient := getGrpcExperimentationAssignmentClient()
	defer experimentationAssignmentGrpcClient.Close()

	experimentationExperimentGrpcClient := getGrpcExperimentationExperimentClient()
	defer experimentationExperimentGrpcClient.Close()

	clickhouseConn, err := clients.NewClickhouseConnection()
	if err != nil {
		logger.Error("Failed to connect to Clickhouse", "error", err)
		os.Exit(1)
	}
	defer clickhouseConn.Close()

	kafkaClient, err := clients.GetKafkaClient()
	if err != nil {
		logger.Error("Failed to create Kafka client: ", "error", err)
		os.Exit(1)
	}

	microBatchSizeInt, err := strconv.Atoi(os.Getenv("DATA_COOKING_SERVICE_MICROBATCH_SIZE"))
	if err != nil {
		logger.Error("Invalid DATA_COOKING_SERVICE_MICROBATCH_SIZE, must be an integer", "error", err)
		os.Exit(1)
	}

	// repository
	cookedEventsRepository := repository.NewCookedEventsRepositoryClickhouse(clickhouseConn)

	// service
	salt, bucketCount := prismhash.GetBucketConfig()
	bucketService := prismhash.NewBucketService(salt, bucketCount)
	microBatchProcessor := services.NewMicroBatchProcessorImp(cookedEventsRepository, experimentationExperimentGrpcClient, experimentationAssignmentGrpcClient, bucketService, logger)
	eventReader := microbatcher.NewKafkaEventReader(kafkaClient, logger)
	microBatchService := microbatcher.NewMicroBatchingService(microBatchSizeInt, utils.GetFlushTimeoutDuration(), eventReader, microBatchProcessor, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Data Cooking Service Started")
	microBatchService.Start(ctx)
}

func initLogger() *slog.Logger {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "production" {
		log.Fatal("APP_ENV must be set to development or production")
	}

	prismLog.InitLogger(env, "clickhouse-writer")
	return prismLog.GetLogger()
}

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}

func getGrpcExperimentationAssignmentClient() clients.ExperimentationAssignmentClient {
	address := fmt.Sprintf("%s:%s", os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS"), os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT"))
	client, err := clients.NewGrpcExperimentationAssignmentClient(address)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getGrpcExperimentationExperimentClient() clients.ExperimentationExperimentClient {
	address := fmt.Sprintf("%s:%s", os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_ADDRESS"), os.Getenv("EXPERIMENTATION_SERVICE_GRPC_SERVER_PORT"))
	client, err := clients.NewGrpcExperimentationExperimentClient(address)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
