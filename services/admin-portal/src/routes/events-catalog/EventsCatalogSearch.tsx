import SearchIcon from "../../components/icons/SearchIcon";

interface EventsCatalogSearchProps {
  onSearch?: (query: string) => void;
}

const EventsCatalogSearch = ({ onSearch }: EventsCatalogSearchProps) => {
  return (
    <div className="flex h-12 w-full flex-row items-center gap-2 rounded-md border border-gray-300 bg-white px-4 py-2 shadow-xs focus-within:border-gray-400 lg:max-w-96">
      <SearchIcon className="size-5 shrink-0 text-gray-400" />
      <input
        type="text"
        placeholder="Search"
        className="w-full text-sm outline-none"
        onChange={(e) => onSearch?.(e.target.value)}
      />
    </div>
  );
};

export default EventsCatalogSearch;
