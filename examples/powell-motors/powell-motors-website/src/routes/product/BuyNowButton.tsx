import { useStringFlagValue } from "@openfeature/react-sdk";
import PrimaryButton from "../../components/buttons/PrimaryButton";

type ButtonVariant = "button-black" | "button-red";

const CONTROL_VARIANT: ButtonVariant = "button-black";
const TREATMENT_VARIANT: ButtonVariant = "button-red";

const variantClassName: Record<ButtonVariant, string> = {
  "button-black": "bg-black-700!",
  "button-red": "bg-red-700!",
};

interface BuyNowButtonProps {
  handlePurchase: () => void;
}

const BuyNowButton = (props: BuyNowButtonProps) => {
  const { handlePurchase } = props;
  const evaluatedVariant = useStringFlagValue(
    "the-homer-product-page-button",
    CONTROL_VARIANT,
  );

  const buttonVariant: ButtonVariant =
    evaluatedVariant === TREATMENT_VARIANT
      ? TREATMENT_VARIANT
      : CONTROL_VARIANT;

  return (
    <PrimaryButton
      className={variantClassName[buttonVariant]}
      onClick={handlePurchase}
    >
      Buy Now
    </PrimaryButton>
  );
};

export default BuyNowButton;
