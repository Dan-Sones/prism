package event

type Numerator string
type Denominator string

type BinaryMetricQueryResult struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}
