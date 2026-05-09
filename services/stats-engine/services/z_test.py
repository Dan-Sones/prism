import pandas as pd
import spotify_confidence as conf
from models.z_test import ZTestResult
from services.helpers import calculate_interval_from_alpha


def perform_z_test_for_binary_metric(df: pd.DataFrame, alpha: float) -> ZTestResult:
    ztest = conf.ZTest(data_frame=df,
                       numerator_column='success',
                       numerator_sum_squares_column='success',
                       denominator_column='total',
                       interval_size=calculate_interval_from_alpha(alpha),
                       correction_method='bonferroni',
                       categorical_group_columns='variation_name')

    return parse_result(ztest.difference('control', 'treatment'))


def parse_result(df: pd.DataFrame) -> ZTestResult:
    row = df.iloc[0]
    return ZTestResult(
        absolute_difference=row['absolute_difference'],
        ci_lower=row['ci_lower'],
        ci_upper=row['ci_upper'],
        p_value=row['p-value'],
        adjusted_ci_lower=row['adjusted ci_lower'],
        adjusted_ci_upper=row['adjusted ci_upper'],
        adjusted_p_value=row['adjusted p-value'],
        is_significant=row['is_significant'],
        powered_effect=row['powered_effect'],
    )
