package main

import (
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/services"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}

func main() {
	loadEnv()
	fmt.Print(
		"  _____                      _                      _     ____  _                 _       _             \n" +
			" | ____|_  ___ __   ___ _ __(_)_ __ ___   ___ _ __ | |_  / ___|(_)_ __ ___  _   _| | __ _| |_ ___  _ __ \n" +
			" |  _| \\ \\/ / '_ \\ / _ \\ '__| | '_ " + "`" + " _ \\ / _ \\ '_ \\| __| \\___ \\| | '_ " + "`" + " _ \\| | | | |/ _" + "`" + " | __/ _ \\| '__|\n" +
			" | |___ >  <| |_) |  __/ |  | | | | | | |  __/ | | | |_   ___) | | | | | | | |_| | | (_| | || (_) | |   \n" +
			" |_____/_/\\_\\ .__/ \\___|_|  |_|_| |_| |_|\\___|_| |_|\\__| |____/|_|_| |_| |_|\\__,_|_|\\__,_|\\__\\___/|_|   \n" +
			"            |_|\n",
	)

	port := os.Getenv("EVENTS_SERVICE_SERVER_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("Error converting port to int: %v\n", err)
		return
	}

	performer := services.NewActionPerformerHttp(os.Getenv("EVENTS_SERVICE_SERVER_HOST"), portInt)

	simDetails := services.GetSimulation()
	// TODO: maybe add support for conucrrent, but for now just get the first one.
	for _, experimentConfig := range simDetails {
		vuids := make(model.VariantUserIds)

		for _, variantKey := range experimentConfig.VariantKeys {
			vuids[variantKey] = services.GetUserIdsForVariant(variantKey)
		}

		simulation := model.NewExperimentSimulation(experimentConfig, vuids, performer)
		simulation.BeginExperiment()
		return
	}

}

//func beginSimulation(experimentName string, experimentConfig model.ExperimentConfig) {
//	var wg sync.WaitGroup
//
//	variantUserIds := make(model.VariantUserIds)
//	for _, variantKey := range experimentConfig.VariantKeys {
//		variantUserIds[variantKey] = services.GetUserIdsForVariant(variantKey)
//	}
//
//	simulation := model.NewExperimentSimulation(experimentConfig, variantUserIds)
//
//	for _, variantKey := range experimentConfig.VariantKeys {
//		wg.Add(1)
//		go simulation.PerformForVariant(variantKey, &wg)
//	}
//
//	wg.Wait()
//
//}
