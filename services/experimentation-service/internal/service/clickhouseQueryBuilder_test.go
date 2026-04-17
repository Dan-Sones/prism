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
