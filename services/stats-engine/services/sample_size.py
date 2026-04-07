import spotify_confidence

def get_absolute_sample_size(absolute_percentage_mde,baseline_proportion, alpha, power,treatments):
    total, per_variant, allocations = spotify_confidence.SampleSize.binomial(
        absolute_percentage_mde=absolute_percentage_mde,
        baseline_proportion=baseline_proportion,
        alpha=alpha,
        power=power,
        treatments=treatments,
    )
    return total, per_variant, allocations
