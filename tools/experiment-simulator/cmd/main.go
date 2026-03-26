package main

import (
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/services"
	"fmt"
)

func main() {
	fmt.Print(
		"  _____                      _                      _     ____  _                 _       _             \n" +
			" | ____|_  ___ __   ___ _ __(_)_ __ ___   ___ _ __ | |_  / ___|(_)_ __ ___  _   _| | __ _| |_ ___  _ __ \n" +
			" |  _| \\ \\/ / '_ \\ / _ \\ '__| | '_ " + "`" + " _ \\ / _ \\ '_ \\| __| \\___ \\| | '_ " + "`" + " _ \\| | | | |/ _" + "`" + " | __/ _ \\| '__|\n" +
			" | |___ >  <| |_) |  __/ |  | | | | | | |  __/ | | | |_   ___) | | | | | | | |_| | | (_| | || (_) | |   \n" +
			" |_____/_/\\_\\ .__/ \\___|_|  |_|_| |_| |_|\\___|_| |_|\\__| |____/|_|_| |_| |_|\\__,_|_|\\__,_|\\__\\___/|_|   \n" +
			"            |_|\n",
	)

	performer := services.ActionPerformerPrint{}

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
