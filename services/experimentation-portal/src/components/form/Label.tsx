interface LabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
  required?: boolean;
}

const Label = ({ required, children, ...rest }: LabelProps) => {
  return (
    <label className="text-sm text-gray-600" {...rest}>
      {children}
      {required && <span className="pl-1 text-red-400">*</span>}
    </label>
  );
};

export default Label;
