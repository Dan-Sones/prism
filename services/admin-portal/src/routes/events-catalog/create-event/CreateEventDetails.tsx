import TextInput from "../../../components/form/TextInput";

const CreateEventDetails = () => {
  return (
    <section className="rounded-md bg-white p-6 shadow-xs">
      <h2 className="mb-4 text-sm font-semibold text-gray-700">
        Event Details
      </h2>
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-1">
          <label className="text-sm text-gray-600" htmlFor="name">
            Name <span className="text-red-400">*</span>
          </label>
          <TextInput
            id="name"
            name="name"
            placeholder="e.g. purchase_completed"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-sm text-gray-600" htmlFor="description">
            Description
          </label>
          <textarea
            id="description"
            name="description"
            rows={3}
            placeholder="Describe when this event is fired..."
            className="w-full rounded-md border border-slate-200 bg-gray-50 px-3 py-2 text-sm text-slate-800 transition duration-300 placeholder:text-slate-400 hover:border-slate-300 focus:border-slate-400 focus:outline-none"
          />
        </div>
      </div>
    </section>
  );
};

export default CreateEventDetails;
