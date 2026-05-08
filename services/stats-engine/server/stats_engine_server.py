from services.sample_size import get_sample_size
import stats_engine.v1.stats_engine_pb2_grpc as ses_grpc
import stats_engine.v1.stats_engine_pb2 as ses_pb
import pandas as pd
from services.z_test import perform_z_test_for_binary_metric
from services.decision import make_decision

DIRECTION_MAP = {
      0: "increase",
      1: "decrease",
  }

def calculate_relative_mde(mde: float, baseline: float) -> float:
    return mde / baseline

def calculate_variance(metric: ses_pb.MetricDetails) -> float:
    if metric.is_binary:
        return metric.baseline * (1 - metric.baseline)
    else:
        raise NotImplementedError("Variance calculation for non-binary metrics not yet implemented.")

class StatsEngineServer(ses_grpc.StatsEngineServicer):
    def CalculateSampleSize(self, request, context):
        df = pd.DataFrame([], columns=["metric_name", "baseline", "variance", "is_binary", "mde", "direction", "nim"])

        for metric in request.metrics:
            df.loc[len(df)] = {
                "metric_name": metric.metric_key,
                "baseline": metric.baseline,
                "variance": calculate_variance(metric),
                "is_binary": metric.is_binary,
                "mde": calculate_relative_mde(metric.absolute_percentage_mde, metric.baseline),
                "direction": DIRECTION_MAP.get(metric.direction),
                "nim": None,  # Just looking at success metrics atm so no need for this
            }

        total_sample_size = get_sample_size(df, request.power, request.alpha)
        per_variant_sample_size = total_sample_size / 2 # Only control and treatment atm, so divide by 2

        res = {
            "total_sample_size": total_sample_size.astype(int),
            "sample_size_per_variant": [per_variant_sample_size.astype(int), per_variant_sample_size.astype(int)],
            "split": [0.5, 0.5], # Only control and treatment atm
        }

        return ses_pb.CalculateSampleSizeResponse(**res)

    def PerformZTestBinaryMetric(self, request, context):
        df = pd.DataFrame([], columns=["variation_name", "success", "total"])

        df.loc[0]= {
            "variation_name": request.control_name,
            "success": request.numerator_control,
            "total": request.denominator_control,
        }

        df.loc[1]= {
            "variation_name": request.treatment_name,
            "success": request.numerator_treatment,
            "total": request.denominator_treatment,
        }

        z_test_result = perform_z_test_for_binary_metric(df, request.alpha)
        decision = make_decision(z_test_result, request.absolute_percentage_mde)

        return decision.to_proto()





