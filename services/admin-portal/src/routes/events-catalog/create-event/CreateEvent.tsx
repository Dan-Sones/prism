import { useNavigate } from "react-router";
import { useForm, FormProvider } from "react-hook-form";
import PrimaryButton from "../../../components/button/PrimaryButton";
import CreateEventDetails from "./CreateEventDetails";
import CreateEventFields from "./CreateEventFields";
import {
  createEventType,
  type CreateEventTypeRequest,
} from "../../../api/eventsCatalog";
import { useMutation } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import type { ProblemDetail } from "../../../api/base/problem";
import { useErrorBanner } from "../../../context/ErrorBannerContext";

const CreateEvent = () => {
  const navigate = useNavigate();

  const { setErrorMessage } = useErrorBanner();

  const form = useForm<CreateEventTypeRequest>({
    mode: "onChange",
    defaultValues: {
      fields: [{ name: "", fieldKey: "", dataType: "string" }],
    },
  });

  const onCancel = () => {
    navigate("/events-catalog");
  };

  const onSubmit = (data: CreateEventTypeRequest) => {
    mutation.mutate(data);
  };

  const mutation = useMutation<
    void,
    AxiosError<ProblemDetail>,
    CreateEventTypeRequest
  >({
    mutationFn: createEventType,
    onSuccess: () => {
      // TODO: redirect to event type details page after creation instead of just going back to the list
      navigate("/events-catalog");
    },
    onError: (error) => {
      const baseErrorMessage = "Failed to create event type:";

      if (error.response?.data.detail) {
        setErrorMessage(baseErrorMessage + " " + error.response.data.detail);
        return;
      }

      setErrorMessage(baseErrorMessage);
    },
  });

  return (
    <div className="flex h-full grow flex-col items-center justify-center">
      <div className="grow overflow-y-auto">
        <FormProvider {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="flex max-w-3xl flex-col gap-6 px-4 py-6 md:px-10 md:py-8 lg:px-20 lg:py-10">
              <section className="rounded-md bg-white p-6 shadow-xs">
                <h1 className="text-2xl font-semibold lg:text-3xl">
                  Create Event Type
                </h1>
                <p className="mt-1 text-sm text-gray-400">
                  Define a new event type and its fields for use in metrics.
                </p>
              </section>

              <CreateEventDetails />
              <CreateEventFields />
              <div className="flex max-w-3xl justify-end gap-3">
                <button
                  type="button"
                  onClick={onCancel}
                  className="cursor-pointer rounded-lg px-3 py-2.5 text-sm text-gray-500 hover:text-gray-700"
                >
                  Cancel
                </button>
                <PrimaryButton type="submit" className="text-sm">
                  Create Event Type
                </PrimaryButton>
              </div>
            </div>
          </form>
        </FormProvider>
      </div>
    </div>
  );
};

export default CreateEvent;
