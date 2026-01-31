import ProductAbout from "./ProductAbout";
import ProductDetails from "./ProductDetails";

const ProductPage = () => {
  return (
    <main className="mx-auto my-12 flex max-w-7xl flex-col px-4 sm:px-6 lg:px-8">
      <div className="flex w-full flex-col gap-8 lg:flex-row">
        <div className="flex-1">
          <img
            src="../../../public/the-homer.jpg"
            alt="Car"
            className="h-auto w-full rounded-lg shadow-md"
          />
        </div>
        <ProductDetails
          title="The Homer"
          price="£10,000"
          description="Built To Order - 6-8 Week Delivery Time - Extremely Large Beverage Holder"
        />
      </div>
      <ProductAbout />
    </main>
  );
};

export default ProductPage;
