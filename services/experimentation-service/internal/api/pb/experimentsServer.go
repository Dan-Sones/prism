package pb

import (
	"context"
	pb "experimentation-service/internal/grpc/generated/experimentation_service_experiments/v1"
	"experimentation-service/internal/service"

	dbExperiment "github.com/Dan-Sones/prismdbmodels/model/experiment"
	dbMetric "github.com/Dan-Sones/prismdbmodels/model/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExperimentsServer struct {
	pb.UnimplementedExperimentationServiceExperimentsServer
	experimentService *service.ExperimentService
}

func NewExperimentsServer(experimentService *service.ExperimentService) *ExperimentsServer {
	return &ExperimentsServer{
		experimentService: experimentService,
	}
}

func (s *ExperimentsServer) GetEnrichedExperimentByKey(ctx context.Context, req *pb.GetEnrichedExperimentByKeyRequest) (*pb.GetEnrichedExperimentByKeyResponse, error) {
	enrichedExp, err := s.experimentService.GetEnrichedExperimentByKey(ctx, req.GetExperimentKey())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.GetEnrichedExperimentByKeyResponse{
		Experiment: toProtoEnrichedExperiment(enrichedExp),
	}, nil
}

func toProtoEnrichedExperiment(exp dbExperiment.EnrichedExperiment) *pb.EnrichedExperiment {
	proto := &pb.EnrichedExperiment{
		Id:            exp.ID.String(),
		Name:          exp.Name,
		CreatedAt:     timestamppb.New(exp.CreatedAt),
		FeatureFlagId: exp.FeatureFlagID,
		AaStartTime:   timestamppb.New(exp.AAStartTime),
		AaEndTime:     timestamppb.New(exp.AAEndTime),
		Hypothesis:    exp.Hypothesis,
		Description:   exp.Description,
		Status:        toProtoExperimentStatus(exp.Status),
	}

	if exp.StartTime != nil {
		proto.StartTime = timestamppb.New(*exp.StartTime)
	}
	if exp.EndTime != nil {
		proto.EndTime = timestamppb.New(*exp.EndTime)
	}

	proto.Variants = make([]*pb.ExperimentVariant, len(exp.Variants))
	for i, v := range exp.Variants {
		proto.Variants[i] = toProtoVariant(v)
	}

	proto.Metrics = make([]*pb.EnrichedMetric, len(exp.Metrics))
	for i, m := range exp.Metrics {
		proto.Metrics[i] = toProtoEnrichedMetric(m)
	}

	return proto
}

func toProtoVariant(v dbExperiment.ExperimentVariant) *pb.ExperimentVariant {
	return &pb.ExperimentVariant{
		Name:          v.Name,
		FeatureFlagId: v.FeatureFlagID,
		VariantKey:    v.VariantKey,
		UpperBound:    int32(v.UpperBound),
		LowerBound:    int32(v.LowerBound),
		VariantType:   toProtoVariantType(v.VariantType),
	}
}

func toProtoEnrichedMetric(m dbMetric.EnrichedMetric) *pb.EnrichedMetric {
	proto := &pb.EnrichedMetric{
		Id:           m.ID.String(),
		Name:         m.Name,
		MetricKey:    m.MetricKey,
		Description:  m.Description,
		CreatedAt:    timestamppb.New(m.CreatedAt),
		MetricType:   toProtoMetricType(m.MetricType),
		AnalysisUnit: toProtoAnalysisUnit(m.AnalysisUnit),
		IsBinary:     m.IsBinary,
	}

	proto.MetricComponents = make([]*pb.EnrichedMetricComponent, len(m.MetricComponents))
	for i, c := range m.MetricComponents {
		proto.MetricComponents[i] = toProtoMetricComponent(c)
	}

	return proto
}

func toProtoMetricComponent(c dbMetric.EnrichedMetricComponent) *pb.EnrichedMetricComponent {
	proto := &pb.EnrichedMetricComponent{
		Id:                   c.ID.String(),
		Role:                 toProtoComponentRole(c.Role),
		EventType:            c.EventType.EventKey,
		AggregationOperation: toProtoAggregationOperation(c.AggregationOperation),
	}

	if c.AggregationField != nil {
		fieldKey := c.AggregationField.FieldKey
		proto.AggregationField = &fieldKey
	}
	if c.SystemColumnName != nil {
		proto.SystemColumnName = c.SystemColumnName
	}

	return proto
}

