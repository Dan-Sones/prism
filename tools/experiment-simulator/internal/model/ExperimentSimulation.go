package model

import (
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

type ExperimentSimulation struct {
	ExperimentConfig ExperimentConfig
	VariantUserIds   VariantUserIds
	Performer        ActionPerformer
}

func NewExperimentSimulation(experimentConfig ExperimentConfig, variantUserIds VariantUserIds, performer ActionPerformer) ExperimentSimulation {
	return ExperimentSimulation{
		ExperimentConfig: experimentConfig,
		VariantUserIds:   variantUserIds,
		Performer:        performer,
	}
}

func (es *ExperimentSimulation) BeginExperiment() {
	var experimentParticipantsByVariant = make(map[string][]ExperimentParticipant)
	for _, variantKey := range es.ExperimentConfig.VariantKeys {
		experimentParticipantsByVariant[variantKey] = es.GetParticipantsForVariant(variantKey)
	}

	// combine all participants into a single slice and then shuffle the slice to simulate events happening in a random order across variants
	var allParticipants []ExperimentParticipant
	for _, participantsForVariant := range experimentParticipantsByVariant {
		allParticipants = append(allParticipants, participantsForVariant...)
	}
	rand.Shuffle(len(allParticipants), func(i, j int) {
		allParticipants[i], allParticipants[j] = allParticipants[j], allParticipants[i]
	})

	durationSeconds := es.ExperimentConfig.DurationSeconds
	totalActions := getTotalActions(allParticipants)

	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(totalActions))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After(time.Duration(durationSeconds) * time.Second)
	currentActionIndex := 0

	fmt.Println("This Experiment Will Produce a total of ", totalActions, "actions across all variants and the experiment will run for", durationSeconds, "seconds")
	fmt.Println("")

	fmt.Println("Variation Split")

	for _, variantKey := range es.ExperimentConfig.VariantKeys {
		numExposureEvents := getTotalEventsForVariantAndEventType(es.ExperimentConfig, variantKey, "experiment_exposure")
		for eventTypeKey := range es.ExperimentConfig.Events {
			if eventTypeKey == "experiment_exposure" {
				continue
			}
			numEvents := getTotalEventsForVariantAndEventType(es.ExperimentConfig, variantKey, eventTypeKey)
			fmt.Println("Variant", variantKey, ":", numExposureEvents, "exposures and", numEvents, eventTypeKey, "events")
		}
	}

	fmt.Println("")
	fmt.Println("Are you ready to begin the experiment simulation? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Experiment simulation aborted.")
		return
	}

	fmt.Println("")
	fmt.Println("----Experiment Simulation In Progress----")
	fmt.Println("")

	totalActionsPerformed := 0

	var mu sync.Mutex

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			if currentActionIndex >= len(allParticipants) {
				return
			}
			participant := allParticipants[currentActionIndex]
			go participant.PerformActionsWithDelay()
			mu.Lock()
			totalActionsPerformed++
			fmt.Printf("Total Actions Performed: %d/%d\r", totalActionsPerformed, totalActions)
			mu.Unlock()

			currentActionIndex++
		}
	}

	fmt.Println("----Experiment simulation completed.----")

}

func (es *ExperimentSimulation) GetParticipantsForVariant(variantKey string) []ExperimentParticipant {
	numExposureEventsToPublishForVariant := es.GetNumberOfExposureEventsToPublishForVariant(variantKey)
	participantUserIds := es.VariantUserIds.GetFirstXUserIdsForVariant(variantKey, numExposureEventsToPublishForVariant)

	usersForVariant := make([]ExperimentParticipant, len(participantUserIds))
	for i, userId := range participantUserIds {
		usersForVariant[i] = NewExperimentParticipant(userId, variantKey, es.ExperimentConfig.FeatureFlagKey, es.Performer)
	}

	for i := range numExposureEventsToPublishForVariant {
		// We need to add the exposure event as the first action for each user in the variant
		usersForVariant[i].AddAction(NewParticipantEventParameters("experiment_exposure", map[string]any{
			"feature_flag_key": es.ExperimentConfig.FeatureFlagKey,
		}))
	}

	// for each eventType (that is not experiment_exposure) get the count to publish for that variant
	for eventTypeKey, config := range es.ExperimentConfig.Events {
		if eventTypeKey == "experiment_exposure" {
			// We don't want to publish exposures seperate from interactions
			continue
		}

		// Get the count to publish for that variant
		countToPublishForVariant := config.CountToPublishForVariant[variantKey]

		// Create the properties for the event based on the config and add the event as an action for the first x users in the variant where x is the count to publish for that variant
		for i := range countToPublishForVariant {
			properties := make(map[string]any)
			for fieldKey, config := range config.Fields {

				if fieldKey == "feature_flag_key" {
					properties[fieldKey] = es.ExperimentConfig.FeatureFlagKey
					continue
				}

				if fieldKey == "variant_key" {
					properties[fieldKey] = variantKey
					continue
				}

				properties[fieldKey] = generateDataForField(config)
			}
			// Allows us to cycle through the users that are assigned
			usersForVariant[i%len(usersForVariant)].AddAction(NewParticipantEventParameters(eventTypeKey, properties))
		}

	}

	return usersForVariant
}

func getTotalActions(participants []ExperimentParticipant) int {
	totalActions := 0
	for _, participant := range participants {
		totalActions += len(participant.Actions)
	}
	return totalActions
}

func (es *ExperimentSimulation) GetNumberOfExposureEventsToPublishForVariant(variantKey string) int {
	exposureEventConfig, ok := es.ExperimentConfig.Events["experiment_exposure"]
	if !ok {
		log.Fatal("experiment_exposure event type config is required but not found in experiment config!!!")
	}

	countToPublish, ok := exposureEventConfig.CountToPublishForVariant[variantKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Count to publish for variant %s is required but not found in experiment_exposure event config!!!", variantKey))
	}

	return countToPublish
}

func generateDataForField(config FieldConfig) any {
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

func getTotalEventsForVariantAndEventType(experimentConfig ExperimentConfig, variantKey string, eventTypeKey string) int {
	eventConfig, ok := experimentConfig.Events[eventTypeKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Event type %s config is required but not found in experiment config!!!", eventTypeKey))
	}

	countToPublish, ok := eventConfig.CountToPublishForVariant[variantKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Count to publish for variant %s is required but not found in event type %s config!!!", variantKey, eventTypeKey))
	}

	return countToPublish
}
