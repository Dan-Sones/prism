import { motion, AnimatePresence } from "framer-motion";
import { useErrorBanner } from "../../context/ErrorBannerContext";
import ExclamationCircleIcon from "../icons/ExclamationCircleIcon";
import XMarkIcon from "../icons/XMarkIcon";

const ErrorBanner = () => {
  const { setErrorMessage, errorMessage } = useErrorBanner();

  const onClose = () => {
    setErrorMessage(null);
  };

  return (
    <AnimatePresence>
      {errorMessage && (
        <motion.section
          initial={{ y: -40, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          exit={{ y: -40, opacity: 0 }}
          transition={{ type: "spring", stiffness: 400, damping: 30 }}
          className="absolute inset-0 top-1 mx-5 mt-2 flex max-h-12 rounded-xl bg-red-500 sm:max-h-10"
          style={{ zIndex: "10000" }}
        >
          <div className="flex grow items-center justify-between">
            <div className="flex grow items-center justify-between px-4">
              <div className="flex items-center justify-between gap-2">
                <ExclamationCircleIcon className="size-6 text-white" />

                <p className="text-sm text-white">{errorMessage}</p>
              </div>
              <button type="button" onClick={onClose}>
                <XMarkIcon className="size-6 cursor-pointer text-white" />
              </button>
            </div>
          </div>
        </motion.section>
      )}
    </AnimatePresence>
  );
};

export default ErrorBanner;
