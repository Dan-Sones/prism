import SearchIcon from "../../components/icons/SearchIcon";

interface EventsCatalogSearchProps {
  onSearch?: (query: string) => void;
}

const EventsCatalogSearch = ({ onSearch }: EventsCatalogSearchProps) => {
  return (
    <div className="flex w-full flex-row items-center gap-2 rounded-md border border-gray-400 bg-white px-4 py-2 shadow-sm">
      <SearchIcon className="size-5 shrink-0 text-gray-400" />
      <input
        type="text"
        placeholder="Search"
        className="w-full outline-none"
        onChange={(e) => onSearch?.(e.target.value)}
      />
    </div>
  );
};

export default EventsCatalogSearch;
