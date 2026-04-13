package services

import (
	"context"
	"experiment-simulator/internal/clients"
	"experiment-simulator/internal/model"
	"fmt"
)

type UserIdService interface {
	GetXUserIdsWithinExperimentAndVariant(count int, experimentKey string, wantVariantKey model.VariantKey) ([]string, error)
}

type UserIdServiceImp struct {
	AssignmentClient clients.AssignmentClient
}

func NewUserIdService(client clients.AssignmentClient) *UserIdServiceImp {
	return &UserIdServiceImp{
		AssignmentClient: client,
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

		res, err := u.AssignmentClient.GetExperimentsAndVariantsForUser(context.Background(), id)
		if err != nil {
			return nil, fmt.Errorf("failed to get experiments and variants for id %s: %w", id, err)
		}

		variantKey, exists := res.Assignments[experimentKey]
		if !exists {
			continue
		}

		if variantKey == string(wantVariantKey) {
			foundIds = append(foundIds, id)
		}
		fmt.Printf("User IDs Found: %d/%d\r", len(foundIds), count)

		curr++
	}

	return foundIds, nil
}
