import React from "react";
import type {
  ExperimentResponse,
  ExperimentResultsResponse,
} from "../../../../../api/experiments";
import MetricCard from "./MetricResultCard";
import ShipDecision from "./ShipDecision";
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

  return (
    <React.Fragment>
      <ShipDecision
        decisionRecommendation={experimentResults?.decision_recommendation}
        recommendationReason={experimentResults?.recommendation_reason}
        recommendationLoading={resultsLoading}
        recommendationError={resultsError}
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
