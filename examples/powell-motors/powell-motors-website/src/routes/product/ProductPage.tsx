import ProductAbout from "./ProductAbout";
import ProductDetails from "./ProductDetails";

const ProductPage = () => {
  return (
    <main className="mx-auto my-12 max-w-7xl px-4 sm:px-6 lg:px-8 flex flex-col ">
      <div className="flex flex-col lg:flex-row w-full gap-8">
        <div className="flex-1">
          <img
            src="../../../public/the-homer.jpg"
            alt="Car"
            className="h-auto rounded-lg shadow-md w-full"
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
