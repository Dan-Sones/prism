from dataclasses import dataclass
from enum import Enum
from models.z_test import ZTestResult
from models.binary_observation import BinaryObservation
import stats_engine.v1.stats_engine_pb2 as pb

class DecisionRecommendation(Enum):
    RECOMMEND = "Recommend"
    NOT_RECOMMEND = "Not Recommend"
    INCONCLUSIVE = "Inconclusive"

@dataclass
class DecisionOutput:
    recommendation: str
    recommendation_reason: str
    practically_significant: bool
    statistically_significant: bool
    z_test_result: ZTestResult
    control_observation: BinaryObservation
    treatment_observation: BinaryObservation

    def to_proto(self) -> pb.DecisionOutput:
        return pb.DecisionOutput(
            recommendation=pb.DecisionRecommendation.Value(
                f"DECISION_RECOMMENDATION_{self.recommendation.name}"
            ),
            recommendation_reason=self.recommendation_reason,
            z_test_result=self.z_test_result.to_proto(),
            control_observation=self.control_observation.to_proto(),
            treatment_observation=self.treatment_observation.to_proto(),
            practically_significant = self.practically_significant,
            statistically_significant = self.statistically_significant
        )