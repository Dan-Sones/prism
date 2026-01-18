package main

import (
	"assignment-service/internal/api/http"
	"log"
	"os"
	"strconv"

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

	//salt, bucketCount := loadBucketConfig()

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{})

	//bucketService := service.NewBucketService(salt, bucketCount)

}
