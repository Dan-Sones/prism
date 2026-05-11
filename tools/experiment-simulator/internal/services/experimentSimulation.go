package services

import (
	"context"
	"encoding/binary"

	"errors"
	"experiment-simulator/internal/assertors"
	"experiment-simulator/internal/model"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type ExperimentSimulation struct {
	ExperimentConfig          model.ExperimentConfig
	AssertionService          assertors.AssertionService
	Performer                 model.ActionPerformer
	UserIdService             UserIdService
	CacheInvalidationProducer *CacheInvalidationProducer
	ExperimentTimingService   *ExperimentTimingService
}

func NewExperimentSimulation(experimentConfig model.ExperimentConfig,
	performer model.ActionPerformer,
	userIdService UserIdService,
	assertionService assertors.AssertionService,
	cacheInvalidationProducer *CacheInvalidationProducer,
	experimentTimingService *ExperimentTimingService,
) ExperimentSimulation {
	return ExperimentSimulation{
		UserIdService:             userIdService,
		ExperimentConfig:          experimentConfig,
		Performer:                 performer,
		AssertionService:          assertionService,
		CacheInvalidationProducer: cacheInvalidationProducer,
		ExperimentTimingService:   experimentTimingService,
	}
}

func (es *ExperimentSimulation) SimulateExperiment() {
	// todo: improve formatting between stages - messy at the moment
	err := es.ExperimentTimingService.MoveAAStartToNow(es.ExperimentConfig.ExperimentUUID)
	if err != nil {
		log.Fatalf("Failed to move AA start time to now, simulation aborted: %v", err)
	}

	aaParticipantsWithActions := es.GetParticipantsWithActions(model.ExperimentPhaseAA)
	es.PerformAATest(*aaParticipantsWithActions)

	// Wait for the flush to take place so we are asserting against complete results
	fmt.Println("Waiting for buffer flush cooldown before performing assertions...")
	time.Sleep(time.Duration(70 * time.Second))

	fail := es.AssertionService.PerformAssertionsFor(es.ExperimentConfig.AA.PublishAmounts, es.ExperimentConfig.ExperimentKey, model.ExperimentPhaseAA)
	if fail {
		log.Fatalf("A/A Phase Failed - Simulation Aborted")
	}

	abParticipantsWithActions := es.GetParticipantsWithActions(model.ExperimentPhaseAB)
	es.PerformABTest(*abParticipantsWithActions)

	// Wait for the flush to take place so we are asserting against complete results
	fmt.Println("Waiting for buffer flush cooldown before performing assertions...")
	time.Sleep(time.Duration(70 * time.Second))

	fail = es.AssertionService.PerformAssertionsFor(es.ExperimentConfig.AB.PublishAmounts, es.ExperimentConfig.ExperimentKey, model.ExperimentPhaseAB)
	if fail {
		log.Fatalf("A/B Phase Failed")
	}

	fmt.Println("A/B Phase Complete!")
}

func (es *ExperimentSimulation) PerformAATest(aaParticipantsWithActions []model.ExperimentParticipant) {

	durationSeconds := es.ExperimentConfig.AA.DurationSeconds
	totalActions := getTotalActions(aaParticipantsWithActions)

	fmt.Println("This A/A Phase Will Produce a total of", totalActions, "actions across all variants and the A/A Phase will run for", durationSeconds, "seconds")
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

func (es *ExperimentSimulation) PerformABTest(abParticipantsWithActions []model.ExperimentParticipant) {
	durationSeconds := es.ExperimentConfig.AB.DurationSeconds
	totalActions := getTotalActions(abParticipantsWithActions)

	fmt.Println("This A/B Phase Will Produce a total of", totalActions, "actions across all variants and the A/B Phase will run for", durationSeconds, "seconds")
	fmt.Println("")

	fmt.Println("Are you ready to begin the experiment simulation? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		fmt.Println("Experiment simulation aborted.")
		return
	}

	fmt.Println("")
	fmt.Println("----A/B Phase In Progress----")
	fmt.Println("")

	err := es.ExperimentTimingService.ProgressExperimentToABPhase(es.ExperimentConfig.ExperimentUUID)
	if err != nil {
		log.Fatalf("Failed to progress experiment to AB phase: %v", err)
	}

	fmt.Println("Invalidating assignment cache...")
	if err := es.CacheInvalidationProducer.InvalidateExperiment(context.Background(), es.ExperimentConfig.ExperimentKey); err != nil {
		log.Fatalf("Failed to invalidate cache: %v", err)
	}

	if es.PerformActions(abParticipantsWithActions, durationSeconds, totalActions) {
		return
	}

	es.ExperimentTimingService.EndABPhase(es.ExperimentConfig.ExperimentUUID)

	fmt.Println("----Experiment simulation completed.----")
}

func (es *ExperimentSimulation) PerformActions(particpantsWithActions []model.ExperimentParticipant, durationSeconds int, totalActions int) bool {
	interval := time.Duration(float64(time.Second) * float64(durationSeconds) / float64(len(particpantsWithActions)))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timeout := time.After((time.Duration(durationSeconds) * time.Second) + (10 * time.Second))
	totalActionsPerformed := 0
	currentParticipant := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for {
		select {
		case <-timeout:
			wg.Wait()
			return true
		case <-ticker.C:
			if totalActionsPerformed >= totalActions {
				wg.Wait()
				return true
			}

			participant := particpantsWithActions[currentParticipant]
			wg.Add(1)
			go func(p model.ExperimentParticipant) {
				defer wg.Done()
				p.PerformActionsWithDelay(es.Performer)
			}(participant)
			mu.Lock()
			totalActionsPerformed += len(participant.Actions)
			fmt.Printf("Total Actions Performed: %d/%d\r", totalActionsPerformed+1, totalActions)
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

func (es *ExperimentSimulation) GetParticipantsWithActions(phase model.ExperimentPhaseType) *[]model.ExperimentParticipant {
	fmt.Println("Reading Experiment Config")

	controlVariantKey, err := es.GetControlVariantKey()
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	treatmentVariantKey, err := es.GetTreatmentVariantKey()
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	numExposureEventsForControlVariant, err := es.GetNumberOfEventTypeToPublishForVariantAndPhase("experiment_exposure", controlVariantKey, phase)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	numExposureEventsForTreatmentVariant, err := es.GetNumberOfEventTypeToPublishForVariantAndPhase("experiment_exposure", treatmentVariantKey, phase)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	variantToEventToAmountExcludingExposure := es.GetVariantToEventToAmountExcludingExposure(phase)

	fmt.Println("Fetching User IDs From Grpc")

	controlUserIds, err := es.UserIdService.GetXUserIdsWithinExperimentAndVariant(*numExposureEventsForControlVariant, es.ExperimentConfig.ExperimentKey, controlVariantKey)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	treatmentUserIds, err := es.UserIdService.GetXUserIdsWithinExperimentAndVariant(*numExposureEventsForTreatmentVariant, es.ExperimentConfig.ExperimentKey, treatmentVariantKey)
	if err != nil {
		log.Fatalf("Failed to complete simulation: %v \n", err)
	}

	userIds := append(controlUserIds, treatmentUserIds...)

	experimentParticipantsWithExposureEvents := es.createExperimentParticipantsWithExposureEventsForUserIds(userIds, es.ExperimentConfig.ExperimentKey)

	numExperimentParticipants := len(experimentParticipantsWithExposureEvents)

	addActions := func(variantKey model.VariantKey, offset int) {
		for eventKey, numOfEventToPublish := range variantToEventToAmountExcludingExposure[variantKey] {
			for i := 0; i < numOfEventToPublish; i++ {
				idx := offset + i
				if idx >= numExperimentParticipants {
					log.Fatalf("You defined %d experiment_exposure events but then asked for %d %s events - That's not possible! ", numExperimentParticipants, numOfEventToPublish, eventKey)
				}

				eventFields := es.ExperimentConfig.Events[eventKey].Fields
				fieldValues := make(map[model.EventField]interface{})

				for fieldName, fieldConfig := range eventFields {
					fieldSeed := DeriveSeed(es.ExperimentConfig.RandomSeed, string(variantKey), string(eventKey), string(fieldName), strconv.Itoa(i))
					randSouce := rand.New(rand.NewSource(fieldSeed))
					var phaseFieldConfig map[model.VariantKey]model.FieldConfigMinMax
					if phase == model.ExperimentPhaseAA {
						phaseFieldConfig = fieldConfig.AA
					} else {
						phaseFieldConfig = fieldConfig.AB
					}
					fieldValues[fieldName] = GenerateDataForField(fieldConfig.Type, phaseFieldConfig[variantKey], randSouce)
				}
				experimentParticipantsWithExposureEvents[idx].AddAction(model.NewParticipantEventParameters(eventKey, fieldValues))
			}
		}
	}

	addActions(controlVariantKey, 0)
	addActions(treatmentVariantKey, len(controlUserIds))

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

func (es *ExperimentSimulation) createExperimentParticipantsWithExposureEventsForUserIds(userIds []string, experimentKey string) []model.ExperimentParticipant {
	var experimentParticipants []model.ExperimentParticipant

	for _, userId := range userIds {
		participant := model.NewExperimentParticipant(userId, experimentKey)
		participant.AddAction(model.NewParticipantEventParameters("experiment_exposure", map[model.EventField]interface{}{}))
		experimentParticipants = append(experimentParticipants, participant)
	}

	return experimentParticipants
}

func (es *ExperimentSimulation) GetVariantToEventToAmountExcludingExposure(phase model.ExperimentPhaseType) map[model.VariantKey]map[model.EventKey]int {
	variantToEventKeyToAmount := make(map[model.VariantKey]map[model.EventKey]int)

	for variantKey, _ := range es.GetVariantKeys() {
		eventToAmountForVariant, err := es.GetNumberOfEventsToPublishForVariantAndPhase(variantKey, phase)
		if err != nil {
			log.Fatalf("Failed to complete simulation: %v \n", err)
		}

		variantToEventKeyToAmount[variantKey] = eventToAmountForVariant
	}

	return variantToEventKeyToAmount
}

func (es *ExperimentSimulation) GetNumberOfEventsToPublishForVariantAndPhase(variantKey model.VariantKey, phase model.ExperimentPhaseType) (map[model.EventKey]int, error) {
	eventToAmount := make(map[model.EventKey]int)

	var publishAmounts map[model.EventKey]map[model.VariantKey]int
	switch phase {
	case model.ExperimentPhaseAA:
		publishAmounts = es.ExperimentConfig.AA.PublishAmounts
	case model.ExperimentPhaseAB:
		publishAmounts = es.ExperimentConfig.AB.PublishAmounts
	}

	for eventKey, variantToAmountMap := range publishAmounts {
		if eventKey == "experiment_exposure" {
			continue
		}

		amount, exists := variantToAmountMap[variantKey]
		if !exists {
			return nil, fmt.Errorf("number of %s events to publish for variant_key %s in phase %s not in yml", eventKey, variantKey, phase)
		}

		eventToAmount[eventKey] = amount
	}

	return eventToAmount, nil
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

func (es *ExperimentSimulation) GetTreatmentVariantKey() (model.VariantKey, error) {
	for key, role := range es.ExperimentConfig.Variants {
		if role == model.Treatment {
			return key, nil
		}
	}

	return "", errors.New("no Treatment Variant Specified")
}

func (es *ExperimentSimulation) GetNumberOfEventTypeToPublishForVariantAndPhase(event_key model.EventKey, variant_key model.VariantKey, phase model.ExperimentPhaseType) (*int, error) {
	_, variantExists := es.ExperimentConfig.Variants[variant_key]
	if !variantExists {
		return nil, errors.New("variant Key Not Found in variants section")
	}

	var publishAmounts map[model.EventKey]map[model.VariantKey]int
	switch phase {
	case model.ExperimentPhaseAA:
		publishAmounts = es.ExperimentConfig.AA.PublishAmounts
	case model.ExperimentPhaseAB:
		publishAmounts = es.ExperimentConfig.AB.PublishAmounts
	default:
		return nil, fmt.Errorf("unknown phase: %s", phase)
	}

	variantAmounts, eventKeyExists := publishAmounts[event_key]
	if !eventKeyExists {
		return nil, fmt.Errorf("event_key %s not found in %s publish_amounts", event_key, phase)
	}

	amountToPublish, amountExists := variantAmounts[variant_key]
	if !amountExists {
		return nil, fmt.Errorf("publish_amount not found for variant %s in phase %s", variant_key, phase)
	}

	return &amountToPublish, nil

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
