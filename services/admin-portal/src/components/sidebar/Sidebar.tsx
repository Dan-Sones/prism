import PrismLogo from "../../assets/prism-logo-minimal-dark.svg";
import EventCatalogIcon from "./icons/EventCatalogIcon";
import HomeIcon from "./icons/HomeIcon";
import SliderIcon from "./icons/SliderIcon";
import NavItem, { type NavItemProps } from "./NavItem";

interface SidebarProps {
  isOpen: boolean;
}

const Sidebar = (props: SidebarProps) => {
  const { isOpen } = props;

  const navItems: Array<NavItemProps> = [
    {
      name: "Home",
      href: "/home",
      icon: HomeIcon,
    },
    {
      name: "Settings",
      href: "/test",
      icon: SliderIcon,
    },
    {
      name: "Event Catalog",
      href: "/events-catalog",
      icon: EventCatalogIcon,
    },
  ];

  return (
    <aside
      id="sidebar"
      className={`fixed z-30 h-full w-80 bg-white transition-transform duration-300 ease-in-out lg:relative lg:translate-x-0 ${isOpen ? "translate-x-0" : "-translate-x-full"} `}
    >
      <div className="flex flex-col gap-3 pt-1 pl-1">
        <a
          href="/home"
          className="font-brand flex cursor-pointer flex-row items-center gap-2 p-4 pb-1.5 text-2xl tracking-wide text-slate-900"
        >
          <span>
            <img src={PrismLogo} alt="Prism Logo" className="size-9" />
          </span>
          <span className="mt-1">Prism</span>
        </a>
        <nav>
          {navItems.map((item) => (
            <NavItem key={item.name} {...item} />
          ))}
        </nav>
      </div>
    </aside>
  );
};

export default Sidebar;
