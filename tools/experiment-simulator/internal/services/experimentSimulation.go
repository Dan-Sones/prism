package services

import (
	"encoding/binary"
	"errors"
	"experiment-simulator/internal/model"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"sync"
	"time"
)

type ExperimentSimulation struct {
	ExperimentConfig model.ExperimentConfig
	Performer        model.ActionPerformer
	UserIdService    UserIdService
}

func NewExperimentSimulation(experimentConfig model.ExperimentConfig, performer model.ActionPerformer, userIdService UserIdService) ExperimentSimulation {
	return ExperimentSimulation{
		UserIdService:    userIdService,
		ExperimentConfig: experimentConfig,
		Performer:        performer,
	}
}

func (es *ExperimentSimulation) SimulateExperiment() {
	aaParticipantsWithActions := es.GetAATestParticipantsWithActions()

	es.PerformAATest(*aaParticipantsWithActions)
}

func (es *ExperimentSimulation) PerformAATest(aaParticipantsWithActions []model.ExperimentParticipant) {

	durationSeconds := es.ExperimentConfig.AA.DurationSeconds
	totalActions := getTotalActions(aaParticipantsWithActions)

	fmt.Println("This Experiment Will Produce a total of", totalActions, "actions across all variants and the experiment will run for", durationSeconds, "seconds")
	fmt.Println("")

	fmt.Println("Variation Split")

	fmt.Println("")
	fmt.Println("Are you ready to begin the experiment simulation? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Experiment simulation aborted.")
		return
	}

	fmt.Println("")
	fmt.Println("----A/A Phase In Progress----")
	fmt.Println("")

	if es.PerformActions(aaParticipantsWithActions, durationSeconds, totalActions) {
		return
	}

	fmt.Println("----Experiment simulation completed.----")

}

func (es *ExperimentSimulation) PerformActions(particpantsWithActions []model.ExperimentParticipant, durationSeconds int, totalActions int) bool {
	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(len(particpantsWithActions)))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After(time.Duration(durationSeconds) * time.Second)
	totalActionsPerformed := 0
	currentParticipant := 0
	var mu sync.Mutex

	for {
		select {
		case <-timeout:
			return true
		case <-ticker.C:
			if totalActionsPerformed >= totalActions {
				return true
			}

			participant := particpantsWithActions[currentParticipant]
			go participant.PerformActionsWithDelay(es.Performer)
			mu.Lock()
			totalActionsPerformed += len(participant.Actions)
			fmt.Printf("Total Actions Performed: %d/%d\r", totalActionsPerformed, totalActions)
			mu.Unlock()
			currentParticipant++
		}
	}

	return false
}

func getTotalActions(particpantsWithActions []model.ExperimentParticipant) int {
	totalActions := 0
	for _, participant := range particpantsWithActions {
		totalActions += len(participant.Actions)
	}
	return totalActions
}

func (es *ExperimentSimulation) GetAATestParticipantsWithActions() *[]model.ExperimentParticipant {
	fmt.Println("Reading Experiment Config")

	controlVariantKey, err := es.GetControlVariantKey()
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	numExposureEventsForControlVariant, err := es.GetNumberOfEventTypeToPublishForVariantAndPhase("experiment_exposure", controlVariantKey, model.ExperimentPhaseAA)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	variantToEventToAmountExcludingExposure := es.GetVariantToEventToAmountExcludingExposure(model.ExperimentPhaseAA)

	fmt.Println("Fetching User IDs From Grpc")

	userIds, err := es.UserIdService.GetXUserIdsWithinExperimentAndVariant(*numExposureEventsForControlVariant, es.ExperimentConfig.ExperimentKey, controlVariantKey)
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

			fieldValues := make(map[model.EventField]interface{})

			for fieldName, fieldConfig := range eventFields {
				fieldSeed := DeriveSeed(es.ExperimentConfig.RandomSeed, string(controlVariantKey), string(eventKey), string(fieldName))
				randSouce := rand.New(rand.NewSource(fieldSeed))
				value := GenerateDataForField(fieldConfig.Type, fieldConfig.AA[controlVariantKey], randSouce)
				fieldValues[fieldName] = value
			}
			experimentParticipantsWithExposureEvents[i].AddAction(model.NewParticipantEventParameters(eventKey, fieldValues))
		}
	}

	return &experimentParticipantsWithExposureEvents
}

