package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "experimentation-service/internal/grpc/generated/stats_engine/v1"
)

type StatsEngineClient interface {
	CalculateSampleSizeForBinomialMetric(ctx context.Context, absolutePercentageMDE float64, experimentConversions, experimentExposures, treatments int) (total int, per_variant []int, split []float64, error error)
	Close() error
}

type GrpcStatsEngineClient struct {
	connn  *grpc.ClientConn
	client pb.StatsEngineClient
}

func (g *GrpcStatsEngineClient) CalculateSampleSizeForBinomialMetric(ctx context.Context, absolutePercentageMDE float64, experimentConversions, experimentExposures, treatments int) (total int, per_variant []int, split []float64, error error) {
	resp, err := g.client.CalculateSampleSizeForBinomialMetric(ctx, &pb.CalculateSampleSizeForBinomialMetricRequest{
		AbsolutePercentageMde: absolutePercentageMDE,
		ExperimentConversions: int64(experimentConversions),
		ExperimentExposures:   int64(experimentExposures),
		VariantCount:          int32(treatments),
	})
	if err != nil {
		return 0, nil, nil, err
	}

	perVariant := make([]int, len(resp.PerVariant))
	for i, v := range resp.PerVariant {
		perVariant[i] = int(v)
	}

	return int(resp.Total), perVariant, resp.Split, nil
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
