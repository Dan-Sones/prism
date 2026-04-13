package model

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
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

func (es *ExperimentSimulation) GetAATestParticipantsWithActions() *[]ExperimentParticipant {
	controlVariantKey, err := es.GetControlVariantKey()
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	numExposureEventsForControlVariant, err := es.GetNumberOfEventTypeToPublishForVariantAndPhase("experiment_exposure", controlVariantKey, ExperimentPhaseAA)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	variantToEventToAmountExcludingExposure := es.GetVariantToEventToAmountExcludingExposure(ExperimentPhaseAA)

	userIds, err := es.VariantUserIds.getUserIdsForVariant(controlVariantKey, *numExposureEventsForControlVariant)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	experimentParticipantsWithExposureEvents := es.createExperimentParticipantsWithExposureEventsForUserIds(userIds, controlVariantKey, es.ExperimentConfig.ExperimentKey)

	numExperimentParticipants := len(experimentParticipantsWithExposureEvents)

	for eventKey, numOfEventToPublish := range variantToEventToAmountExcludingExposure[controlVariantKey] {
		for i := 0; i < numOfEventToPublish; i++ {
			if i >= numExperimentParticipants {
				log.Fatalf("You defined %d experiment_exposure events but then asked for %d %s events - That's not possible! ", numExperimentParticipants, numOfEventToPublish, eventKey)
			}

			eventFields := es.ExperimentConfig.Events[eventKey].Fields

			fieldValues := make(map[EventField]interface{})

			for fieldName, fieldConfig := range eventFields {
				fieldSeed := DeriveSeed(es.ExperimentConfig.RandomSeed, string(controlVariantKey), string(eventKey), string(fieldName))
				randSouce := rand.New(rand.NewSource(fieldSeed))
				value := GenerateDataForField(fieldConfig.Type, fieldConfig.AA[controlVariantKey], randSouce)
				fieldValues[fieldName] = value
			}
			experimentParticipantsWithExposureEvents[i].AddAction(NewParticipantEventParameters(eventKey, fieldValues))
		}
	}

	return &experimentParticipantsWithExposureEvents
}

func GenerateDataForField(fieldType FieldType, fieldConfig FieldConfigMinMax, randSource rand.Source) any {

	rng := rand.New(randSource)

	if fieldType == FieldTypeFloat {
		mx := *fieldConfig.Max
		mn := *fieldConfig.Min
		value := rng.Float64()*(mx-mn) + mn
		return value
	}

	// TODO: support datatypes other than float
	log.Fatalf("We only support a float field type, you need to change the simulation config - sorry ):")
	return nil
}

func DeriveSeed(parent int64, parts ...string) int64 {
	h := fnv.New64a()
	binary.Write(h, binary.LittleEndian, parent)
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return int64(h.Sum64())
}

func (es *ExperimentSimulation) createExperimentParticipantsWithExposureEventsForUserIds(userIds []string, variantKey VariantKey, experimentKey string) []ExperimentParticipant {
	var experimentParticipants []ExperimentParticipant

	for _, userId := range userIds {
		participant := NewExperimentParticipant(userId, variantKey, experimentKey)
		participant.AddAction(NewParticipantEventParameters("experiment_exposure", map[EventField]interface{}{}))
		experimentParticipants = append(experimentParticipants, participant)
	}

	return experimentParticipants
}

func (es *ExperimentSimulation) GetVariantToEventToAmountExcludingExposure(phase ExperimentPhaseType) map[VariantKey]map[EventKey]int {
	variantToEventKeyToAmount := make(map[VariantKey]map[EventKey]int)

	for variantKey, varType := range es.GetVariantKeys() {
		if phase == ExperimentPhaseAA && varType == Treatment {
			continue
		}

		eventToAmountForVariant, err := es.GetNumberOfEventsToPublishForVariantAndPhase(variantKey, ExperimentPhaseAA)
		if err != nil {
			log.Fatalf("Failed to complete simulation: %v \n", err)
		}

		variantToEventKeyToAmount[variantKey] = eventToAmountForVariant
	}

	return variantToEventKeyToAmount
}

func (es *ExperimentSimulation) GetNumberOfEventsToPublishForVariantAndPhase(variantKey VariantKey, phase ExperimentPhaseType) (map[EventKey]int, error) {
	eventToAmount := make(map[EventKey]int)

	if phase == ExperimentPhaseAA {
		for eventKey, variantToAmountMap := range es.ExperimentConfig.AA.PublishAmounts {
			if eventKey == "experiment_exposure" {
				// Experiment Exposure events dealt with elsewhere
				continue
			}

			amount, exists := variantToAmountMap[variantKey]
			if !exists {
				return nil, errors.New(fmt.Sprintf("number of %s events to publish for variant_key %s in phase %s not in yml", eventKey, variantKey, ExperimentPhaseAA))
			}

			eventToAmount[eventKey] = amount

			return eventToAmount, nil
		}
	}

	// TODO: implement equiv for AB

	return nil, errors.New("it's not possible to reach here! we we are using an enum with only 2 options")

}

