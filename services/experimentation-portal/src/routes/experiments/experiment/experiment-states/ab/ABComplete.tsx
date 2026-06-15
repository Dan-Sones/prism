import React from "react";
import type {
  ExperimentResponse,
  ExperimentResultsResponse,
} from "../../../../../api/experiments";
import MetricCard from "./MetricResultCard";
import ShipDecision from "./ShipDecision";
import SampleSizeAchieved from "./SampleSizeAchieved";
import LoadingPlaceholder from "../../../../../components/spinner/LoadingPlaceholder";

interface ABCompleteProps {
  experimentDetails?: ExperimentResponse;
  experimentResults?: ExperimentResultsResponse;
  resultsLoading?: boolean;
  resultsError?: boolean;
}

const ABComplete = (props: ABCompleteProps) => {
  const { experimentResults, resultsLoading, resultsError, experimentDetails } =
    props;

  const firstMetricObservations = React.useMemo(() => {
    if (!experimentResults?.metric_observations) return undefined;
    const firstMetricId = Object.keys(experimentResults.metric_observations)[0];
    return firstMetricId
      ? experimentResults.metric_observations[firstMetricId]
      : undefined;
  }, [experimentResults]);

  const controlKey =
    experimentDetails?.variants.find((v) => v.type === "control")?.key ?? "";
  const treatmentKey =
    experimentDetails?.variants.find((v) => v.type === "treatment")?.key ?? "";
  const requiredPerVariant =
    experimentDetails?.total_required_sample_size != null
      ? Math.floor(experimentDetails.total_required_sample_size / 2)
      : undefined;

  return (
    <React.Fragment>
      <ShipDecision
        decisionRecommendation={experimentResults?.decision_recommendation}
        recommendationReason={experimentResults?.recommendation_reason}
        recommendationLoading={resultsLoading}
        recommendationError={resultsError}
      />
      <SampleSizeAchieved
        requiredPerVariant={requiredPerVariant}
        controlActual={firstMetricObservations?.control?.denominator ?? 0}
        controlKey={controlKey}
        treatmentActual={firstMetricObservations?.treatment?.denominator ?? 0}
        treatmentKey={treatmentKey}
      />
      {!experimentResults || resultsLoading || !experimentDetails ? (
        <LoadingPlaceholder />
      ) : (
        Object.keys(experimentResults.metrics).map((metricId) => (
          <MetricCard
            key={metricId}
            metricDetails={experimentResults.metrics[metricId]}
            metricObservations={experimentResults.metric_observations[metricId]}
            zTestResult={experimentResults.test_results[metricId]}
            practicallySignificant={experimentResults.practically_significant}
            statisticallySignificant={
              experimentResults.statistically_significant
            }
            resultsLoading={resultsLoading}
            resultsError={resultsError}
            variants={experimentDetails?.variants}
          />
        ))
      )}
    </React.Fragment>
  );
};

export default ABComplete;
