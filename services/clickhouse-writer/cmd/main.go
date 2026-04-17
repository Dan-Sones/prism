package main

import (
	"clickhouse-writer/internal/clients"
	"clickhouse-writer/internal/repository"
	"clickhouse-writer/internal/services"
	"clickhouse-writer/internal/utils"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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
		"CLICKHOUSE_WRITER_KAFKA_CONSUMER_GROUP_ID",
		"KAFKA_EVENTS_TOPIC",
		"CLICKHOUSE_HOST",
		"CLICKHOUSE_USER",
		"CLICKHOUSE_PASSWORD",
		"CLICKHOUSE_DB",
		"CLICKHOUSE_NATIVE_PORT",
		"CLICKHOUSE_WRITER_MICROBATCH_SIZE",
		"CLICK_HOUSE_WRITER_MICROBATCH_FLUSH_TIMEOUT_SECONDS",
	)

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

	microBatchSizeInt, err := strconv.Atoi(os.Getenv("CLICKHOUSE_WRITER_MICROBATCH_SIZE"))
	if err != nil {
		logger.Error("Invalid CLICKHOUSE_WRITER_MICROBATCH_SIZE, must be an integer", "error", err)
		os.Exit(1)
	}

	// repositories
	eventsRepository := repository.NewEventsRepositoryClickhouse(clickhouseConn)

	// services
	microBatchProcessor := services.NewMicroBatchProcessorImp(eventsRepository)
	eventReader := microbatcher.NewEventReaderImp(kafkaClient, logger)
	microBatchService := microbatcher.NewMicroBatchingService(microBatchSizeInt, utils.GetFlushTimeoutDuration(), eventReader, microBatchProcessor, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
