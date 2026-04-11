package experiment

import "fmt"

type ExperimentVariant struct {
	FeatureFlagID string      `json:"feature_flag_id"`
	VariantKey    string      `json:"variant_id"`
	UpperBound    int         `json:"upper_bound"`
	LowerBound    int         `json:"lower_bound"`
	VariantType   VariantType `json:"variant_type"`
}

type VariantType string

const (
	VariantTypeControl   VariantType = "control"
	VariantTypeTreatment VariantType = "treatment"
)

func (a *VariantType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := VariantType(s)
	switch dt {
	case VariantTypeControl, VariantTypeTreatment:
		*a = dt
		return nil
	default:
		return fmt.Errorf("invalid VariantType: %s", s)
	}
}
