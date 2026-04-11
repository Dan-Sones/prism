package validators

import (
	metricreq "experimentation-service/internal/model/metric"
	"experimentation-service/internal/problems"
	"strings"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

func TestValidateCreateMetricRequest(t *testing.T) {
	tests := []struct {
		name    string
		request metricreq.CreateMetricRequest
		want    []problems.Violation
	}{
		{
			name: "Valid request",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				MetricKey:    "revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: nil,
		},
		{
			name: "Missing name",
			request: metricreq.CreateMetricRequest{
				Name:         "",
				MetricKey:    "revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
			},
		},
		{
			name: "Name exceeds max length",
			request: metricreq.CreateMetricRequest{
				Name:         strings.Repeat("a", 101),
				MetricKey:    "revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 100 characters",
				},
			},
		},
		{
			name: "Metric key exceeds max length",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				MetricKey:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "metric_key",
					Message: "Metric key must be less than 50 characters",
				},
			},
		},
		{
			name: "Missing metric key",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				MetricKey:    "",
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "metric_key",
					Message: "Metric key is required",
				},
			},
		},
		{
			name: "Missing analysis unit",
			request: metricreq.CreateMetricRequest{
				Name:      "Revenue",
				MetricKey: "revenue",
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "analysis_unit",
					Message: "Analysis unit is required",
				},
			},
		},
		{
			name: "Metric key starts with number",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				MetricKey:    "123_revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "metric_key",
					Message: "Metric key must start with a letter and only contain letters, numbers, underscores, or hyphens",
				},
			},
		},
		{
			name: "Metric key contains invalid characters",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				MetricKey:    "revenue@value",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components: []metricreq.CreateMetricRequestComponent{
					{
						Role:                 metric.ComponentRoleBaseEvent,
						EventTypeID:          uuid.New(),
						FieldKeyID:           uuid.New(),
						AggregationOperation: metric.AggregationOperationCount,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "metric_key",
					Message: "Metric key must start with a letter and only contain letters, numbers, underscores, or hyphens",
				},
			},
		},
		{
			name: "No components",
			request: metricreq.CreateMetricRequest{
				Name:         "Revenue",
				MetricKey:    "revenue",
				AnalysisUnit: metric.AnalysisUnitUser,
				Components:   []metricreq.CreateMetricRequestComponent{},
			},
			want: []problems.Violation{
				{
					Field:   "components",
					Message: "At least one component is required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCreateMetricRequest(tt.request)

			if len(got) != len(tt.want) {
				t.Fatalf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}

			for i, violation := range got {
				if violation != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], violation)
				}
			}
		})
	}
}

func TestValidateCreateMetricRequestComponent(t *testing.T) {
	tests := []struct {
		name      string
		component metricreq.CreateMetricRequestComponent
		index     int
		want      []problems.Violation
	}{
		{
			name: "Valid component",
			component: metricreq.CreateMetricRequestComponent{
				Role:                 metric.ComponentRoleBaseEvent,
				EventTypeID:          uuid.New(),
				FieldKeyID:           uuid.New(),
				AggregationOperation: metric.AggregationOperationCount,
			},
			index: 0,
			want:  nil,
		},
		{
			name: "Missing role",
			component: metricreq.CreateMetricRequestComponent{
				Role:                 "",
				EventTypeID:          uuid.New(),
				FieldKeyID:           uuid.New(),
				AggregationOperation: metric.AggregationOperationCount,
			},
			index: 0,
			want: []problems.Violation{
				{
					Field:   "components[0].role",
					Message: "Role is required",
				},
			},
		},
		{
			name: "Missing event type ID",
			component: metricreq.CreateMetricRequestComponent{
				Role:                 metric.ComponentRoleBaseEvent,
				EventTypeID:          uuid.Nil,
				FieldKeyID:           uuid.New(),
				AggregationOperation: metric.AggregationOperationCount,
			},
			index: 0,
			want: []problems.Violation{
				{
					Field:   "components[0].event_type_id",
					Message: "Event type ID is required",
				},
			},
		},
		{
			name: "Missing field key",
			component: metricreq.CreateMetricRequestComponent{
				Role:                 metric.ComponentRoleBaseEvent,
				EventTypeID:          uuid.New(),
				FieldKeyID:           uuid.Nil,
				AggregationOperation: metric.AggregationOperationCount,
			},
			index: 0,
			want: []problems.Violation{
				{
					Field:   "components[0].field_key_id",
					Message: "Field key is required",
				},
			},
		},
		{
			name: "Missing aggregation operation",
			component: metricreq.CreateMetricRequestComponent{
				Role:                 metric.ComponentRoleBaseEvent,
				EventTypeID:          uuid.New(),
				FieldKeyID:           uuid.New(),
				AggregationOperation: "",
			},
			index: 0,
			want: []problems.Violation{
				{
					Field:   "components[0].aggregation_operation",
					Message: "Aggregation operation is required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCreateMetricRequestComponent(tt.component, tt.index)

			if len(got) != len(tt.want) {
				t.Fatalf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}

			for i, violation := range got {
				if violation != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], violation)
				}
			}
		})
	}
}
