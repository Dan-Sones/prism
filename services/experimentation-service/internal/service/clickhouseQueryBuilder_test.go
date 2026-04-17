package service

import "testing"

func TestClickhouseQueryBuilder_BuildInEventKeyWhere(t *testing.T) {
	tests := []struct {
		name      string
		eventKeys []string
		want      string
	}{
		{
			name:      "Single event key",
			eventKeys: []string{"eventA"},
			want:      "event_key in ('eventA')",
		},
		{
			name:      "Multiple event keys",
			eventKeys: []string{"eventA", "eventB", "eventC"},
			want:      "event_key in ('eventA', 'eventB', 'eventC')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &ClickhouseQueryBuilder{}
			got := builder.BuildInEventKeyWhere(tt.eventKeys)
			if got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}

}
