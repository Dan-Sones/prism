import Card from "../../../../../components/card/Card";

interface SampleSizeRequiredProps {
  requiredSampleSize: number | undefined;
}

const SampleSizeRequired = ({
  requiredSampleSize,
}: SampleSizeRequiredProps) => {
  const perVariant =
    requiredSampleSize != null
      ? Math.floor(requiredSampleSize / 2).toLocaleString()
      : "—";

  return (
    <Card>
      <h2 className="text-lg font-semibold">Required Sample Size</h2>
      <p className="text-3xl font-bold">{perVariant}</p>
      <p className="text-sm">
        <span className="font-semibold">Per Variant</span> — based on observed
        metric variance and experiment configuration
      </p>
    </Card>
  );
};

export default SampleSizeRequired;
