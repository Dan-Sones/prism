import prismLogo from "../../assets/prism-logo.svg";
import HomeIcon from "./icons/HomeIcon";
import SliderIcon from "./icons/SliderIcon";
import NavItem, { type NavItemProps } from "./NavItem";

const Sidebar = () => {
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
  ];

  return (
    <section id="sidebar" className="w-72 border-r border-gray-200 bg-gray-50">
      <div className="flex flex-col gap-3 pt-3">
        <a
          href="/home"
          className="font-brand flex cursor-pointer flex-row items-center gap-2 p-4 text-2xl tracking-wide text-slate-900"
        >
          <span>
            <img src={prismLogo} alt="Prism Logo" className="size-9" />
          </span>
          <span className="mt-1">Prism</span>
        </a>
        <nav>
          {navItems.map((item) => (
            <NavItem key={item.name} {...item} />
          ))}
        </nav>
      </div>
    </section>
  );
};

export default Sidebar;
