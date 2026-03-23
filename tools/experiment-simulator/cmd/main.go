package main

import (
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/parsers"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	simDetails := getSimulation()
	// TODO: maybe add support for conucrrent, but for now just get the first one.
	for experimentName, experimentConfig := range simDetails {
		beginSimulation(experimentName, experimentConfig)
		return
	}

}

func beginSimulation(experimentName string, experimentConfig model.ExperimentConfig) {
	var wg sync.WaitGroup

	for _, variantKey := range experimentConfig.VariantKeys {
		wg.Add(1)
		go performForVariant(variantKey, experimentConfig.Events, experimentConfig.DurationSeconds, experimentConfig.FeatureFlagKey, &wg)
	}

	wg.Wait()

}

func getSimulation() model.SimulationConfig {
	data, err := os.ReadFile(filepath.Join("resources", "experiment.yml"))
	if err != nil {
		log.Fatal(err)
	}
	return parsers.ParseExperimentConfig(data)
}

func performForVariant(variantName string, events map[string]model.EventConfig, durationSeconds int, flagKey string, parentWg *sync.WaitGroup) {

	var variantWg sync.WaitGroup

	for eventTypeKey, config := range events {
		count := config.CountToPublishForVariant[variantName]
		variantWg.Add(1)
		go performEventTypeForVariant(eventTypeKey, count, config.Fields, durationSeconds, &variantWg)
	}

	variantWg.Wait()
	parentWg.Done()
}

func performEventTypeForVariant(eventTypeKey string, countToPublish int, fields []map[string]model.FieldConfig, durationSeconds int, variantWg *sync.WaitGroup) {
	defer variantWg.Done()

	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(countToPublish))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After(time.Duration(durationSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			return
		case t := <-ticker.C:
			fmt.Printf("publishing", eventTypeKey, "at", t, "\n")
		}
	}
}
