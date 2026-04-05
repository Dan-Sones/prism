package metric

import "fmt"

type AnalysisUnit string

const (
	AnalysisUnitUser AnalysisUnit = "user"
)

func (a *AnalysisUnit) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := AnalysisUnit(s)
	switch dt {
	case AnalysisUnitUser:
		*a = dt
		return nil
	default:
		return fmt.Errorf("invalid AnalysisUnit: %s", s)
	}
}
