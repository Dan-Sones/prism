package clients

import (
	"context"
	"experimentation-service/internal/model/experiment/sampleSize"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "experimentation-service/internal/grpc/generated/stats_engine/v1"
)

type StatsEngineClient interface {
	CalculateSampleSize(ctx context.Context, enrichedMetrics []sampleSize.MetricForExperiment, alpha, power float64, variants int) (totalSampleSize int, sampleSizePerVariant []int, split []float64, error error)
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
