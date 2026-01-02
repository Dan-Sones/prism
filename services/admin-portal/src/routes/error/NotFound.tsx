const NotFound = () => {
  return (
    <div className="flex h-full grow flex-col items-center justify-center">
      <div className="text-center">
        <h1 className="mb-4 text-5xl font-bold text-slate-50">404</h1>
        <h2 className="mb-2 text-2xl text-slate-50">Page Not Found</h2>
        <p className="font-secondary mb-4 text-lg text-slate-50">
          Sorry, the page you are looking for does not exist.
        </p>
      </div>
    </div>
  );
};

export default NotFound;
