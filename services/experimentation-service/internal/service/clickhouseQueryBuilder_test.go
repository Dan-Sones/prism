package service

import (
	event2 "experimentation-service/internal/model/event"
	"fmt"
	"testing"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
)

func TestClickhouseQueryBuilder_BuildInEventKeyWhere(t *testing.T) {
	tests := []struct {
		name string
		m    metric.EnrichedMetric
		want string
	}{
		{
			name: "Single event key",
			m: metric.EnrichedMetric{
				MetricComponents: []metric.EnrichedMetricComponent{
					{
						EventType: event.EventType{
							EventKey: "eventA",
						},
					},
				},
			},
			want: "event_key IN ('eventA')",
		},
		{
			name: "Multiple event keys",
			m: metric.EnrichedMetric{
				MetricComponents: []metric.EnrichedMetricComponent{
					{
						EventType: event.EventType{
							EventKey: "eventA",
						},
					},
					{
						EventType: event.EventType{
							EventKey: "eventB",
						},
					},
					{
						EventType: event.EventType{
							EventKey: "eventC",
						},
					},
				},
			},
			want: "event_key IN ('eventA', 'eventB', 'eventC')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &ClickhouseQueryBuilder{}
			got := builder.BuildInEventKeyWhere(tt.m)
			if got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}

}

func TestClickhouseQueryBuilder_BuildSelectItemForCountDistinct(t *testing.T) {
	sysColName := "user_id"

	tests := []struct {
		name      string
		component metric.EnrichedMetricComponent
		want      string
	}{
		{
			name: "Count distinct on system column",
			component: metric.EnrichedMetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				SystemColumnName: &sysColName,
				Role:             metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(user_id, event_key = 'eventA') AS numerator",
		},
		{
			name: "Count distinct on string event field",
			component: metric.EnrichedMetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldA",
					DataType: event.DataTypeString,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(string_properties['fieldA'], event_key = 'eventA') AS numerator",
		},
		{
			name: "Count distinct on float event field",
			component: metric.EnrichedMetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldB",
					DataType: event.DataTypeFloat,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(float_properties['fieldB'], event_key = 'eventA') AS numerator",
		},
		{
			name: "Count distinct on int event field",
			component: metric.EnrichedMetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldC",
					DataType: event.DataTypeInt,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(int_properties['fieldC'], event_key = 'eventA') AS numerator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &ClickhouseQueryBuilder{}
			got, err := builder.BuildSelectItemForCountDistinct(tt.component)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestClickhouseQueryBuilder_BuildQueryForExperimentMetric_Ratio(t *testing.T) {
	sysColName := "user_id"

	startTime, err := time.Parse(time.RFC3339, "2026-04-23T14:30:00Z")
	if err != nil {
		t.Fatalf("Error parsing start time: %v", err)
	}
	endTime := startTime.Add(1 * time.Hour)

	tests := []struct {
		name          string
		experimentKey string
		m             metric.EnrichedMetric
		startTime     time.Time
		endTime       time.Time
		isAA          bool
		expectedQuery event2.QueryString
	}{
		{
			name:          "BINARY Ratio Metric System Column Name, AA: purchase conversion rate",
			experimentKey: "button_color_v1",
			startTime:     startTime,
			endTime:       endTime,
			isAA:          true,
			m: metric.EnrichedMetric{
				MetricType: metric.MetricTypeRatio,
				MetricComponents: []metric.EnrichedMetricComponent{
					{
						Role: metric.ComponentRoleNumerator,
						EventType: event.EventType{
							EventKey: "purchase",
						},
						AggregationOperation: metric.AggregationOperationCountDistinct,
						SystemColumnName:     &sysColName,
					},
					{
						Role: metric.ComponentRoleDenominator,
						EventType: event.EventType{
							EventKey: "experiment_exposure",
						},
						AggregationOperation: metric.AggregationOperationCountDistinct,
						SystemColumnName:     &sysColName,
					},
				},
			},
			expectedQuery: `SELECT variant_key, uniqExactIf(user_id, event_key = 'purchase') AS numerator, uniqExactIf(user_id, event_key = 'experiment_exposure') AS denominator FROM cooked_events WHERE experiment_key = 'button_color_v1' AND event_key IN ('purchase', 'experiment_exposure') AND sent_at >= '2026-04-23 14:30:00' AND sent_at <= '2026-04-23 15:30:00' AND is_aa = true GROUP BY variant_key;`,
		},
		{
			name:          "BINARY Ratio Metric System Column Name, NOT AA: purchase conversion rate",
			experimentKey: "button_color_v1",
			startTime:     startTime,
			endTime:       endTime,
			isAA:          false,
			m: metric.EnrichedMetric{
				MetricType: metric.MetricTypeRatio,
				MetricComponents: []metric.EnrichedMetricComponent{
					{
						Role: metric.ComponentRoleNumerator,
						EventType: event.EventType{
							EventKey: "purchase",
						},
						AggregationOperation: metric.AggregationOperationCountDistinct,
						SystemColumnName:     &sysColName,
					},
					{
						Role: metric.ComponentRoleDenominator,
						EventType: event.EventType{
							EventKey: "experiment_exposure",
						},
						AggregationOperation: metric.AggregationOperationCountDistinct,
						SystemColumnName:     &sysColName,
					},
				},
			},
			expectedQuery: `SELECT variant_key, uniqExactIf(user_id, event_key = 'purchase') AS numerator, uniqExactIf(user_id, event_key = 'experiment_exposure') AS denominator FROM cooked_events WHERE experiment_key = 'button_color_v1' AND event_key IN ('purchase', 'experiment_exposure') AND sent_at >= '2026-04-23 14:30:00' AND sent_at <= '2026-04-23 15:30:00' AND is_aa = false GROUP BY variant_key;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &ClickhouseQueryBuilder{}
			query, err := builder.BuildQueryFor(tt.experimentKey, tt.m, startTime, endTime, tt.isAA)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if query != tt.expectedQuery {
				t.Errorf("Expected query:\n%v\nGot:\n%v", tt.expectedQuery, query)
			}
		})
	}
}

func TestClickhouseQueryBuilder_BuildQueryFor_NoMetricComponents(t *testing.T) {
	experimentKey := "button_color_v1"
	m := metric.EnrichedMetric{
		MetricType: metric.MetricTypeRatio,
	}
	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour)

	builder := &ClickhouseQueryBuilder{}
	_, err := builder.BuildQueryFor(experimentKey, m, startTime, endTime, false)
	if err.Error() != "metric must have at least one component" {
		t.Fatalf("Expected error about missing metric components, got: %v", err)
	}
}

func TestClickhouseQueryBuilder_BuildTimeRangeWhere(t *testing.T) {
	builder := &ClickhouseQueryBuilder{}
	const layout = "2006-01-02 15:04:05"

	startTime := time.Now().UTC()
	endTime := startTime.UTC().Add(1 * time.Hour)

	expected := fmt.Sprintf("sent_at >= '%s' AND sent_at <= '%s'",
		startTime.Format(layout), endTime.Format(layout))

	got := builder.BuildTimeRangeWhere(startTime, endTime)
	if got != expected {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestClickhouseQueryBuilder_BuildIsAAWhere(t *testing.T) {
	builder := &ClickhouseQueryBuilder{}

	tests := []struct {
		name string
		isAA bool
		want string
	}{
		{
			name: "isAA true",
			isAA: true,
			want: "is_aa = true",
		},
		{
			name: "isAA false",
			isAA: false,
			want: "is_aa = false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.BuildIsAAWhere(tt.isAA)
			if got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}
}
