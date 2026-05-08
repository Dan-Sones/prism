from dataclasses import dataclass
from enum import Enum
from models.z_test import ZTestResult
import stats_engine.v1.stats_engine_pb2 as pb

class DecisionRecommendation(Enum):
    RECOMMEND = "Recommend"
    NOT_RECOMMEND = "Not Recommend"
    INCONCLUSIVE = "Inconclusive"

@dataclass
class DecisionOutput:
    recommendation: DecisionRecommendation
    recommendation_reason: str
    ZTestResult: ZTestResult

    def to_proto(self) -> pb.DecisionOutput:
        return pb.DecisionOutput(
            recommendation=pb.DecisionRecommendation.Value(f"DECISION_RECOMMENDATION_{self.recommendation.name}"),
            recommendation_reason=self.recommendation_reason,
            z_test_result=self.ZTestResult.to_proto()
        )

