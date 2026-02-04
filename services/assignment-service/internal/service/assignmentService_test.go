package service

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"
)

const salt = "ULTRA_SECRET_SALT"
const bucketCount = 10000

func TestAssignmentService_GetVariantsForUserId_ShouldUpdateCache(t *testing.T) {
	assignmentCache := NewStubAssignmentCache()
	assignmentClient := NewStubAssignmentClient()
	bucketService := NewBucketService(salt, bucketCount)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	assignmentService := NewAssignmentService(logger, bucketService, assignmentClient, assignmentCache)
	userId := "21"
	expectedBucket := int32(3930)

	assignmentsFromClient := map[string]string{
		"flag_c": "variant3",
		"flag_d": "variant4",
	}
	assignmentClient.SetAssignmentsForBucket(expectedBucket, assignmentsFromClient)

	_, err := assignmentService.GetVariantsForUserId(context.Background(), userId)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The cache is written to asynchronously as to not block the request, so wait a bit or the test may fail
	time.Sleep(5 * time.Second)

	cachedAssignments, err := assignmentCache.GetAssignmentsForBucket(context.Background(), expectedBucket)
	if err != nil {
		t.Fatalf("Unexpected error retrieving from cache: %v", err)
	}

	if !reflect.DeepEqual(cachedAssignments, assignmentsFromClient) {
		t.Errorf("Expected cached assignments %v, got %v", assignmentsFromClient, cachedAssignments)
	}
}

func TestAssignmentService_GetVariantsForUserId(t *testing.T) {

	tests := []struct {
		name                  string
		userId                string
		expectedBucket        int32
		assignmentsInCache    map[string]string
		assignmentsFromClient map[string]string
		expectedAssignments   map[string]string
	}{
		{
			name:           "Should return assignments from cache when available",
			userId:         "21",
			expectedBucket: 3930,
			assignmentsInCache: map[string]string{
				"flag_a": "variant1",
				"flag_b": "variant2",
			},
			assignmentsFromClient: nil,
			expectedAssignments: map[string]string{
				"flag_a": "variant1",
				"flag_b": "variant2",
			},
		}, {
			name:               "Should return assignments from client when not in cache",
			userId:             "21",
			expectedBucket:     3930,
			assignmentsInCache: nil, // No cache
			assignmentsFromClient: map[string]string{
				"flag_c": "variant3",
				"flag_d": "variant4",
			},
			expectedAssignments: map[string]string{
				"flag_c": "variant3",
				"flag_d": "variant4",
			},
		}, {
			name:                  "Should return empty assignments when neither cache nor client have data",
			userId:                "99999",
			expectedBucket:        9842,
			assignmentsInCache:    nil, // No cache
			assignmentsFromClient: nil, // No data from client
			expectedAssignments:   nil,
		},
	}

	assignmentCache := NewStubAssignmentCache()
	assignmentClient := NewStubAssignmentClient()
	bucketService := NewBucketService(salt, bucketCount)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	assignmentService := NewAssignmentService(logger, bucketService, assignmentClient, assignmentCache)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignmentCache.ClearCache()
			assignmentClient.ClearAssignments()

			if tt.assignmentsInCache != nil {
				err := assignmentService.assignmentCache.SetAssignmentsForBucket(context.Background(), tt.expectedBucket, tt.assignmentsInCache)
				if err != nil {
					t.Fatal("Failed to set up cache:", err)
				}
			}

			if tt.assignmentsFromClient != nil {
				assignmentClient.SetAssignmentsForBucket(tt.expectedBucket, tt.assignmentsFromClient)
			}

			assignments, err := assignmentService.GetVariantsForUserId(context.Background(), tt.userId)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(assignments, tt.expectedAssignments) {
				t.Errorf("Expected assignments %v, got %v", tt.expectedAssignments, assignments)
			}
		})
	}

}
