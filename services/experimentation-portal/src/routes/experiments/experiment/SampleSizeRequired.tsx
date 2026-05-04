import Card from "../../../components/card/Card";

interface SampleSizeRequiredProps {
  requiredSampleSize: number | undefined;
}

const SampleSizeRequired = (props: SampleSizeRequiredProps) => {
  return (
    <Card>
      <h2 className="text-lg font-semibold">Sample Size Required</h2>
      <p className="text-3xl font-bold">
        {props.requiredSampleSize?.toLocaleString()}
      </p>
      <p className="text-sm">
        <span className="font-semibold">Per Variant</span> - Based on observed
        metricvariance and experiment configuration parameters
      </p>
    </Card>
  );
};

export default SampleSizeRequired;
