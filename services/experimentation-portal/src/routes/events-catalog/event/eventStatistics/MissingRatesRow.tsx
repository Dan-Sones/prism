import FieldKeyDataTypePill from "../../../../components/fieldKey/FieldKeyDataTypePill";
import type { MissingTableRateRow } from "./MissingRatesTable";

const missingRateColor = (rate: number) => {
  if (rate > 10) return "text-red-500";
  if (rate > 1) return "text-yellow-500";
  return "text-green-500";
};

const MissingRatesRow = (props: MissingTableRateRow) => {
  const { fieldKey, missingRate, fieldType } = props;
  return (
    <tr
      key={fieldKey}
      className="border-b border-gray-100 text-xs last:border-0"
    >
      <td className="flex items-center gap-2 py-2">
        <p className="font-mono">{fieldKey}</p>
        <FieldKeyDataTypePill dataType={fieldType} />
      </td>
      <td className={`py-2 ${missingRateColor(missingRate)}`}>
        {missingRate}%
      </td>
    </tr>
  );
};

export default MissingRatesRow;
