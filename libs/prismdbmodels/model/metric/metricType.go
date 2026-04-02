package metric

import "fmt"

type MetricType string

const (
	MetricTypeSimple MetricType = "simple"
	MetricTypeRatio  MetricType = "ratio"
)

func (m *MetricType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := MetricType(s)
	switch dt {
	case MetricTypeSimple, MetricTypeRatio:
		*m = dt
		return nil
	default:
		return fmt.Errorf("invalid MetricType: %s", s)
	}
}
