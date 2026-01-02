package main

import (
	"admin-service/internal/api/http"
	"admin-service/internal/clients"
	"admin-service/internal/environment"
	"context"
	"fmt"
	"log"
	http2 "net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load("../../infrastructure/.env")
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

	pgPool := clients.GetPostgresConnectionPool()
	defer pgPool.Close()

	var greeting string
	err = pgPool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{})

	http2.ListenAndServe(":8080", router)

}
