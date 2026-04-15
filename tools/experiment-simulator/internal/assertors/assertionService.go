package assertors

import (
	"context"
	"experiment-simulator/internal/model"
	"experiment-simulator/internal/repository"
	"fmt"
)

type AssertionService interface {
	PerformAssertionsFor(publishAmounts map[model.EventKey]map[model.VariantKey]int, experimentKey string) (hasFailures bool)
}

type AssertionServiceClickhouse struct {
	eventsRepository repository.EventsRepository
}

func NewAssertionServiceClickhouse(eventsRepository repository.EventsRepository) *AssertionServiceClickhouse {
	return &AssertionServiceClickhouse{
		eventsRepository: eventsRepository,
	}
}

func (a *AssertionServiceClickhouse) PerformAssertionsFor(publishAmounts map[model.EventKey]map[model.VariantKey]int, experimentKey string) (hasFailures bool) {
	ctx := context.Background()

	results := make(map[model.EventKey]map[model.VariantKey]bool)

	fmt.Printf("---- Performing assertions for experiment key %s ----\n", experimentKey)

	for eventKey, variantKeys := range publishAmounts {
		if results[eventKey] == nil {
			results[eventKey] = make(map[model.VariantKey]bool)
		}

		for variantKey, amount := range variantKeys {
			count, err := a.eventsRepository.GetCountOfEventForVariantAndExperiment(ctx, eventKey, variantKey, experimentKey)
			if err != nil {
				fmt.Errorf("Failed to perform assertions for event key %s ", eventKey)
				continue
			}
			if amount != count {
				results[eventKey][variantKey] = false
				fmt.Printf("For event key %s and variant key %s expected amount was %d and actual amount was %d, assertion failed\n", eventKey, variantKey, amount, count)
			} else {
				results[eventKey][variantKey] = true
				fmt.Printf("For event key %s and variant key %s expected amount was %d and actual amount was %d, assertion passed\n", eventKey, variantKey, amount, count)
			}
		}
	}

	for _, variants := range results {
		for _, v := range variants {
			if !v {
				return true
			}
		}
	}
	return false
}
