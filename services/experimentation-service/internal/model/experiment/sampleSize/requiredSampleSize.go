package sampleSize

type RequiredSampleSizeResponse struct {
	TotalRequiredSampleSize int            `json:"total_required_sample_size"`
	SampleSizePerVariant    map[string]int `json:"sample_size_per_variant"`
}
