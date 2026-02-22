import { useFormContext } from "react-hook-form";
import TextInput from "../../../components/form/TextInput";
import type { CreateEventTypeRequest } from "../../../api/eventsCatalog";

const CreateEventDetails = () => {
  const { register, formState } = useFormContext<CreateEventTypeRequest>();
  const { errors } = formState;

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
            placeholder="e.g. Purchase Completed"
            {...register("name", {
              required: "Name is required",
              maxLength: {
                value: 100,
                message: "Name must be less than 100 characters",
              },
            })}
          />
          {errors.name && (
            <p className="text-xs text-red-500">{errors.name.message}</p>
          )}
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-sm text-gray-600" htmlFor="eventKey">
            Event Key <span className="text-red-400">*</span>
          </label>
          <p className="text-xs text-gray-400">
            The event key must match the key found in the event payload EXACTLY.
          </p>
          <TextInput
            id="eventKey"
            placeholder="e.g. purchase_completed"
            {...register("eventKey", {
              required: "Event key is required",
              maxLength: {
                value: 50,
                message: "Event key must be less than 50 characters",
              },
              pattern: {
                value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
                message:
                  "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
              },
            })}
          />
          {errors.eventKey && (
            <p className="text-xs text-red-500">{errors.eventKey.message}</p>
          )}
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-sm text-gray-600" htmlFor="description">
            Description
          </label>
          <textarea
            {...register("description")}
            id="description"
            name="description"
            rows={3}
            placeholder="Described the context in which this event is triggered, what it represents, and any other relevant details."
            className="w-full rounded-md border border-slate-200 bg-gray-50 px-3 py-2 text-sm text-slate-800 transition duration-300 placeholder:text-slate-400 hover:border-slate-300 focus:border-slate-400 focus:outline-none"
          />
        </div>
      </div>
    </section>
  );
};

export default CreateEventDetails;
