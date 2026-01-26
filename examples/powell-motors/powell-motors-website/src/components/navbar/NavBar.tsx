const NavBar = () => {
  return (
    <nav className="bg-gray-800 h-16 text-white">
      <div className="flex justify-between h-full w-full p-4">
        <h1 className="text-2xl">Powell Motors</h1>
        <button className="border-white border-2 rounded-xl px-2 cursor-pointer">
          Login
        </button>
      </div>
    </nav>
  );
};

export default NavBar;
