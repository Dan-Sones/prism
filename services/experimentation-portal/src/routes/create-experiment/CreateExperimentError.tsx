interface CreateExperimentErrorProps {
  message: string | null;
}

const CreateExperimentError = (props: CreateExperimentErrorProps) => {
  const { message } = props;

  if (!message) {
    return null;
  }

  return (
    <span className="flex w-full flex-row gap-4 rounded-md bg-red-100 p-2 text-center text-red-500">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth={1.5}
        stroke="currentColor"
        className="size-14"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
        />
      </svg>
      <span className="flex flex-col items-start justify-start">
        <p className="font-semibold">Something Went Wrong: </p>
        <p className="text-sm font-light">{message}</p>
      </span>
    </span>
  );
};

export default CreateExperimentError;
