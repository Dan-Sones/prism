from dataclasses import dataclass, asdict
import stats_engine.v1.stats_engine_pb2 as pb

@dataclass
class ZTestResult:
    absolute_difference: bool
    ci_lower: float
    ci_upper: float
    p_value: float
    adjusted_ci_lower: float
    adjusted_ci_upper: float
    adjusted_p_value: float
    is_significant: bool
    powered_effect: float

    def to_proto(self) -> pb.ZTestResult:
        return pb.ZTestResult(**asdict(self))