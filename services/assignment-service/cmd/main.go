package main

import (
	"assignment-service/internal/api/http"
	"assignment-service/internal/clients"
	"assignment-service/internal/repository"
	"assignment-service/internal/service"
	"log"
	"os"
	"strconv"

	prismLog "github.com/Dan-Sones/prismlogger"
	"github.com/joho/godotenv"
)

func loadBucketConfig() (string, int) {

	salt := os.Getenv("SALT_VALUE")
	bucketCountStr := os.Getenv("BUCKET_COUNT")
	bucketCount, err := strconv.Atoi(bucketCountStr)
	if err != nil {
		log.Fatalf("Invalid BUCKET_COUNT: %v", err)
	}
	return salt, bucketCount
}

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

	//salt, bucketCount := loadBucketConfig()

	// Repositories
	experimentRepository := repository.NewExperimentRepository(pgPool, logger)

	// Services
	//bucketService := service.NewBucketService(salt, bucketCount)
	experimentService := service.NewExperimentService(experimentRepository, logger)

	// Controllers

	experimentService.GetVariantsForBucket(1)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{})

}