func toProtoExperimentStatus(s dbExperiment.ExperimentStatus) pb.ExperimentStatus {
	switch s {
	case dbExperiment.ExperimentStatusAAPlanned:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AA_PLANNED
	case dbExperiment.ExperimentStatusAA:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AA
	case dbExperiment.ExperimentStatusAAComplete:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AA_COMPLETE
	case dbExperiment.ExperimentStatusABPlanned:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AB_PLANNED
	case dbExperiment.ExperimentStatusAB:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AB
	case dbExperiment.ExperimentStatusComplete:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_AB_COMPLETE
	default:
		return pb.ExperimentStatus_EXPERIMENT_STATUS_UNSPECIFIED
	}
}

func toProtoVariantType(v dbExperiment.VariantType) pb.VariantType {
	switch v {
	case dbExperiment.VariantTypeControl:
		return pb.VariantType_VARIANT_TYPE_CONTROL
	case dbExperiment.VariantTypeTreatment:
		return pb.VariantType_VARIANT_TYPE_TREATMENT
	default:
		return pb.VariantType_VARIANT_TYPE_UNSPECIFIED
	}
}

func toProtoMetricType(m dbMetric.MetricType) pb.MetricType {
	switch m {
	case dbMetric.MetricTypeSimple:
		return pb.MetricType_METRIC_TYPE_SIMPLE
	case dbMetric.MetricTypeRatio:
		return pb.MetricType_METRIC_TYPE_RATIO
	default:
		return pb.MetricType_METRIC_TYPE_UNSPECIFIED
	}
}

func toProtoAnalysisUnit(a dbMetric.AnalysisUnit) pb.AnalysisUnit {
	switch a {
	case dbMetric.AnalysisUnitUser:
		return pb.AnalysisUnit_ANALYSIS_UNIT_USER
	default:
		return pb.AnalysisUnit_ANALYSIS_UNIT_UNSPECIFIED
	}
}

func toProtoComponentRole(r dbMetric.ComponentRole) pb.ComponentRole {
	switch r {
	case dbMetric.ComponentRoleBaseEvent:
		return pb.ComponentRole_COMPONENT_ROLE_BASE_EVENT
	case dbMetric.ComponentRoleNumerator:
		return pb.ComponentRole_COMPONENT_ROLE_NUMERATOR
	case dbMetric.ComponentRoleDenominator:
		return pb.ComponentRole_COMPONENT_ROLE_DENOMINATOR
	default:
		return pb.ComponentRole_COMPONENT_ROLE_UNSPECIFIED
	}
}

func toProtoAggregationOperation(a dbMetric.AggregationOperation) pb.AggregationOperation {
	switch a {
	case dbMetric.AggregationOperationCount:
		return pb.AggregationOperation_AGGREGATION_OPERATION_COUNT
	case dbMetric.AggregationOperationSum:
		return pb.AggregationOperation_AGGREGATION_OPERATION_SUM
	case dbMetric.AggregationOperationAvg:
		return pb.AggregationOperation_AGGREGATION_OPERATION_AVG
	case dbMetric.AggregationOperationMin:
		return pb.AggregationOperation_AGGREGATION_OPERATION_MIN
	case dbMetric.AggregationOperationMax:
		return pb.AggregationOperation_AGGREGATION_OPERATION_MAX
	case dbMetric.AggregationOperationCountDistinct:
		return pb.AggregationOperation_AGGREGATION_OPERATION_COUNT_DISTINCT
	case dbMetric.AggregationOperationPercentile95:
		return pb.AggregationOperation_AGGREGATION_OPERATION_PERCENTILE_95
	case dbMetric.AggregationOperationPercentile99:
		return pb.AggregationOperation_AGGREGATION_OPERATION_PERCENTILE_99
	default:
		return pb.AggregationOperation_AGGREGATION_OPERATION_UNSPECIFIED
	}
}
