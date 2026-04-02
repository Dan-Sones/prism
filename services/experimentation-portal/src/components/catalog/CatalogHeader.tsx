import PageTitle from "../title/PageTitle";
import PrimaryButton from "../button/PrimaryButton";
import PlusCircleIcon from "../icons/PlusCircleIcon";

interface CatalogHeaderProps {
  title: string;
  createButtonText: string;
  onCreate: () => void;
}

const CatalogHeader = (props: CatalogHeaderProps) => {
  const { title, createButtonText, onCreate } = props;

  return (
    <div className="mb-6 min-w-full">
      <div className="mb-3 flex flex-row items-center justify-between">
        <PageTitle>{title}</PageTitle>
        <PrimaryButton onClick={onCreate} className="text-sm">
          <span className="flex flex-row items-center gap-1.5">
            <PlusCircleIcon className="size-5" />
            {createButtonText}
          </span>
        </PrimaryButton>
      </div>
    </div>
  );
};

export default CatalogHeader;
