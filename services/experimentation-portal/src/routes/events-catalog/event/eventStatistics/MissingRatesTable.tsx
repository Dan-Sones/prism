import type { DataType } from "../../../../api/eventsCatalog";
import MissingRatesRow from "./MissingRatesRow";

export interface MissingTableRateRow {
  fieldKey: string;
  missingRate: number;
  fieldType: DataType;
}

interface MissingRatesTableProps {
  missingRates?: Array<MissingTableRateRow>;
}

const MissingRatesTable = (props: MissingRatesTableProps) => {
  const { missingRates } = props;
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-gray-200 text-left text-xs text-gray-400">
          <th className="pb-2 font-normal">Field Key</th>
          <th className="pb-2 font-normal">Missing Rate</th>
        </tr>
      </thead>
      <tbody>
        {missingRates?.map(({ fieldKey, missingRate, fieldType }) => (
          <MissingRatesRow
            key={fieldKey}
            fieldKey={fieldKey}
            missingRate={missingRate}
            fieldType={fieldType}
          />
        ))}
      </tbody>
    </table>
  );
};

export default MissingRatesTable;
