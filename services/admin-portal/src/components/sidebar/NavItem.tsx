import { useLocation } from "react-router";

export interface NavItemProps {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
}

const NavItem = (item: NavItemProps) => {
  const { name, href, icon: Icon } = item;

  const location = useLocation();

  const isActive = location.pathname === href;

  return (
    <a
      key={name}
      href={href}
      className="flex flex-row items-center gap-2 rounded px-4 py-2 transition-all duration-200 hover:bg-gray-200"
    >
      <span className="flex items-center justify-center">
        <Icon
          className={`size-6 ${isActive ? "text-green-600" : "text-slate-600"}`}
        />
      </span>
      <p className={`${isActive ? "font-semibold" : "text-slate-600"}`}>
        {name}
      </p>
    </a>
  );
};

export default NavItem;
