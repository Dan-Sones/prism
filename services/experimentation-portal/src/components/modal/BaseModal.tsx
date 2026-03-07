import Modal from "react-modal";
interface BaseModalProps extends Modal.Props {
  children?: React.ReactNode;
}
const BaseModal = (props: BaseModalProps) => {
  return <Modal {...props}>{props.children}</Modal>;
};

export default BaseModal;
