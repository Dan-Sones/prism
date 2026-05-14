import math

import pandas as pd
import spotify_confidence as conf
from models.z_test import ZTestResult
from services.helpers import calculate_interval_from_alpha


def _sanitize(value: float) -> float:
    return 0.0 if math.isnan(value) or math.isinf(value) else value


def perform_z_test_for_binary_metric(df: pd.DataFrame, alpha: float, control_name: str, treatment_name: str) -> ZTestResult:
    ztest = conf.ZTest(data_frame=df,
                       numerator_column='success',
                       numerator_sum_squares_column='success',
                       denominator_column='total',
                       interval_size=calculate_interval_from_alpha(alpha),
                       correction_method='bonferroni',
                       categorical_group_columns='variation_name')

    return parse_result(ztest.difference(control_name, treatment_name))


def parse_result(df: pd.DataFrame) -> ZTestResult:
    row = df.iloc[0]
    return ZTestResult(
        absolute_difference=row['absolute_difference'],
        ci_lower=_sanitize(row['ci_lower']),
        ci_upper=_sanitize(row['ci_upper']),
        p_value=_sanitize(row['p-value']),
        adjusted_ci_lower=_sanitize(row['adjusted ci_lower']),
        adjusted_ci_upper=_sanitize(row['adjusted ci_upper']),
        adjusted_p_value=_sanitize(row['adjusted p-value']),
        is_significant=row['is_significant'],
        powered_effect=_sanitize(row['powered_effect']),
    )
