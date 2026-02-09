package service

//
//import (
//	"assignment-service/internal/model"
//	"context"
//	"io"
//	"log/slog"
//	"testing"
//	"time"
//)
//
//func TestListenForCacheInvalidation_shouldInvalidateGivenFullInvalidate(t *testing.T) {
//	cache := NewStubAssignmentCache()
//	consumer := NewMockKafkaConsumer()
//	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
//
//	service := NewCacheInvalidationServiceKafka(consumer, logger, cache)
//
//	existingCacheSet := map[int32]map[string]string{
//		1001: {
//			"feature_dark_mode":     "enabled",
//			"feature_new_dashboard": "control",
//			"experiment_checkout":   "variant_a",
//		},
//		1002: {
//			"feature_dark_mode":  "disabled",
//			"experiment_pricing": "variant_b",
//			"rollout_new_api":    "50_percent",
//		},
//		1003: {
//			"feature_notifications": "enabled",
//			"experiment_onboarding": "variant_c",
//		},
//	}
//
//	for bucketId, assignments := range existingCacheSet {
//		err := cache.SetAssignmentsForBucket(nil, bucketId, assignments)
//		if err != nil {
//			t.Fatalf("Failed to set up initial cache for bucket %d: %v", bucketId, err)
//		}
//	}
//
//	invalidationMessage := model.InvalidationMessage{
//		Action: "FULL_INVALIDATE",
//		Bucket: 1002,
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	consumer.AddMessage(invalidationMessage)
//
//	go func() {
//		_ = service.ListenForInvalidations(ctx)
//	}()
//
//	time.Sleep(200 * time.Millisecond)
//
//	assignments, err := cache.GetAssignmentsForBucket(nil, 1002)
//	if err != nil {
//		t.Fatalf("Error retrieving assignments for bucket 1002: %v", err)
//	}
//
//	if len(assignments) != 0 {
//		t.Errorf("Expected cache for bucket 1002 to be invalidated, but found assignments: %v", assignments)
//	}
//}
//
//func TestListenForCacheInvalidation_shouldUpdateGivenUpdateAction(t *testing.T) {
//	cache := NewStubAssignmentCache()
//	consumer := NewMockKafkaConsumer()
//	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
//
//	service := NewCacheInvalidationServiceKafka(consumer, logger, cache)
//
//	existingCacheSet := map[int32]map[string]string{
//		2001: {
//			"feature_search":    "enabled",
//			"experiment_signup": "variant_x",
//		},
//	}
//
//	for bucketId, assignments := range existingCacheSet {
//		err := cache.SetAssignmentsForBucket(nil, bucketId, assignments)
//		if err != nil {
//			t.Fatalf("Failed to set up initial cache for bucket %d: %v", bucketId, err)
//		}
//	}
//
//	invalidationMessage := model.InvalidationMessage{
//		Action:   "UPDATE",
//		Bucket:   2001,
//		Flag:     "experiment_signup",
//		NewValue: "variant_y",
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	consumer.AddMessage(invalidationMessage)
//
//	go func() {
//		_ = service.ListenForInvalidations(ctx)
//	}()
//
//	time.Sleep(200 * time.Millisecond)
//
//	assignments, err := cache.GetAssignmentsForBucket(nil, 2001)
//	if err != nil {
//		t.Fatalf("Error retrieving assignments for bucket 2001: %v", err)
//	}
//
//	expectedValue := "variant_y"
//	if assignments["experiment_signup"] != expectedValue {
//		t.Errorf("Expected 'experiment_signup' to be updated to '%s', but got '%s'", expectedValue, assignments["experiment_signup"])
//	}
//}
//
//func TestListenForCacheInvalidation_shouldRemoveGivenRemoveAction(t *testing.T) {
//	cache := NewStubAssignmentCache()
//	consumer := NewMockKafkaConsumer()
//	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
//
//	service := NewCacheInvalidationServiceKafka(consumer, logger, cache)
//
//	existingCacheSet := map[int32]map[string]string{
//		3001: {
//			"feature_beta_access": "enabled",
//			"experiment_ui_test":  "variant_z",
//		},
//	}
//
//	for bucketId, assignments := range existingCacheSet {
//		err := cache.SetAssignmentsForBucket(nil, bucketId, assignments)
//		if err != nil {
//			t.Fatalf("Failed to set up initial cache for bucket %d: %v", bucketId, err)
//		}
//	}
//
//	invalidationMessage := model.InvalidationMessage{
//		Action: "REMOVE",
//		Bucket: 3001,
//		Flag:   "experiment_ui_test",
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	consumer.AddMessage(invalidationMessage)
//
//	go func() {
//		_ = service.ListenForInvalidations(ctx)
//	}()
//
//	time.Sleep(200 * time.Millisecond)
//
//	assignments, err := cache.GetAssignmentsForBucket(nil, 3001)
//	if err != nil {
//		t.Fatalf("Error retrieving assignments for bucket 3001: %v", err)
//	}
//
//	if _, exists := assignments["experiment_ui_test"]; exists {
//		t.Errorf("Expected 'experiment_ui_test' to be removed from cache, but it still exists with value '%s'", assignments["experiment_ui_test"])
//	}
//}
//
//// TODO: Write tests for retry logic