func (es *ExperimentSimulation) GetVariantKeys() map[VariantKey]VariantRole {
	return es.ExperimentConfig.Variants
}

func (es *ExperimentSimulation) GetControlVariantKey() (VariantKey, error) {
	for key, role := range es.ExperimentConfig.Variants {
		if role == Control {
			return key, nil
		}
	}

	return "", errors.New("no Control Variant Specified")
}

func (es *ExperimentSimulation) GetNumberOfEventTypeToPublishForVariantAndPhase(event_key EventKey, variant_key VariantKey, phase ExperimentPhaseType) (*int, error) {
	_, variantExists := es.ExperimentConfig.Variants[variant_key]
	if !variantExists {
		return nil, errors.New("variant Key Not Found in variants section")
	}

	if phase == ExperimentPhaseAA {
		_, eventKeyExistsInSection := es.ExperimentConfig.AA.PublishAmounts[event_key]
		if !eventKeyExistsInSection {
			return nil, errors.New("event_key not found in aa 'publish_amounts' section")
		}

		amountToPublish, amountExistsInSection := es.ExperimentConfig.AA.PublishAmounts[event_key][variant_key]
		if !amountExistsInSection {
			return nil, errors.New("publish_amount not found for variant")
		}

		return &amountToPublish, nil
	}

	// TODO: AB Equivalent

	return nil, errors.New("it's not possible to reach here! we we are using an enum with only 2 options")

}

