package assertors

import (
	"context"
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/repository"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type AssertionService struct {
	eventsRepository repository.EventsRepository
}

func NewAssertionService(eventsRepository repository.EventsRepository) *AssertionService {
	return &AssertionService{
		eventsRepository: eventsRepository,
	}
}

func (a *AssertionService) WaitForFlush() {
	flushTimeStr := os.Getenv("DATA_COOKING_SERVICE_MICROBATCH_FLUSH_TIMEOUT_SECONDS")

	flushTime, err := strconv.Atoi(flushTimeStr)
	if err != nil {
		log.Fatalf("Error converting flush time to int: %v", err)
	}

	timeToWait := (time.Duration(flushTime) * time.Second) + (5 * time.Second) // Add an extra 5 seconds to be safe
	time.Sleep(timeToWait)
}

func (a *AssertionService) PerformAssertionsFor(experimentSimulation *model.ExperimentSimulation) {
	ctx := context.Background()

	for _, variantKey := range experimentSimulation.ExperimentConfig.ExperimentConfigAB.VariantKeys {
		for eventTypeKey := range experimentSimulation.ExperimentConfig.ExperimentConfigAB.Events {
			expected := experimentSimulation.GetTotalEventsForVariantAndEventType(variantKey, eventTypeKey)

			got, err := a.eventsRepository.GetCountOfEventForVariantAndExperiment(ctx, eventTypeKey, variantKey, experimentSimulation.ExperimentConfig.ExperimentConfigAB.FeatureFlagKey)
			if err != nil {
				log.Fatal("Error getting count of event for variant and experiment: ", err)
			}
			fmt.Printf("Assertion for experiment %s, variant %s, event %s: expected %d, got %d, assertion passed: %t\n", experimentSimulation.ExperimentConfig.ExperimentConfigAB.FeatureFlagKey, variantKey, eventTypeKey, expected, got, expected == got)
		}
	}
}
