package clients

import (
	"context"
	"experimentation-service/internal/model/experiment/sampleSize"

	"github.com/Dan-Sones/prismdbmodels/model/experimentResults"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "experimentation-service/internal/grpc/generated/stats_engine/v1"
)

type StatsEngineClient interface {
	CalculateSampleSize(ctx context.Context, enrichedMetrics []sampleSize.MetricForExperiment, alpha, power float64, variants int) (totalSampleSize int, sampleSizePerVariant []int, split []float64, error error)
	PerformZTestBinaryMetric(ctx context.Context, controlName,
		treatmentName string,
		controlSuccesses,
		controlTrials,
		treatmentSuccesses,
		treatmentTrials int64,
		absolutePercentageMde float64,
	) (recommendation experimentResults.DecisionRecommendation,
		reason string,
		ztestResult *experimentResults.ZTestResult,
		practicallySignificant bool,
		statisticallySignificant bool,
		err error)
	Close() error
}

type GrpcStatsEngineClient struct {
	connn  *grpc.ClientConn
	client pb.StatsEngineClient
}

func (g *GrpcStatsEngineClient) CalculateSampleSize(ctx context.Context, enrichedMetrics []sampleSize.MetricForExperiment, alpha, power float64, variants int) (totalSampleSize int, sampleSizePerVariant []int, split []float64, error error) {
	grpcFormattedMetrics, err := ConvertToGrpcMetrics(enrichedMetrics)
	if err != nil {
		return 0, nil, nil, err
	}

	resp, err := g.client.CalculateSampleSize(ctx, &pb.CalculateSampleSizeRequest{
		Metrics: grpcFormattedMetrics,
		Alpha:   alpha,
		Power:   power,
	})
	if err != nil {
		return 0, nil, nil, err
	}

	perVariant := make([]int, len(resp.SampleSizePerVariant))
	for i, v := range resp.SampleSizePerVariant {
		perVariant[i] = int(v)
	}

	return int(resp.TotalSampleSize), perVariant, resp.Split, nil
}

func (g *GrpcStatsEngineClient) Close() error {
	return g.connn.Close()
}

func NewStatsEngineClient(statsEngineAddr string) (*GrpcStatsEngineClient, error) {
	conn, err := grpc.NewClient(statsEngineAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcStatsEngineClient{
		connn:  conn,
		client: pb.NewStatsEngineClient(conn),
	}, nil
}

func ConvertToGrpcMetrics(enrichedMetrics []sampleSize.MetricForExperiment) ([]*pb.MetricDetails, error) {
	grpcMetrics := make([]*pb.MetricDetails, len(enrichedMetrics))
	for i, m := range enrichedMetrics {
		formattedDirection, err := m.DirectionGrpcFormat()
		if err != nil {
			return nil, err
		}

		grpcMetrics[i] = &pb.MetricDetails{
			MetricKey:             m.MetricKey(),
			Baseline:              m.MetricValue(),
			AbsolutePercentageMde: m.MDE(),
			IsBinary:              m.IsBinary(),
			Direction:             formattedDirection,
		}
	}

	return grpcMetrics, nil
}

func (g *GrpcStatsEngineClient) PerformZTestBinaryMetric(ctx context.Context,
	controlName,
	treatmentName string,
	controlSuccesses,
	controlTrials,
	treatmentSuccesses,
	treatmentTrials int64,
	absolutePercentageMde float64,
) (reccommendation experimentResults.DecisionRecommendation,
	reason string,
	zTestResult *experimentResults.ZTestResult,
	practicallySignificant bool,
	statisticallySignificant bool,
	err error) {

	resp, err := g.client.PerformZTestBinaryMetric(ctx, &pb.PerformZTestBinaryMetricRequest{
		ControlName:   controlName,
		TreatmentName: treatmentName,
		ControlObservation: &pb.BinaryObservation{
			Numerator:   int32(controlSuccesses),
			Denominator: int32(controlTrials),
		},
		TreatmentObservation: &pb.BinaryObservation{
			Numerator:   int32(treatmentSuccesses),
			Denominator: int32(treatmentTrials),
		},
		AbsolutePercentageMde: absolutePercentageMde,
		Alpha:                 0.05,
	})
	if err != nil {
		return experimentResults.DecisionRecommendationUnspecified, "", nil, false, false, err
	}

	recommendationMap := map[pb.DecisionRecommendation]experimentResults.DecisionRecommendation{
		pb.DecisionRecommendation_DECISION_RECOMMENDATION_UNSPECIFIED:   experimentResults.DecisionRecommendationUnspecified,
		pb.DecisionRecommendation_DECISION_RECOMMENDATION_RECOMMEND:     experimentResults.DecisionRecommendationRecommend,
		pb.DecisionRecommendation_DECISION_RECOMMENDATION_NOT_RECOMMEND: experimentResults.DecisionRecommendationNotRecommend,
		pb.DecisionRecommendation_DECISION_RECOMMENDATION_INCONCLUSIVE:  experimentResults.DecisionRecommendationInconclusive,
	}

	return recommendationMap[resp.GetRecommendation()], resp.GetRecommendationReason(), TransformZTestResult(resp.GetZTestResult()), resp.GetPracticallySignificant(), resp.GetStatisticallySignificant(), nil
}

func TransformZTestResult(grpcResult *pb.ZTestResult) *experimentResults.ZTestResult {
	if grpcResult == nil {
		return nil
	}

	return &experimentResults.ZTestResult{
		AbsoluteDifference: grpcResult.GetAbsoluteDifference(),
		CILower:            grpcResult.GetCiLower(),
		CIUpper:            grpcResult.GetCiUpper(),
		PValue:             grpcResult.GetPValue(),
		AdjustedCILower:    grpcResult.GetAdjustedCiLower(),
		AdjustedCIUpper:    grpcResult.GetAdjustedCiUpper(),
		AdjustedPValue:     grpcResult.GetAdjustedPValue(),
		IsSignificant:      grpcResult.GetIsSignificant(),
		PoweredEffect:      grpcResult.GetPoweredEffect(),
	}
}
