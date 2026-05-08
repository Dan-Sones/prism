from models.z_test import ZTestResult
from models.decision_output import DecisionOutput, DecisionRecommendation

# TODO:
# - look into see if my use of MDE is correct
# - - book talks about practical significance boundary but I'm not sure if that concept is fully interchangable wih mde
# - - - We use MDE to determine our desired sample size earlier on in the exp lifecyle,
def make_decision(result: ZTestResult, mde: float) -> DecisionOutput:
    if not result.is_significant:
        return DecisionOutput(
            recommendation=DecisionRecommendation.NOT_RECOMMEND,
            ZTestResult=result,
        )

    if result.powered_effect > mde:
        return DecisionOutput(
            recommendation=DecisionRecommendation.INCONCLUSIVE,
            recommendation_reason="Test was underpowered, we couldn't detect the MDE",
            ZTestResult=result,
        )

    if result.ci_lower >= mde:
        return DecisionOutput(
            recommendation=DecisionRecommendation.RECOMMEND,
            recommendation_reason="The lower bound of the confidence interval is confidently above the MDE.",
            ZTestResult=result,
        )

    if result.ci_upper <= mde:
        return DecisionOutput(
            recommendation=DecisionRecommendation.NOT_RECOMMEND,
            recommendation_reason="The upper bound of the confidence interval is confidently below the MDE.",
            ZTestResult=result,
        )

    else:
        return DecisionOutput(
            recommendation=DecisionRecommendation.INCONCLUSIVE,
            recommendation_reason="Statistically significant results observed, but potentially not practically significant",
            ZTestResult=result,
        )