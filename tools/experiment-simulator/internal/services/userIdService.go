package services

import (
	"context"
	"experiment-simulator/internal/clients"
	"experiment-simulator/internal/model"
	"fmt"
	"slices"
	"time"
)

type UserIdService interface {
	GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error)
}

type UserIdServiceImp struct {
	AssignmentClient                clients.AssignmentClient
	ExperimentationAssignmentClient clients.ExperimentationAssignmentClient
}

func NewUserIdService(
	assignmentClient clients.AssignmentClient,
	experimentationAssignmentClient clients.ExperimentationAssignmentClient,

) *UserIdServiceImp {
	return &UserIdServiceImp{
		AssignmentClient:                assignmentClient,
		ExperimentationAssignmentClient: experimentationAssignmentClient,
	}
}

func (u *UserIdServiceImp) GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error) {

	foundIds := []string{}
	curr := 1
	for {
		if len(foundIds) == count {
			break
		}

		id := fmt.Sprintf("user-%d", curr)

		bucketForUser, err := u.AssignmentClient.GetBucketForUser(context.Background(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get bucket for user: %w", err)
		}

		allExperiments, err := u.ExperimentationAssignmentClient.GetExperimentsAndVariantsForBucketAtTime(context.Background(), bucketForUser, "experiment-simulator", time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to get experiment assignments for bucket: %w", err)
		}

		activeExperiment := slices.Collect(func(yield func(clients.ExperimentWithVariants) bool) {
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
		variantKey, err := u.AssignmentClient.GetVariantForUserFromExperimentDetails(context.Background(), id, activeExperiment[0])
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