//func (es *ExperimentSimulation) BeginExperiment() {
//	var experimentParticipantsByVariant = make(map[string][]ExperimentParticipant)
//	for _, variantKey := range es.ExperimentConfig.ExperimentConfigAB.VariantKeys {
//		experimentParticipantsByVariant[variantKey] = es.GetParticipantsForVariant(variantKey)
//	}
//
//	// combine all participants into a single slice and then shuffle the slice to simulate events happening in a random order across variants
//	var allParticipants []ExperimentParticipant
//	for _, participantsForVariant := range experimentParticipantsByVariant {
//		allParticipants = append(allParticipants, participantsForVariant...)
//	}
//	rand.Shuffle(len(allParticipants), func(i, j int) {
//		allParticipants[i], allParticipants[j] = allParticipants[j], allParticipants[i]
//	})
//
//	durationSeconds := es.ExperimentConfig.ExperimentConfigAB.DurationSeconds
//	totalActions := getTotalActions(allParticipants)
//
//	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(totalActions))
//	ticker := time.NewTicker(interval)
//	defer ticker.Stop()
//	timeout := time.After(time.Duration(durationSeconds) * time.Second)
//	currentActionIndex := 0
//
//	fmt.Println("This Experiment Will Produce a total of ", totalActions, "actions across all variants and the experiment will run for", durationSeconds, "seconds")
//	fmt.Println("")
//
//	fmt.Println("Variation Split")
//
//	for _, variantKey := range es.ExperimentConfig.ExperimentConfigAB.VariantKeys {
//		numExposureEvents := es.GetTotalEventsForVariantAndEventType(variantKey, "experiment_exposure")
//		for eventTypeKey := range es.ExperimentConfig.ExperimentConfigAB.Events {
//			if eventTypeKey == "experiment_exposure" {
//				continue
//			}
//			numEvents := es.GetTotalEventsForVariantAndEventType(variantKey, eventTypeKey)
//			fmt.Println("Variant", variantKey, ":", numExposureEvents, "exposures and", numEvents, eventTypeKey, "events")
//		}
//	}
//
//	fmt.Println("")
//	fmt.Println("Are you ready to begin the experiment simulation? (y/n)")
//	var input string
//	fmt.Scanln(&input)
//	if input != "y" {
//		fmt.Println("Experiment simulation aborted.")
//		return
//	}
//
//	fmt.Println("")
//	fmt.Println("----Experiment Simulation In Progress----")
//	fmt.Println("")
//
//	totalActionsPerformed := 0
//	var mu sync.Mutex
//
//	for {
//		select {
//		case <-timeout:
//			return
//		case <-ticker.C:
//			if currentActionIndex >= len(allParticipants) {
//				return
//			}
//			participant := allParticipants[currentActionIndex]
//			go participant.PerformActionsWithDelay()
//			mu.Lock()
//			totalActionsPerformed += len(participant.Actions)
//			fmt.Printf("Total Actions Performed: %d/%d\r", totalActionsPerformed, totalActions)
//			mu.Unlock()
//
//			currentActionIndex++
//		}
//	}
//
//	fmt.Println("----Experiment simulation completed.----")
//
//}
//
//func (es *ExperimentSimulation) GetParticipantsForVariant(variantKey string) []ExperimentParticipant {
//	numExposureEventsToPublishForVariant := es.GetNumberOfExposureEventsToPublishForVariant(variantKey)
//	participantUserIds := es.VariantUserIds.GetFirstXUserIdsForVariant(variantKey, numExposureEventsToPublishForVariant)
//
//	usersForVariant := make([]ExperimentParticipant, len(participantUserIds))
//	for i, userId := range participantUserIds {
//		usersForVariant[i] = NewExperimentParticipant(userId, variantKey, es.ExperimentConfig.ExperimentConfigAB.FeatureFlagKey, es.Performer)
//	}
//
//	for i := range numExposureEventsToPublishForVariant {
//		// We need to add the exposure event as the first action for each user in the variant
//		usersForVariant[i].AddAction(NewParticipantEventParameters("experiment_exposure", map[string]any{}))
//	}
//
//	// for each eventType (that is not experiment_exposure) get the count to publish for that variant
//	for eventTypeKey, config := range es.ExperimentConfig.ExperimentConfigAB.Events {
//		if eventTypeKey == "experiment_exposure" {
//			// We don't want to publish exposures seperate from interactions
//			continue
//		}
//
//		// Get the count to publish for that variant
//		countToPublishForVariant := config.CountToPublishForVariant[variantKey]
//
//		// Create the properties for the event based on the config and add the event as an action for the first x users in the variant where x is the count to publish for that variant
//		for i := range countToPublishForVariant {
//			properties := make(map[string]any)
//			for fieldKey, config := range config.Fields {
//				properties[fieldKey] = generateDataForField(config)
//			}
//			// Allows us to cycle through the users that are assigned
//			usersForVariant[i%len(usersForVariant)].AddAction(NewParticipantEventParameters(eventTypeKey, properties))
//		}
//
//	}
//
//	return usersForVariant
//}
//
//func (es *ExperimentSimulation) GetTotalEventsForVariantAndEventType(variantKey string, eventTypeKey string) int {
//	eventConfig, ok := es.ExperimentConfig.ExperimentConfigAB.Events[eventTypeKey]
//	if !ok {
//		log.Fatal(fmt.Sprintf("Event type %s config is required but not found in experiment config!!!", eventTypeKey))
//	}
//
//	countToPublish, ok := eventConfig.CountToPublishForVariant[variantKey]
//	if !ok {
//		log.Fatal(fmt.Sprintf("Count to publish for variant %s is required but not found in event type %s config!!!", variantKey, eventTypeKey))
//	}
//
//	return countToPublish
//}
//
//func getTotalActions(participants []ExperimentParticipant) int {
//	totalActions := 0
//	for _, participant := range participants {
//		totalActions += len(participant.Actions)
//	}
//	return totalActions
//}
//
//func (es *ExperimentSimulation) GetNumberOfExposureEventsToPublishForVariant(variantKey string) int {
//	exposureEventConfig, ok := es.ExperimentConfig.ExperimentConfigAB.Events["experiment_exposure"]
//	if !ok {
//		log.Fatal("experiment_exposure event type config is required but not found in experiment config!!!")
//	}
//
//	countToPublish, ok := exposureEventConfig.CountToPublishForVariant[variantKey]
//	if !ok {
//		log.Fatal(fmt.Sprintf("Count to publish for variant %s is required but not found in experiment_exposure event config!!!", variantKey))
//	}
//
//	return countToPublish
//}
//
//func generateDataForField(config FieldConfig) any {
//	if config.Min != nil && config.Max != nil {
//		min := *config.Min
//		max := *config.Max
//
//		if config.Type == "int" {
//			return rand.IntN(int(max-min)) + int(min)
//		}
//
//		if config.Type == "float" {
//			return rand.Float64()*(max-min) + min
//		}
//	}
//
//	if config.Type == "int" {
//		min := 0
//		max := 100
//		return rand.IntN(int(max-min)) + int(min)
//	}
//
//	if config.Type == "float" {
//		min := float64(0)
//		max := float64(100)
//		return rand.Float64()*(max-min) + min
//	}
//
//	// I've realised for strings this is going to be even more complicated. we need to allow the user to have a list of potential values and then weights of each value being chosen? Maybe I just make this
//	//initial version support int and float values only?
//	log.Fatal("This script only supports int and float event fields!!!")
//	return nil
//}
