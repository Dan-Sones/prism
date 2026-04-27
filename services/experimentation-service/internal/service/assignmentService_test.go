package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type mockExperimentRepository struct {
	getFunc func(ctx context.Context, bucketId int32, atTime time.Time) ([]*experiment.ExperimentWithVariants, error)
}

func (m *mockExperimentRepository) GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int32, atTime time.Time) ([]*experiment.ExperimentWithVariants, error) {
	return m.getFunc(ctx, bucketId, atTime)
}

func TestGetExperimentsAndVariantsForBucket_RepoErrorsShouldPropagate(t *testing.T) {
	mockRepo := &mockExperimentRepository{
		getFunc: func(ctx context.Context, bucketId int32, atTime time.Time) ([]*experiment.ExperimentWithVariants, error) {
			return nil, errors.New("repo error")
		},
	}
	service := &AssignmentService{
		experimentRepository: mockRepo,
		bucketCount:          100,
		logger:               slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	_, _, err := service.GetExperimentsAndVariantsForBucketAtTime(context.Background(), 10, "assignment-service", time.Now())
	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "repo error" {
		t.Errorf("Expected 'repo error', got %v", err)
	}
}

func GetNewAssignmentServiceForTest(bucketCount int32) *AssignmentService {
	return &AssignmentService{
		bucketCount: bucketCount,
	}
}

func TestGetExperimentsAndVariantsForBucket_Validation(t *testing.T) {
	tests := []struct {
		name          string
		bucketId      int32
		wantField     string
		wantViolation bool
	}{
		{
			name:          "Negative bucketId returns validation violation",
			bucketId:      -1,
			wantField:     "bucket_id",
			wantViolation: true,
		},
		{
			name:          "bucketId exceeds bucketCount returns validation violation",
			bucketId:      101,
			wantField:     "bucket_id",
			wantViolation: true,
		},
	}

	service := GetNewAssignmentServiceForTest(100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violations, err := service.GetExperimentsAndVariantsForBucketAtTime(context.Background(), tt.bucketId, "assignment-service", time.Now())
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if tt.wantViolation {
				if len(violations) == 0 {
					t.Fatal("Expected violations, got none")
				}
				found := false
				for _, v := range violations {
					if v.Field == tt.wantField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected violation on field %q, got %v", tt.wantField, violations)
				}
			}
		})
	}
}

func TestGetExperimentsAndVariantsForBucket_AAOverride(t *testing.T) {
	mockRepo := &mockExperimentRepository{
		getFunc: func(ctx context.Context, bucketId int32, atTime time.Time) ([]*experiment.ExperimentWithVariants, error) {

			exp1Id, err := uuid.NewUUID()
			if err != nil {
				t.Fatalf("Failed to generate UUID: %v", err)
			}

			return []*experiment.ExperimentWithVariants{
				{
					Experiment: experiment.Experiment{
						ID:          exp1Id,
						Name:        "Experiment 1",
						AAStartTime: time.Now().Add(-1 * time.Hour),
						AAEndTime:   time.Now().Add(1 * time.Hour),
						UniqueSalt:  "abc123",
					},
					Variants: []experiment.ExperimentVariant{
						{
							VariantKey:  "control",
							VariantType: experiment.VariantTypeControl,
							UpperBound:  50,
							LowerBound:  0,
						},
						{
							VariantKey:  "treatment",
							VariantType: experiment.VariantTypeTreatment,
							UpperBound:  100,
							LowerBound:  50,
						},
					},
				},
			}, nil
		},
	}

	service := &AssignmentService{
		experimentRepository: mockRepo,
		bucketCount:          100,
	}

	results, _, err := service.GetExperimentsAndVariantsForBucketAtTime(context.Background(), 10, "assignment-service", time.Now())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 experiment, got %d", len(results))
	}

	exp := results[0]

	// Assert that the treatment variant was stripped out and the control variant was adjusted to 100% of the traffic
	if len(exp.Variants) != 1 {
		t.Fatalf("Expected 1 variant after AA override, got %d", len(exp.Variants))
	}

	controlVariant := exp.Variants[0]
	if controlVariant.VariantType != experiment.VariantTypeControl {
		t.Errorf("Expected control variant type, got %v", controlVariant.VariantType)
	}
	if controlVariant.UpperBound != 100 {
		t.Errorf("Expected control variant upper bound to be 100, got %d", controlVariant.UpperBound)
	}
	if controlVariant.LowerBound != 0 {
		t.Errorf("Expected control variant lower bound to be 0, got %d", controlVariant.LowerBound)
	}
}
