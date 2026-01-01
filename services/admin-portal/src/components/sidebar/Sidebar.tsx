const Sidebar = () => {
  const navItems = [{ name: "Home", href: "/home" }];

  return (
    <section
      id="sidebar"
      className="w-64 bg-gray-800 text-white rounded-4xl m-4 p-4"
    >
      <div className="flex flex-col pt-4 gap-3">
        <nav>
          {navItems.map((item) => (
            <a
              key={item.name}
              href={item.href}
              className="block py-2 px-4 rounded hover:bg-gray-700"
            >
              {item.name}
            </a>
          ))}
        </nav>
      </div>
    </section>
  );
};

export default Sidebar;