func GenerateDataForField(fieldType model.FieldType, fieldConfig model.FieldConfigMinMax, randSource rand.Source) any {

	rng := rand.New(randSource)

	if fieldType == model.FieldTypeFloat {
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

func (es *ExperimentSimulation) createExperimentParticipantsWithExposureEventsForUserIds(userIds []string, variantKey model.VariantKey, experimentKey string) []model.ExperimentParticipant {
	var experimentParticipants []model.ExperimentParticipant

	for _, userId := range userIds {
		participant := model.NewExperimentParticipant(userId, variantKey, experimentKey)
		participant.AddAction(model.NewParticipantEventParameters("experiment_exposure", map[model.EventField]interface{}{}))
		experimentParticipants = append(experimentParticipants, participant)
	}

	return experimentParticipants
}

func (es *ExperimentSimulation) GetVariantToEventToAmountExcludingExposure(phase model.ExperimentPhaseType) map[model.VariantKey]map[model.EventKey]int {
	variantToEventKeyToAmount := make(map[model.VariantKey]map[model.EventKey]int)

	for variantKey, varType := range es.GetVariantKeys() {
		if phase == model.ExperimentPhaseAA && varType == model.Treatment {
			continue
		}

		eventToAmountForVariant, err := es.GetNumberOfEventsToPublishForVariantAndPhase(variantKey, model.ExperimentPhaseAA)
		if err != nil {
			log.Fatalf("Failed to complete simulation: %v \n", err)
		}

		variantToEventKeyToAmount[variantKey] = eventToAmountForVariant
	}

	return variantToEventKeyToAmount
}

func (es *ExperimentSimulation) GetNumberOfEventsToPublishForVariantAndPhase(variantKey model.VariantKey, phase model.ExperimentPhaseType) (map[model.EventKey]int, error) {
	eventToAmount := make(map[model.EventKey]int)

	if phase == model.ExperimentPhaseAA {
		for eventKey, variantToAmountMap := range es.ExperimentConfig.AA.PublishAmounts {
			if eventKey == "experiment_exposure" {
				// Experiment Exposure events dealt with elsewhere
				continue
			}

			amount, exists := variantToAmountMap[variantKey]
			if !exists {
				return nil, errors.New(fmt.Sprintf("number of %s events to publish for variant_key %s in phase %s not in yml", eventKey, variantKey, model.ExperimentPhaseAA))
			}

			eventToAmount[eventKey] = amount

			return eventToAmount, nil
		}
	}

	// TODO: implement equiv for AB

	return nil, errors.New("it's not possible to reach here! we we are using an enum with only 2 options")

}

func (es *ExperimentSimulation) GetVariantKeys() map[model.VariantKey]model.VariantRole {
	return es.ExperimentConfig.Variants
}

func (es *ExperimentSimulation) GetControlVariantKey() (model.VariantKey, error) {
	for key, role := range es.ExperimentConfig.Variants {
		if role == model.Control {
			return key, nil
		}
	}

	return "", errors.New("no Control Variant Specified")
}

func (es *ExperimentSimulation) GetNumberOfEventTypeToPublishForVariantAndPhase(event_key model.EventKey, variant_key model.VariantKey, phase model.ExperimentPhaseType) (*int, error) {
	_, variantExists := es.ExperimentConfig.Variants[variant_key]
	if !variantExists {
		return nil, errors.New("variant Key Not Found in variants section")
	}

	if phase == model.ExperimentPhaseAA {
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

// Initially ALL users will be assigned to the control Variant.
// Then generate user id, then make call to assignment service via grpc to check if user is in desired experiment.
// If user is in experiment, add to list of user ids that can be used for aa test
// Once number of user ids == num exposure events
// Perform AA test

// wait buffer flush cooldown

// Validate correct number of events found. Validate variance is correct.
// THEN make http call to exp-service to confirm start of the experiment
// EITHER add an override start time so the experiment can begin immediately OR overwrite db record so it start immediatley rather than at midnight
// Then delete existing user ids that we had for a/a test
// the exp-service will have assigned new buckets and variant bounds.
// Then perform the same logic to get user ids. -> gen user id, fetch assignment from assignment service, if user is in either variant, record so.
// when we have users for each variant = exposure events for each variant
// gen those random events
// perform
// validate exp results are the same bruv

// Other thoughts:
// - To calculate variance well do we need to have more than one event per person rather than a 1->1 relationship like atm
