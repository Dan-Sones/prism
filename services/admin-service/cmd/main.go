package main

import (
	"admin-service/internal/api/http"
	"admin-service/internal/environment"
	"fmt"
	"log"
	http2 "net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "development" {
		environment.InitLogger("development", "admin-service")
	} else if env == "production" {
		environment.InitLogger("production", "admin-service")
	} else {
		log.Fatal("APP_ENV must be set to development or production")
	}

	logger := environment.GetLogger()
	logger.Info("admin-service started")

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{})

	http2.ListenAndServe(":8080", router)

}
