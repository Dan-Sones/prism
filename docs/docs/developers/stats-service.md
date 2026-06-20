---
sidebar_position: 8
---

# Stats Service

The stats service is a Python GRPC service that performs the statistical operations for Prism.

It currently exposes a sample size calculator for binomial metrics and a z-test implementation for binary metrics.

The service is designed to work exclusively with pre-aggregated results. This means you should not attempt to send an experiment's results to the service. You should instead send the representation of the metric itself with numerator and denominator values as dictated in the proto file.

The service leverages the [Spotify Confidence Library](https://github.com/spotify/confidence).

## Z Test Decision Rule

Based on the results of the z-test, the following decision rule is implemented On the basis of the outputted confidence intervals from the libraries Z-Test function.

| Outcome           | Statistical Significance | Practical Significance | Recommendation       | Action                                                                                                                                     |
| ----------------- | ------------------------ | ---------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| Clear Winner      | ✅                       | ✅                     | **Recommend**        | Ship the treatment variant. It met your MDE and was statistically significant                                                              |
| No Effect         | ❌                       | ❌                     | **Do Not Recommend** | Do not ship. Results were not Practically OR significantly significant. Rethink your hypothesis and maybe re-run with a larger sample size |
| Negligible Effect | ✅                       | ❌                     | **Do Not Recommend** | The effect is real but does not meet you MDE so can't be recommended                                                                       |
| Regression        | ✅                       | ❌                     | **Do Not Recommend** | Treatment caused a statistically significant regression vs control. Do not ship.                                                           |
| Borderline        | ✅                       | ⚠️                     | **Inconclusive**     | Re-run the experiment with a larger sample size                                                                                            |
