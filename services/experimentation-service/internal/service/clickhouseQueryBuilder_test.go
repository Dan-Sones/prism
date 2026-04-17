package service

import (
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
)

func TestClickhouseQueryBuilder_BuildInEventKeyWhere(t *testing.T) {
	tests := []struct {
		name string
		m    metric.Metric
		want string
	}{
		{
			name: "Single event key",
			m: metric.Metric{
				MetricComponents: []metric.MetricComponent{
					{
						EventType: event.EventType{
							EventKey: "eventA",
						},
					},
				},
			},
			want: "event_key in ('eventA')",
		},
		{
			name: "Multiple event keys",
			m: metric.Metric{
				MetricComponents: []metric.MetricComponent{
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
			want: "event_key in ('eventA', 'eventB', 'eventC')",
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
		component metric.MetricComponent
		want      string
	}{
		{
			name: "Count distinct on system column",
			component: metric.MetricComponent{
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
			component: metric.MetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldA",
					DataType: event.DataTypeString,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(string_properties[fieldA], event_key = 'eventA') AS numerator",
		},
		{
			name: "Count distinct on float event field",
			component: metric.MetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldB",
					DataType: event.DataTypeFloat,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(float_properties[fieldB], event_key = 'eventA') AS numerator",
		},
		{
			name: "Count distinct on int event field",
			component: metric.MetricComponent{
				EventType: event.EventType{
					EventKey: "eventA",
				},
				AggregationField: &event.EventField{
					FieldKey: "fieldC",
					DataType: event.DataTypeInt,
				},
				Role: metric.ComponentRoleNumerator,
			},
			want: "uniqExactIf(int_properties[fieldC], event_key = 'eventA') AS numerator",
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
