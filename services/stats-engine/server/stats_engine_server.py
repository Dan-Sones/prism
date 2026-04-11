from services.sample_size import get_absolute_sample_size
import stats_engine.v1.stats_engine_pb2 as ses_pb2
import stats_engine.v1.stats_engine_pb2_grpc as ses_grpc


class StatsEngineServer(ses_grpc.StatsEngineServicer):
    def CalculateSampleSizeAbsoluteMetric(self, request, context):
        total, per_variant, split = get_absolute_sample_size(request.sample_size)
        return ses_pb2.CalculateSampleSizeAbsoluteMetricResponse(
            total=total,
            per_variant=per_variant,
            split=split,
        )

