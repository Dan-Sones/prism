package main

import (
	"encoding/json"
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/parsers"
	"fmt"
	"log"
	"math/rand/v2"
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

	//CONTINUE FROM HERE:
	// we need to work out how we decide what user ids to use
	// these then need to be passed down all the way to the publish function.

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

func performForVariant(variantKey string, events map[string]model.EventConfig, durationSeconds int, flagKey string, parentWg *sync.WaitGroup) {

	var variantWg sync.WaitGroup

	for eventTypeKey, config := range events {
		count := config.CountToPublishForVariant[variantKey]
		variantWg.Add(1)
		go performEventTypeForVariant(variantKey, eventTypeKey, count, config.Fields, durationSeconds, &variantWg)
	}

	variantWg.Wait()
	parentWg.Done()
}

func performEventTypeForVariant(variantKey, eventTypeKey string, countToPublish int, fields map[string]model.FieldConfig, durationSeconds int, variantWg *sync.WaitGroup) {
	defer variantWg.Done()

	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(countToPublish))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After(time.Duration(durationSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			publishEventTypeForVariant(variantKey, eventTypeKey, fields)
		}
	}
}

func publishEventTypeForVariant(variantKey, eventTypeKey string, fields map[string]model.FieldConfig) {
	properties := make(map[string]any)

	for fieldKey, config := range fields {

		if fieldKey == "feature_flag_key" {
			properties[fieldKey] = eventTypeKey
			continue
		}

		if fieldKey == "variant_key" {
			properties[fieldKey] = variantKey
			continue
		}

		properties[fieldKey] = generateDataForField(config)
	}

	eventRequest := model.EventRequest{
		EventKey:    eventTypeKey,
		UserDetails: model.UserDetails{},
		SentAt:      time.Now(),
		Properties:  properties,
	}

	eventRequestStr, err := json.Marshal(eventRequest)
	if err != nil {
		log.Fatal("Failed to marshal json... aborting!!!!")
	}

	fmt.Println(string(eventRequestStr))
}

func generateDataForField(config model.FieldConfig) any {

	if config.Min != nil && config.Max != nil {
		min := *config.Min
		max := *config.Max

		if config.Type == "int" {
			return rand.IntN(int(max-min)) + int(min)
		}

		if config.Type == "float" {
			return rand.Float64()*(max-min) + min
		}
	}

	if config.Type == "int" {
		min := 0
		max := 100
		return rand.IntN(int(max-min)) + int(min)
	}

	if config.Type == "float" {
		min := float64(0)
		max := float64(100)
		return rand.Float64()*(max-min) + min
	}

	// I've realised for strings this is going to be even more complicated. we need to allow the user to have a list of potential values and then weights of each value being chosen? Maybe I just make this
	//initial version support int and float values only?
	log.Fatal("This script only supports int and float event fields!!!")
	return nil
}
