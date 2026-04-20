from services.sample_size import get_sample_size_for_binomial_metric
import stats_engine.v1.stats_engine_pb2_grpc as ses_grpc
import stats_engine.v1.stats_engine_pb2 as ses_pb2


class StatsEngineServer(ses_grpc.StatsEngineServicer):
    def CalculateSampleSizeForBinomialMetric(self, request, context):
        total, per_variant, split = get_sample_size_for_binomial_metric(
            request.experiment_exposures,
            request.experiment_conversions,
            request.absolute_percentage_mde,
            request.variant_count,
        )
        return ses_pb2.CalculateSampleSizeForBinomialMetricResponse(  # type: ignore
            total=total,
            per_variant=per_variant,
            split=split,
        )


