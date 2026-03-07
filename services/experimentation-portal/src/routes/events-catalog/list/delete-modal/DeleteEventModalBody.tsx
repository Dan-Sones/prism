import BaseModal from "../../../../components/modal/BaseModal";

interface DeleteEventModalProps {
  isOpen: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

const DeleteEventModal = (props: DeleteEventModalProps) => {
  const { isOpen, onConfirm, onCancel } = props;
  return (
    <BaseModal
      isOpen={isOpen}
      onRequestClose={onCancel}
      overlayClassName="fixed inset-0 z-[9999] flex items-center justify-center bg-black/40"
      className="mx-4 w-full max-w-sm rounded-lg bg-white shadow-lg outline-none"
    >
      <div className="flex flex-col gap-6 p-6">
        <div className="flex flex-col gap-1">
          <h2 className="text-lg font-semibold text-gray-900">
            Delete event type
          </h2>
          <p className="text-sm text-gray-500">
            Are you sure you want to delete this event type? This action cannot
            be undone.
          </p>
        </div>
        <div className="flex justify-end gap-3">
          <button
            type="button"
            onClick={onCancel}
            className="cursor-pointer rounded-lg px-3 py-2.5 text-sm text-gray-500 hover:text-gray-700"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="cursor-pointer rounded-lg bg-red-500 px-3 py-2.5 text-sm text-white transition-colors duration-200 hover:bg-red-600"
          >
            Delete
          </button>
        </div>
      </div>
    </BaseModal>
  );
};

export default DeleteEventModal;
