import { useCallback } from "react";
import type {
  EnrichedExperimentMetric,
  ExperimentVariantResponse,
  MetricValue,
  ZTestResult,
} from "../../../../../api/experiments";
import Card from "../../../../../components/card/Card";
import CheckCircleIcon from "../../../../../components/icons/CheckCircleIcon";
import MetricIndicator from "./MetricIndicator";
import XCircleIcon from "../../../../../components/icons/XCircleIcon";

interface MetricResultCardProps {
  practicallySignificant?: boolean;
  statisticallySignificant?: boolean;
  metricDetails: EnrichedExperimentMetric;
  metricObservations: Record<string, MetricValue>;
  zTestResult: ZTestResult;
  resultsLoading?: boolean;
  resultsError?: boolean;
  variants?: Array<ExperimentVariantResponse>;
}

const MetricResultCard = (props: MetricResultCardProps) => {
  // TODO: Loading and error states
  // TODO: Include variant names NOT just keys in expDetails or even expResults?
  //

  const {
    metricDetails,
    metricObservations,
    zTestResult,
    statisticallySignificant,
    practicallySignificant,
    variants,
  } = props;

  const controlNumerator = metricObservations.control?.numerator || 0;
  const controlDenominator = metricObservations.control?.denominator || 0;
  const treatmentNumerator = metricObservations.treatment?.numerator || 0;
  const treatmentDenominator = metricObservations.treatment?.denominator || 0;

  const calculatePercentage = useCallback(
    (numerator: number, denominator: number): number => {
      if (denominator === 0) return 0;
      return (numerator / denominator) * 100;
    },
    [],
  );

  const formatPercentage = useCallback(
    (value: number): string => `${value.toFixed(2)}%`,
    [],
  );

  const calculatePPChange = useCallback(
    (
      controlValue: number,
      treatmentValue: number,
      controlDenominator: number,
      treatmentDenominator: number,
    ): number => {
      const controlRate = controlValue / controlDenominator;
      const treatmentRate = treatmentValue / treatmentDenominator;
      return (treatmentRate - controlRate) * 100;
    },
    [],
  );

  const formatPPChange = useCallback(
    (value: number): string => `${value > 0 ? "+" : ""}${value.toFixed(2)}pp`,
    [],
  );

  const formattedPValue = useCallback((pValue: number) => {
    if (pValue < 0.0001) return "p < 0.0001";
    return `p = ${pValue.toFixed(4)}`;
  }, []);

  const formattedCI = useCallback((lower: number, upper: number) => {
    return `[${(lower * 100).toFixed(2)}%, ${(upper * 100).toFixed(2)}%]`;
  }, []);

  const formattedMDE = useCallback((mde?: number) => {
    if (mde === undefined) return "—";
    return `+${(mde * 100).toFixed(2)}pp`;
  }, []);

  const formattedPoweredEffect = useCallback((poweredEffect: number) => {
    if (poweredEffect === 0) return "0.00%";
    return `${poweredEffect > 0 ? "+" : ""}${(poweredEffect * 100).toFixed(
      2,
    )}%`;
  }, []);

  const getMdeColor = useCallback(
    (mde?: number) => {
      if (mde === undefined) return "text-gray-500";

      const ppChange = calculatePPChange(
        controlNumerator,
        treatmentNumerator,
        controlDenominator,
        treatmentDenominator,
      );

      return Math.abs(ppChange) >= mde ? "text-green-500" : "text-gray-500";
    },
    [
      calculatePPChange,
      controlNumerator,
      treatmentNumerator,
      controlDenominator,
      treatmentDenominator,
    ],
  );

  const getPoweredEffectColor = useCallback(
    (poweredEffect: number, mde?: number) => {
      if (mde === undefined) return "text-gray-500";
      return poweredEffect >= mde ? "text-red-500" : "text-green-500";
    },
    [],
  );

  const getIncreaseOrDecreaseText = useCallback(() => {
    const ppChange = calculatePPChange(
      controlNumerator,
      treatmentNumerator,
      controlDenominator,
      treatmentDenominator,
    );
    if (ppChange > 0) return "Increase";
    if (ppChange < 0) return "Decrease";
    return "No change";
  }, [
    controlNumerator,
    treatmentNumerator,
    controlDenominator,
    treatmentDenominator,
    calculatePPChange,
  ]);

  return (
    <Card className="!gap-0 pb-0">
      <div className="flex justify-between pb-4">
        <div className="flex flex-col gap-1">
          <h2 className="font-semibold">{metricDetails.metric_id.name}</h2>
          <p className="py-0.5 text-xs text-gray-400">
            {metricDetails.metric_id.description}
          </p>
        </div>
        <div className="flex flex-col items-end justify-end gap-2">
          <MetricIndicator
            isSignificant={statisticallySignificant}
            isSignificantText="Statistically significant"
            isSignificantIcon={
              <CheckCircleIcon className="mr-1 inline h-4 w-4" />
            }
            notSignificantText="Not statistically significant"
            notSignificantIcon={<XCircleIcon className="mr-1 inline h-4 w-4" />}
          />
          <MetricIndicator
            isSignificant={practicallySignificant}
            isSignificantText="Practically significant"
            isSignificantIcon={
              <CheckCircleIcon className="mr-1 inline h-4 w-4" />
            }
            notSignificantText="Not practically significant"
            notSignificantIcon={<XCircleIcon className="mr-1 inline h-4 w-4" />}
          />
        </div>
      </div>
      <div className="-mx-4 grid grid-cols-3 divide-x divide-gray-300 border-y border-gray-300">
        <div className="flex flex-col gap-0.5 px-4 py-4">
          <p className="pb-0.5 font-mono text-xs text-gray-500">
            Control ({variants?.find((v) => v.type === "control")?.key || "—"})
          </p>
          <p className="text-xl font-semibold">
            {formatPercentage(
              calculatePercentage(controlNumerator, controlDenominator),
            )}
          </p>
          <p className="font-mono text-xs text-gray-500">
            {controlNumerator} / {controlDenominator}
          </p>
        </div>
        <div className="flex flex-col gap-0.5 px-4 py-4">
          <p className="pb-0.5 font-mono text-xs text-gray-500">
            Treatment (
            {variants?.find((v) => v.type === "treatment")?.key || "—"})
          </p>
          <p className="text-xl font-semibold">
            {formatPercentage(
              calculatePercentage(treatmentNumerator, treatmentDenominator),
            )}
          </p>
          <p className="font-mono text-xs text-gray-500">
            {treatmentNumerator} / {treatmentDenominator}
          </p>
        </div>
        <div className="flex flex-col gap-0.5 px-4 py-4">
          <p className="pb-0.5 font-mono text-xs text-gray-500">Difference</p>
          <p className="text-xl font-semibold text-blue-600">
            {formatPPChange(
              calculatePPChange(
                controlNumerator,
                treatmentNumerator,
                controlDenominator,
                treatmentDenominator,
              ),
            )}
          </p>
          <p className="font-mono text-xs text-gray-500">
            {getIncreaseOrDecreaseText()}
          </p>
        </div>
      </div>
      <div className="-mx-4 flex flex-col divide-y divide-gray-300 text-xs">
        <div className="flex flex-row items-center justify-between px-4 py-3">
          <p className="text-gray-500">P Value</p>
          <p className="text-gray-500">
            {formattedPValue(zTestResult.p_value)}
          </p>
        </div>
        <div className="flex flex-row items-center justify-between px-4 py-3">
          <p className="text-gray-500">95% Confidence Interval</p>
          <p className="text-gray-500">
            {formattedCI(zTestResult.ci_lower, zTestResult.ci_upper)}
          </p>
        </div>
        <div className="flex flex-row items-center justify-between px-4 py-3">
          <p className="text-gray-500">Required MDE</p>
          <p className={`${getMdeColor(metricDetails.mde)}`}>
            {formattedMDE(metricDetails.mde)}
          </p>
        </div>
        <div className="flex flex-row items-center justify-between px-4 py-3">
          <p className="text-gray-500">Powered Effect</p>
          <p
            className={`${getPoweredEffectColor(zTestResult.powered_effect, metricDetails.mde)}`}
          >
            {formattedPoweredEffect(zTestResult.powered_effect)}
          </p>
        </div>
      </div>
    </Card>
  );
};

export default MetricResultCard;
