from models.z_test import ZTestResult
from models.decision_output import DecisionOutput, DecisionRecommendation

def make_decision_for_z_test(result: ZTestResult, mde: float) -> DecisionOutput:
    # not statistically or practically significant
    if not result.is_significant and result.ci_upper < mde:
        return DecisionOutput(recommendation=DecisionRecommendation.NOT_RECOMMEND,
                              recommendation_reason=
                              "Results are not statistically significant nor practically significant",
                              ZTestResult=result)

    # statistically and practically significant - clear winner
    if result.is_significant and result.ci_lower >= mde:
        return DecisionOutput(recommendation=DecisionRecommendation.RECOMMEND,
                              recommendation_reason=
                              "Results are statistically and practically significant",
                              ZTestResult=result)

    # Statistically significant, but not strong enough results
    if result.is_significant and result.ci_upper >= mde:
        return DecisionOutput(recommendation=DecisionRecommendation.INCONCLUSIVE,
                              recommendation_reason=
                              "Results are statistically significant, but only JUST practically significant, you may want to re-run the experiment a large sample size",
                              ZTestResult=result
                              )


    # statistically significant but not practically significant
    if result.is_significant and result.ci_upper < mde:
        return DecisionOutput(recommendation=DecisionRecommendation.NOT_RECOMMEND,
                              recommendation_reason=
                              "Results are statistically significant, but not practically significant",
                              ZTestResult=result)



    # Inconclusive results
    return DecisionOutput(recommendation=DecisionRecommendation.INCONCLUSIVE,
                          recommendation_reason=
                          "Results are inconclusive",
                          ZTestResult=result)