const NavBar = () => {
  return (
    <nav className="h-16 bg-gray-800 text-white">
      <div className="flex h-full w-full justify-between p-4">
        <h1 className="text-2xl">Powell Motors</h1>
        <button className="cursor-pointer rounded-xl border-2 border-white px-2">
          Login
        </button>
      </div>
    </nav>
  );
};

export default NavBar;
