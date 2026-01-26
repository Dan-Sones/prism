import PrimaryButton from "../../components/buttons/PrimaryButton";
import SecondaryButton from "../../components/buttons/SecondaryButton";

interface ProductDetailsProps {
  title: string;
  price: string;
  description: string;
}

const ProductDetails = ({ title, price, description }: ProductDetailsProps) => {
  return (
    <div className="flex-1 gap-3 flex flex-col justify-center">
      <h2 className="text-4xl font-semibold">{title}</h2>
      <p className="text-xl font-bold">{price}</p>
      <p className="text font-light italic">{description}</p>
      <div className="flex flex-col gap-2 w-full">
        <PrimaryButton>Buy now</PrimaryButton>
        <SecondaryButton>Add to bag</SecondaryButton>
      </div>
      <div className="mt-4 text-gray-800 flex flex-row gap-8 justify-center">
        <div className="flex flex-row gap-2 items-center">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-10"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M8.25 18.75a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 0 1-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 0 0-3.213-9.193 2.056 2.056 0 0 0-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 0 0-10.026 0 1.106 1.106 0 0 0-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12"
            />
          </svg>
          <p className="text-sm">Delivered to your door</p>
        </div>
        <div className="flex flex-row gap-2 items-center">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-10"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z"
            />
          </svg>

          <p className="text-sm">Quality Assured</p>
        </div>
      </div>
    </div>
  );
};

export default ProductDetails;
