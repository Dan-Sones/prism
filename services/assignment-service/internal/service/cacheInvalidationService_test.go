package service

import (
	"assignment-service/internal/model"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestCacheInvalidationServiceKafka_ListenForInvalidations_shouldRemoveExperimentAndBucket(t *testing.T) {
	cache := NewStubExperimentCache()
	consumer := NewMockKafkaConsumer()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	service := NewCacheInvalidationServiceKafka(consumer, logger, cache)

	err := cache.SetExperiment(context.Background(), "checkout_experiment", &model.ExperimentWithVariants{
		ExperimentKey: "checkout_experiment",
		UniqueSalt:    "some_salt",
		Variants: []model.Variant{
			{
				VariantKey: "variant_a",
				LowerBound: 0,
				UpperBound: 49,
			},
			{
				VariantKey: "variant_b",
				LowerBound: 50,
				UpperBound: 99,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to set experiment in cache: %v", err)
	}
	err = cache.AddBucketExperimentKey(context.Background(), 1001, "checkout_experiment")
	if err != nil {
		t.Fatalf("Failed to add experiment in cache: %v", err)
	}
	err = cache.AddBucketExperimentKey(context.Background(), 1002, "checkout_experiment")
	if err != nil {
		t.Fatalf("Failed to add experiment in cache: %v", err)
	}

	removeMessage := model.ExperimentRemoveMessage{
		ExperimentKey: "checkout_experiment",
		Buckets:       []int32{1001, 1002},
	}

	removeMessageBytes, err := json.Marshal(removeMessage)
	if err != nil {
		t.Fatalf("Failed to marshal remove message: %v", err)
	}

	invalidationMessage := model.InvalidationMessage{
		Action: "REMOVE",
		Data:   removeMessageBytes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	consumer.AddMessage(invalidationMessage)

	go func() {
		_ = service.ListenForInvalidations(ctx)
	}()

	time.Sleep(5 * time.Second)

	experiment, err := cache.GetExperiment(context.Background(), "checkout_experiment")
	if err != nil {
		t.Fatalf("Error retrieving experiment from cache: %v", err)
	}
	if experiment != nil {
		t.Errorf("Expected experiment 'checkout_experiment' to be removed from cache, but it still exists")
	}

	for _, bucketId := range []int32{1001, 1002} {
		experimentKeys, err := cache.GetBucketExperimentKeys(context.Background(), bucketId)
		if err != nil {
			t.Fatalf("Error retrieving bucket experiment keys for bucket %d: %v", bucketId, err)
		}
		for _, key := range experimentKeys {
			if key == "checkout_experiment" {
				t.Errorf("Expected 'checkout_experiment' to be removed from bucket %d, but it still exists in bucket experiment keys", bucketId)
			}
		}
	}
}

func TestCacheInvalidationServiceKafka_ListenForInvalidations_shouldHandleActionUpdate(t *testing.T) {
	cache := NewStubExperimentCache()
	consumer := NewMockKafkaConsumer()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	service := NewCacheInvalidationServiceKafka(consumer, logger, cache)

	err := cache.SetExperiment(context.Background(), "checkout_experiment", &model.ExperimentWithVariants{
		ExperimentKey: "checkout_experiment",
		UniqueSalt:    "some_salt",
		Variants: []model.Variant{
			{
				VariantKey: "variant_a",
				LowerBound: 0,
				UpperBound: 49,
			},
			{
				VariantKey: "variant_b",
				LowerBound: 50,
				UpperBound: 99,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to set experiment in cache: %v", err)
	}
	err = cache.AddBucketExperimentKey(context.Background(), 1001, "checkout_experiment")
	if err != nil {
		t.Fatalf("Failed to add experiment in cache: %v", err)
	}

	updateMessage := model.ExperimentUpdateMessage{
		ExperimentKey: "checkout_experiment",
		NewExperiment: model.ExperimentWithVariants{
			ExperimentKey: "checkout_experiment",
			UniqueSalt:    "some_salt",
			Variants: []model.Variant{
				{
					VariantKey: "variant_a",
					LowerBound: 0,
					UpperBound: 24,
				},
				{
					VariantKey: "variant_b",
					LowerBound: 25,
					UpperBound: 100,
				},
			},
		},
	}

	updateMessageBytes, err := json.Marshal(updateMessage)
	if err != nil {
		t.Fatalf("Failed to marshal update message: %v", err)
	}

	invalidationMessage := model.InvalidationMessage{
		Action: "UPDATE",
		Data:   updateMessageBytes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	consumer.AddMessage(invalidationMessage)

	go func() {
		_ = service.ListenForInvalidations(ctx)
	}()

	time.Sleep(5 * time.Second)

	experiment, err := cache.GetExperimentsForBucket(context.Background(), 1001)
	if err != nil {
		t.Fatalf("Error retrieving experiments for bucket 1001: %v", err)
	}

	if experiment[0].Variants[0].UpperBound != 24 || experiment[0].Variants[1].LowerBound != 25 {
		t.Errorf("Experiment variants were not updated correctly in cache")
	}
}

// TODO: Write tests for retry logic
