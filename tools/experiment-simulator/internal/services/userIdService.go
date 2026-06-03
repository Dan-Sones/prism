package services

import (
	"context"
	"experiment-simulator/internal/clients"
	"experiment-simulator/internal/model"
	"fmt"
	"slices"
	"time"

	"github.com/Dan-Sones/prismhash"
	prismmodel "github.com/Dan-Sones/prismhash/model"
)

type UserIdService interface {
	GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error)
}

type UserIdServiceImp struct {
	BucketService                   *prismhash.BucketService
	ExperimentationAssignmentClient clients.ExperimentationAssignmentClient
}

func NewUserIdService(
	bucketService *prismhash.BucketService,
	experimentationAssignmentClient clients.ExperimentationAssignmentClient,
) *UserIdServiceImp {
	return &UserIdServiceImp{
		BucketService:                   bucketService,
		ExperimentationAssignmentClient: experimentationAssignmentClient,
	}
}

func (u *UserIdServiceImp) GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error) {

	foundIds := []string{}
	bucketCache := make(map[int][]prismmodel.ExperimentWithVariants)
	curr := 1
	for {
		if len(foundIds) == count {
			break
		}

		id := fmt.Sprintf("user-%d", curr)
		bucketForUser := int(u.BucketService.GetBucketFor(id))

		allExperiments, ok := bucketCache[bucketForUser]
		if !ok {
			fetched, err := u.ExperimentationAssignmentClient.GetExperimentsAndVariantsForBucketAtTime(context.Background(), bucketForUser, "experiment-simulator", time.Now())
			if err != nil {
				return nil, fmt.Errorf("failed to get experiment assignments for bucket: %w", err)
			}
			bucketCache[bucketForUser] = fetched
			allExperiments = fetched
		}

		activeExperiment := slices.Collect(func(yield func(prismmodel.ExperimentWithVariants) bool) {
			for _, experiment := range allExperiments {
				if experiment.ExperimentKey == experimentKey {
					yield(experiment)
				}
			}
		})

		if len(activeExperiment) == 0 {
			curr++
			continue
		}

		// filter down all experiments to just the experiment we care about
		variantKey, err := prismhash.GetVariantForExperiment(activeExperiment[0], id)
		if err != nil {
			return nil, fmt.Errorf("failed to get variant for user from experiment details: %w", err)
		}

		if variantKey == string(wantVariantKey) {
			foundIds = append(foundIds, id)
		}
		fmt.Printf("User IDs Found for variant key '%s': %d/%d\r", wantVariantKey, len(foundIds), count)

		curr++
	}

	return foundIds, nil
}
