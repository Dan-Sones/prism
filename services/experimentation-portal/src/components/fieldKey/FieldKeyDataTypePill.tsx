const FieldKeyDataTypePill = (props: { dataType: string }) => {
  const dataTypeBadgeColor: Record<string, string> = {
    string: "bg-blue-100 text-blue-700",
    int: "bg-green-100 text-green-700",
    float: "bg-yellow-100 text-yellow-700",
    boolean: "bg-purple-100 text-purple-700",
    timestamp: "bg-orange-100 text-orange-700",
  };

  return (
    <span
      className={`rounded-full px-2 py-0.5 text-xs ${dataTypeBadgeColor[props.dataType] ?? "bg-gray-100 text-gray-600"}`}
    >
      {props.dataType}
    </span>
  );
};

export default FieldKeyDataTypePill;
