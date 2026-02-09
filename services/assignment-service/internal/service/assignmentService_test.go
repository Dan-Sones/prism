package service

import (
	"assignment-service/internal/model"
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
)

const salt = "ULTRA_SECRET_SALT"
const bucketCount = 10000

func TestAssignmentService_GetAssignmentsForUserId_shouldAttemptToReadFromCacheFirst(t *testing.T) {
	experimentClient := NewStubExperimentClient()
	experimentCache := NewStubExperimentCache()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	assignmentService := NewAssignmentService(logger, NewBucketService(salt, bucketCount), experimentClient, experimentCache)

	bucketId := int32(3930)

	experimentKey := "button_color_v1"
	uniqueSalt := "4e770d52-b2b0-42ed-8ccb-2321ff48e143"

	err := experimentCache.SetExperiment(context.Background(), experimentKey, &model.ExperimentWithVariants{
		ExperimentKey: experimentKey,
		UniqueSalt:    uniqueSalt,
		Variants: []model.Variant{
			{
				VariantKey: "button_blue",
				LowerBound: 0,
				UpperBound: 49,
			},
			{
				VariantKey: "button_green",
				LowerBound: 50,
				UpperBound: 99,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to set experiment: %v", err)
	}

	err = experimentCache.AddBucketExperimentKey(context.Background(), bucketId, experimentKey)
	if err != nil {
		t.Fatalf("failed to add experiment key to bucket: %v", err)
	}

	assignments, err := assignmentService.GetAssignmentsForUserId(context.Background(), "21")
	if err != nil {
		t.Fatalf("failed to get assignments: %v", err)
	}

	expectedVariant := "button_green"
	if assignments[experimentKey] != expectedVariant {
		t.Errorf("Expected variant %s for experiment %s, got %s", expectedVariant, experimentKey, assignments[experimentKey])
	}

}

func TestAssignmentService_GetAssignmentsForUserId_shouldCallGrpcOnCacheMiss(t *testing.T) {
	experimentClient := NewStubExperimentClient()
	experimentCache := NewStubExperimentCache()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	assignmentService := NewAssignmentService(logger, NewBucketService(salt, bucketCount), experimentClient, experimentCache)

	bucketId := int32(3930)

	experimentKey := "button_color_v1"
	uniqueSalt := "4e770d52-b2b0-42ed-8ccb-2321ff48e143"

	experimentClient.SetExperimentsForBucket(bucketId, []model.ExperimentWithVariants{{
		ExperimentKey: experimentKey,
		UniqueSalt:    uniqueSalt,
		Variants: []model.Variant{
			{
				VariantKey: "button_blue",
				LowerBound: 0,
				UpperBound: 49,
			},
			{
				VariantKey: "button_green",
				LowerBound: 50,
				UpperBound: 99,
			},
		},
	}})

	assignments, err := assignmentService.GetAssignmentsForUserId(context.Background(), "21")
	if err != nil {
		t.Fatalf("failed to get assignments: %v", err)
	}

	expectedVariant := "button_green"
	if assignments[experimentKey] != expectedVariant {
		t.Errorf("Expected variant %s for experiment %s, got %s", expectedVariant, experimentKey, assignments[experimentKey])
	}
}

func TestAssignmentService_GetAssignmentsForUserId_shouldUpdateCacheOnMissAndGrpcCall(t *testing.T) {
	experimentClient := NewStubExperimentClient()
	experimentCache := NewStubExperimentCache()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	assignmentService := NewAssignmentService(logger, NewBucketService(salt, bucketCount), experimentClient, experimentCache)

	bucketId := int32(3930)

	experimentKey := "button_color_v1"
	uniqueSalt := "4e770d52-b2b0-42ed-8ccb-2321ff48e143"

	experiment := model.ExperimentWithVariants{
		ExperimentKey: experimentKey,
		UniqueSalt:    uniqueSalt,
		Variants: []model.Variant{
			{
				VariantKey: "button_blue",
				LowerBound: 0,
				UpperBound: 49,
			},
			{
				VariantKey: "button_green",
				LowerBound: 50,
				UpperBound: 99,
			},
		},
	}

	experimentClient.SetExperimentsForBucket(bucketId, []model.ExperimentWithVariants{experiment})

	_, err := assignmentService.GetAssignmentsForUserId(context.Background(), "21")
	if err != nil {
		t.Fatalf("failed to get assignments: %v", err)
	}

	// The cache is written to asynchronously as to not block the request, so wait a bit or the test may fail
	time.Sleep(5 * time.Second)

	experiments, err := experimentCache.GetBucketExperimentKeys(context.Background(), bucketId)
	if err != nil {
		return
	}

	fmt.Printf("Cached experiment keys for bucket %d: %v\n", bucketId, experiments)

	if len(experiments) != 1 || experiments[0] != experimentKey {
		t.Errorf("Expected experiment key %s to be cached for bucket %d, got %v", experimentKey, bucketId, experiments)
	}

}
