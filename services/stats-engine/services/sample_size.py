import spotify_confidence

def get_sample_size_for_binomial_metric(experiment_exposures: int, experiment_conversions: int, absolute_percentage_mde: float, treatments: int):
    baseline_proportion = calculate_baseline_proportion(experiment_exposures, experiment_conversions)
    total, per_variant, allocations = spotify_confidence.SampleSize.binomial(
        absolute_percentage_mde=absolute_percentage_mde,
        baseline_proportion=baseline_proportion,
        alpha=0.05,
        power=0.80,
        treatments=treatments,
    )
    return total, per_variant, allocations

def calculate_baseline_proportion(experiment_exposures: int, experiment_conversions: int):
    if experiment_exposures == 0:
        return 0.0
    return experiment_conversions / experiment_exposures