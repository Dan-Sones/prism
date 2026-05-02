package clients

import (
	"context"
	pb "data-cooking-service/internal/grpc/generated/experimentation_service_experiments/v1"
	"time"

	dbEvent "github.com/Dan-Sones/prismdbmodels/model/event"
	dbExperiment "github.com/Dan-Sones/prismdbmodels/model/experiment"
	dbMetric "github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExperimentationExperimentClient interface {
	GetEnrichedExperimentByKey(ctx context.Context, experimentKey string) (dbExperiment.EnrichedExperiment, error)
	Close() error
}

type GrpcExperimentationExperimentClient struct {
	conn   *grpc.ClientConn
	client pb.ExperimentationServiceExperimentsClient
}

func NewGrpcExperimentationExperimentClient(experimentationServiceAddr string) (*GrpcExperimentationExperimentClient, error) {
	conn, err := grpc.NewClient(experimentationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcExperimentationExperimentClient{
		conn:   conn,
		client: pb.NewExperimentationServiceExperimentsClient(conn),
	}, nil
}

func (c *GrpcExperimentationExperimentClient) GetEnrichedExperimentByKey(ctx context.Context, experimentKey string) (dbExperiment.EnrichedExperiment, error) {
	resp, err := c.client.GetEnrichedExperimentByKey(ctx, &pb.GetEnrichedExperimentByKeyRequest{
		ExperimentKey: experimentKey,
	})
	if err != nil {
		return dbExperiment.EnrichedExperiment{}, err
	}

	return fromProtoEnrichedExperiment(resp.Experiment), nil
}

func (c *GrpcExperimentationExperimentClient) Close() error {
	return c.conn.Close()
}

func fromProtoEnrichedExperiment(p *pb.EnrichedExperiment) dbExperiment.EnrichedExperiment {
	id, _ := uuid.Parse(p.Id)

	exp := dbExperiment.EnrichedExperiment{
		ID:            id,
		Name:          p.Name,
		CreatedAt:     p.CreatedAt.AsTime(),
		FeatureFlagID: p.FeatureFlagId,
		AAStartTime:   p.AaStartTime.AsTime(),
		AAEndTime:     p.AaEndTime.AsTime(),
		Hypothesis:    p.Hypothesis,
		Description:   p.Description,
		Status:        fromProtoExperimentStatus(p.Status),
	}

	if p.StartTime != nil {
		t := p.StartTime.AsTime()
		exp.StartTime = &t
	}
	if p.EndTime != nil {
		t := p.EndTime.AsTime()
		exp.EndTime = &t
	}

	exp.Variants = make([]dbExperiment.ExperimentVariant, len(p.Variants))
	for i, v := range p.Variants {
		exp.Variants[i] = fromProtoVariant(v)
	}

	exp.Metrics = make([]dbMetric.EnrichedMetric, len(p.Metrics))
	for i, m := range p.Metrics {
		exp.Metrics[i] = fromProtoEnrichedMetric(m)
	}

	return exp
}

func fromProtoVariant(v *pb.ExperimentVariant) dbExperiment.ExperimentVariant {
	return dbExperiment.ExperimentVariant{
		Name:          v.Name,
		FeatureFlagID: v.FeatureFlagId,
		VariantKey:    v.VariantKey,
		UpperBound:    int(v.UpperBound),
		LowerBound:    int(v.LowerBound),
		VariantType:   fromProtoVariantType(v.VariantType),
	}
}

func fromProtoEnrichedMetric(m *pb.EnrichedMetric) dbMetric.EnrichedMetric {
	id, _ := uuid.Parse(m.Id)

	components := make([]dbMetric.EnrichedMetricComponent, len(m.MetricComponents))
	for i, c := range m.MetricComponents {
		components[i] = fromProtoMetricComponent(c)
	}

	return dbMetric.EnrichedMetric{
		ID:               id,
		Name:             m.Name,
		MetricKey:        m.MetricKey,
		Description:      m.Description,
		CreatedAt:        m.CreatedAt.AsTime(),
		MetricType:       fromProtoMetricType(m.MetricType),
		AnalysisUnit:     fromProtoAnalysisUnit(m.AnalysisUnit),
		IsBinary:         m.IsBinary,
		MetricComponents: components,
	}
}

func fromProtoMetricComponent(c *pb.EnrichedMetricComponent) dbMetric.EnrichedMetricComponent {
	id, _ := uuid.Parse(c.Id)

	comp := dbMetric.EnrichedMetricComponent{
		ID:                   id,
		Role:                 fromProtoComponentRole(c.Role),
		AggregationOperation: fromProtoAggregationOperation(c.AggregationOperation),
		EventType:            fromProtoEventType(c.EventType),
	}

	if c.SystemColumnName != nil {
		comp.SystemColumnName = c.SystemColumnName
	}

	return comp
}

func fromProtoEventType(eventKey string) dbEvent.EventType {
	return dbEvent.EventType{
		EventKey: eventKey,
	}
}

func fromProtoExperimentStatus(s pb.ExperimentStatus) dbExperiment.ExperimentStatus {
	switch s {
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AA_PLANNED:
		return dbExperiment.ExperimentStatusAAPlanned
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AA:
		return dbExperiment.ExperimentStatusAA
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AA_COMPLETE:
		return dbExperiment.ExperimentStatusAAComplete
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AB_PLANNED:
		return dbExperiment.ExperimentStatusABPlanned
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AB:
		return dbExperiment.ExperimentStatusAB
	case pb.ExperimentStatus_EXPERIMENT_STATUS_AB_COMPLETE:
		return dbExperiment.ExperimentStatusComplete
	default:
		return ""
	}
}

func fromProtoVariantType(v pb.VariantType) dbExperiment.VariantType {
	switch v {
	case pb.VariantType_VARIANT_TYPE_CONTROL:
		return dbExperiment.VariantTypeControl
	case pb.VariantType_VARIANT_TYPE_TREATMENT:
		return dbExperiment.VariantTypeTreatment
	default:
		return ""
	}
}

func fromProtoMetricType(m pb.MetricType) dbMetric.MetricType {
	switch m {
	case pb.MetricType_METRIC_TYPE_SIMPLE:
		return dbMetric.MetricTypeSimple
	case pb.MetricType_METRIC_TYPE_RATIO:
		return dbMetric.MetricTypeRatio
	default:
		return ""
	}
}

func fromProtoAnalysisUnit(a pb.AnalysisUnit) dbMetric.AnalysisUnit {
	switch a {
	case pb.AnalysisUnit_ANALYSIS_UNIT_USER:
		return dbMetric.AnalysisUnitUser
	default:
		return ""
	}
}

func fromProtoComponentRole(r pb.ComponentRole) dbMetric.ComponentRole {
	switch r {
	case pb.ComponentRole_COMPONENT_ROLE_BASE_EVENT:
		return dbMetric.ComponentRoleBaseEvent
	case pb.ComponentRole_COMPONENT_ROLE_NUMERATOR:
		return dbMetric.ComponentRoleNumerator
	case pb.ComponentRole_COMPONENT_ROLE_DENOMINATOR:
		return dbMetric.ComponentRoleDenominator
	default:
		return ""
	}
}

func fromProtoAggregationOperation(a pb.AggregationOperation) dbMetric.AggregationOperation {
	switch a {
	case pb.AggregationOperation_AGGREGATION_OPERATION_COUNT:
		return dbMetric.AggregationOperationCount
	case pb.AggregationOperation_AGGREGATION_OPERATION_SUM:
		return dbMetric.AggregationOperationSum
	case pb.AggregationOperation_AGGREGATION_OPERATION_AVG:
		return dbMetric.AggregationOperationAvg
	case pb.AggregationOperation_AGGREGATION_OPERATION_MIN:
		return dbMetric.AggregationOperationMin
	case pb.AggregationOperation_AGGREGATION_OPERATION_MAX:
		return dbMetric.AggregationOperationMax
	case pb.AggregationOperation_AGGREGATION_OPERATION_COUNT_DISTINCT:
		return dbMetric.AggregationOperationCountDistinct
	case pb.AggregationOperation_AGGREGATION_OPERATION_PERCENTILE_95:
		return dbMetric.AggregationOperationPercentile95
	case pb.AggregationOperation_AGGREGATION_OPERATION_PERCENTILE_99:
		return dbMetric.AggregationOperationPercentile99
	default:
		return ""
	}
}

func fromProtoTime(t interface{ AsTime() time.Time }) time.Time {
	return t.AsTime()
}
