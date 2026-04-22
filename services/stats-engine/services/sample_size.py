import spotify_confidence
import pandas as pd
from spotify_confidence import SampleSizeCalculator

def get_sample_size_for_binomial_metric(baseline: float, absolute_percentage_mde: float, treatments: int, power: float, alpha: float):
    total, per_variant, allocations = spotify_confidence.SampleSize.binomial(
        absolute_percentage_mde=absolute_percentage_mde,
        baseline_proportion=baseline,
        alpha=alpha,
        power=power,
        treatments=treatments,
    )
    return total, per_variant, allocations

def get_sample_size(df: pd.DataFrame, power: float, alpha: float):
    calculator = SampleSizeCalculator(
        data_frame=df,
        point_estimate_column="baseline",
        var_column="variance",
        is_binary_column="is_binary",
        metric_column="metric_name",
        correction_method="bonferroni",
        power=power,
        interval_size=calculate_interval_from_alpha(alpha)
    )

    result = calculator.sample_size(
        treatment_weights=[1, 1], # 50/50 Split only for now!!
        mde_column="mde",
        nim_column="nim",
        preferred_direction_column="direction",
    )

    print(result.to_string())

    return result


def calculate_interval_from_alpha(alpha: float) -> float:
    return 1 - alpha