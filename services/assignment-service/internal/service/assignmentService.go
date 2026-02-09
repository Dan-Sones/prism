package service

import (
	"assignment-service/internal/clients"
	"assignment-service/internal/model"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
)

type AssignmentService struct {
	logger           *slog.Logger
	bucketService    *BucketService
	assignmentCache  AssignmentCache
	assignmentClient clients.AssignmentClient
}

func NewAssignmentService(logger *slog.Logger, bService *BucketService, assignmentClient clients.AssignmentClient, assignmentCache AssignmentCache) *AssignmentService {
	return &AssignmentService{
		logger:           logger.With(slog.String("component", "AssignmentService")),
		bucketService:    bService,
		assignmentCache:  assignmentCache,
		assignmentClient: assignmentClient,
	}
}

func (e *AssignmentService) GetAssignmentsForUserId(ctx context.Context, userId string) (map[string]string, error) {
	bucket := e.bucketService.GetBucketFor(userId)

	// REQUIRED CHANGES:
	// 1. Make a request (worry about caching later) to get the active experiments for a bucket.
	// within this response there should be: experiment_id, salt, "flag_id", a list of varisants with their variant_id, flag_value and traffic allocation which is represented as lower_bound and upper_bound for each variant.
	// we should assume that the upper_bound and lower_bound are between 0 and 100, and that is is not possible for there to be overlapping in these traffic allocations.
	// 2. perform an additional md5 hashing step on the user_id + experiment_id + salt. then mod this hash value by 100 to get a number between 0 and 100.
	// 3. compare the number from step 2 with the traffic allocation for each variant to determine which variant the user is bucketed into for each experiment.
	// 4. return a map of flag_id to the flag_value for the variant that the user is bucketed into for each experiment.

	// Will get it working as is first, and then will add back in the caching layer. I need to consider what we would cache in this new instance.
	// I guess just the response from the gRPC call. The actual hash can then be done on the fly without the need to re-request bounds etc.
	// We still need the same operations in the cache, but I am not sure on the structure of the cache....
	// we currently use hmap in redis, but we would now be storing an entire object as the value, I'm not sure if this is an anti-pattern or not?

	experiments, err := e.assignmentClient.GetExperimentsAndVariantsForBucket(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiments and variants for bucket %d: %w", bucket, err)
	}
	e.logger.Info("Fetched assignments from gRPC", "bucket", bucket)

	assignments := make(map[string]string)
	for _, experiment := range experiments {
		variant, err := e.getVariantForExperiment(experiment, userId)
		if err != nil {
			e.logger.Error("Failed to get variant for experiment", "experiment", experiment.ExperimentKey, "user_id", userId, "error", err)
			continue
		}
		assignments[experiment.ExperimentKey] = variant
	}

	return assignments, nil
}

func (e *AssignmentService) getVariantForExperiment(experiments model.ExperimentWithVariants, userId string) (string, error) {
	numberLinePosition := e.getNumberLinePositionForUserAndExperiment(userId, experiments.ExperimentKey, experiments.UniqueSalt)

	for _, variant := range experiments.Variants {
		if numberLinePosition >= variant.LowerBound && numberLinePosition <= variant.UpperBound {
			return variant.VariantKey, nil
		}
	}
	return "", fmt.Errorf("no variant found for user %s in experiment %s with number line position %d", userId, experiments.ExperimentKey, numberLinePosition)
}

func (e *AssignmentService) getNumberLinePositionForUserAndExperiment(userId, experimentId, uniqueSalt string) int32 {
	toHash := fmt.Sprintf("%s:%s:%s", userId, experimentId, uniqueSalt)
	hash := md5.Sum([]byte(toHash))

	hashHex := hex.EncodeToString(hash[:])

	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	numberLinePosition := new(big.Int)
	return int32(numberLinePosition.Mod(hashInt, big.NewInt(100)).Int64())
}
