const Spinner = () => {
  return (
    <div role="status">
      <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-purple-600" />
      <span className="sr-only">Loading...</span>
    </div>
  );
};

export default Spinner;
