package experiment

type GetAbsoluteSampleSizeRequest struct {
	AbsolutePercentageMDE float64 `json:"absolute_percentage_mde"`
	BaselineProportion    float64 `json:"baseline_proportion"`
	Alpha                 float64 `json:"alpha"`
	Power                 float64 `json:"power"`
	Treatments            int     `json:"treatments"`
}

type GetAbsoluteSampleSizeResponse struct {
	TotalSampleSize      int       `json:"total_sample_size"`
	PerVariantSampleSize []int     `json:"per_variant_sample_size"`
	Allocations          []float64 `json:"allocations"`
}
