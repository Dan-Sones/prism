import type { DecisionRecommendation } from "../../../../../api/experiments";
import Card from "../../../../../components/card/Card";
import ExclamationTriangleIcon from "../../../../../components/icons/ExclamationTriangleIcon";
import LoadingPlaceholder from "../../../../../components/spinner/LoadingPlaceholder";
import CheckCircleIcon from "../../../../../components/icons/CheckCircleIcon";

interface ShipDecisionProps {
  decisionRecommendation?: DecisionRecommendation;
  recommendationReason?: string;
  recommendationLoading?: boolean;
  recommendationError?: boolean;
}

interface ShipDecisionStyles {
  fontColor: string;
  bgColor: string;
  borderColor: string;
}

const ShipDecision = (props: ShipDecisionProps) => {
  const {
    decisionRecommendation,
    recommendationReason,
    recommendationLoading,
    recommendationError,
  } = props;
  const dcrToString = (dcr: DecisionRecommendation) => {
    switch (dcr) {
      case "DECISION_RECOMMENDATION_RECOMMEND":
        return "Shipping Recommended";
      case "DECISION_RECOMMENDATION_NOT_RECOMMEND":
        return "Shipping Not Recommended";
      case "DECISION_RECOMMENDATION_INCONCLUSIVE":
        return "Inconclusive";
      default:
        return "No Recommendation";
    }
  };

  const getStyles = (dcr: DecisionRecommendation): ShipDecisionStyles => {
    switch (dcr) {
      case "DECISION_RECOMMENDATION_RECOMMEND":
        return {
          fontColor: "text-green-700",
          bgColor: "!bg-green-100",
          borderColor: "border-green-300",
        };
      case "DECISION_RECOMMENDATION_NOT_RECOMMEND":
        return {
          fontColor: "text-red-700",
          bgColor: "!bg-red-100",
          borderColor: "border-red-300",
        };
      case "DECISION_RECOMMENDATION_INCONCLUSIVE":
        return {
          fontColor: "text-yellow-700",
          bgColor: "!bg-yellow-100",
          borderColor: "border-yellow-300",
        };
      default:
        return {
          fontColor: "text-gray-700",
          bgColor: "!bg-gray-100",
          borderColor: "border-gray-300",
        };
    }
  };

  const styles = getStyles(
    decisionRecommendation ?? "DECISION_RECOMMENDATION_UNSPECIFIED",
  );

  const getIcon = (dcr: DecisionRecommendation) => {
    switch (dcr) {
      case "DECISION_RECOMMENDATION_RECOMMEND":
        return (
          <CheckCircleIcon
            className={`mr-2 inline h-5 w-5 ${styles.fontColor}`}
          />
        );
      case "DECISION_RECOMMENDATION_NOT_RECOMMEND":
        return (
          <ExclamationTriangleIcon
            className={`mr-2 inline h-5 w-5 ${styles.fontColor}`}
          />
        );
      case "DECISION_RECOMMENDATION_INCONCLUSIVE":
        return (
          <ExclamationTriangleIcon
            className={`mr-2 inline h-5 w-5 ${styles.fontColor}`}
          />
        );
      default:
        return null;
    }
  };

  return (
    <Card
      className={`flex flex-row !gap-1 border-1 ${styles.borderColor} ${styles.bgColor} !p-4`}
    >
      {recommendationLoading ||
      recommendationError ||
      !decisionRecommendation ||
      !recommendationReason ? (
        <LoadingPlaceholder />
      ) : (
        <>
          <div className="flex items-center justify-center">
            {getIcon(decisionRecommendation)}
          </div>
          <div className="flex flex-col gap-1">
            <p
              className={`flex items-center gap-2 text-sm font-medium ${styles.fontColor}`}
            >
              {dcrToString(decisionRecommendation)}
            </p>
            <p className="text-muted-foreground text-xs">
              {recommendationReason}
            </p>
          </div>
        </>
      )}
    </Card>
  );
};

export default ShipDecision;
