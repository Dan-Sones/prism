from models.z_test import ZTestResult
from models.decision_output import DecisionOutput, DecisionRecommendation
from models.binary_observation import BinaryObservation

def make_decision_for_z_test(result: ZTestResult, mde: float, control_observation: BinaryObservation, treatment_observation: BinaryObservation) -> DecisionOutput:
    # not statistically or practically significant
    if not result.is_significant and result.ci_upper < mde:
        return DecisionOutput(recommendation=DecisionRecommendation.NOT_RECOMMEND,
                              recommendation_reason=
                              "Results are not statistically significant nor practically significant",
                              z_test_result=result,
                              control_observation=control_observation,
                              treatment_observation=treatment_observation,
                              practically_significant=False,
                              statistically_significant=False,
                              )

    # statistically and practically significant - clear winner
    if result.is_significant and result.ci_lower >= mde:
        return DecisionOutput(recommendation=DecisionRecommendation.RECOMMEND,
                              recommendation_reason=
                              "Results are statistically and practically significant",
                              z_test_result=result,
                              control_observation=control_observation,
                              treatment_observation=treatment_observation,
                              statistically_significant=True,
                              practically_significant=True,
                              )

    # Statistically significant, but not strong enough results
    if result.is_significant and result.ci_upper >= mde:
        return DecisionOutput(recommendation=DecisionRecommendation.INCONCLUSIVE,
                              recommendation_reason=
                              "Results are statistically significant, but only JUST practically significant, you may want to re-run the experiment a large sample size",
                              z_test_result=result,
                              control_observation=control_observation,
                              treatment_observation=treatment_observation,
                              practically_significant=False,
                              statistically_significant=True
                              )

    # statistically significant but not practically significant
    if result.is_significant and result.ci_upper < mde:
        return DecisionOutput(recommendation=DecisionRecommendation.NOT_RECOMMEND,
                              recommendation_reason=
                              "Results are statistically significant, but not practically significant",
                              z_test_result=result,
                              control_observation=control_observation,
                              treatment_observation=treatment_observation,
                              practically_significant=False,
                              statistically_significant=True
                              )

    # Inconclusive results
    return DecisionOutput(recommendation=DecisionRecommendation.INCONCLUSIVE,
                          recommendation_reason=
                          "Results are inconclusive",
                          z_test_result=result,
                          control_observation=control_observation,
                          treatment_observation=treatment_observation,
                          practically_significant=False,
                          statistically_significant=False,
                          )