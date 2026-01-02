package main

import (
	"admin-service/internal/api/http"
	"admin-service/internal/clients"
	"admin-service/internal/controller"
	"admin-service/internal/environment"
	"admin-service/internal/repository"
	"admin-service/internal/service"
	"log"
	http2 "net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
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

	// Repositories
	experimentRepository := repository.NewExperimentRepository(pgPool, logger)

	// Services
	experimentService := service.NewExperimentService(experimentRepository, logger)

	// Controllers
	experimentController := controller.NewExperimentController(experimentService)

	router := http.NewRouter()
	http.RegisterRoutes(router, http.Controllers{
		ExperimentController: experimentController,
	})

	http2.ListenAndServe(":8080", router)

}
