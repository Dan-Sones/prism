package metric

import "fmt"

type AggregationOperation string

const (
	AggregationOperationCount         AggregationOperation = "COUNT"
	AggregationOperationSum           AggregationOperation = "SUM"
	AggregationOperationAvg           AggregationOperation = "AVG"
	AggregationOperationMin           AggregationOperation = "MIN"
	AggregationOperationMax           AggregationOperation = "MAX"
	AggregationOperationCountDistinct AggregationOperation = "COUNT_DISTINCT"
	AggregationOperationPercentile95  AggregationOperation = "PERCENTILE_95"
	AggregationOperationPercentile99  AggregationOperation = "PERCENTILE_99"
)

func (a *AggregationOperation) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := AggregationOperation(s)
	switch dt {
	case AggregationOperationCount, AggregationOperationSum, AggregationOperationAvg, AggregationOperationMin, AggregationOperationMax, AggregationOperationCountDistinct, AggregationOperationPercentile95, AggregationOperationPercentile99:
		*a = dt
		return nil
	default:
		return fmt.Errorf("invalid AggregationOperation: %s", s)
	}
}
