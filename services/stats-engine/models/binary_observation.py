import stats_engine.v1.stats_engine_pb2 as pb
from dataclasses import dataclass

@dataclass
class BinaryObservation:
    numerator: int
    denominator: int

    def to_proto(self) -> pb.BinaryObservation:
        return pb.BinaryObservation(
            numerator=self.numerator,
            denominator=self.denominator
        )