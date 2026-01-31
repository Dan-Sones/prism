interface ProductItem {
  imageSrc: string;
  title: string;
  description: string;
}

const ProductAbout = () => {
  const items: ProductItem[] = [
    {
      imageSrc: "../../../public/homer-with-design.webp",
      title: "Design Language",
      description:
        "The Homer is a revolutionary vehicle designed with the modern driver in mind. Featuring an expansive beverage holder, it ensures that you stay refreshed on all your journeys. Its sleek design and advanced engineering make it not just a car, but a statement on wheels.",
    },
    {
      imageSrc: "../../../public/homer-in-factory.png",
      title: "For the Modern Man",
      description:
        "The Homer is a revolutionary vehicle designed with the modern driver in mind. Featuring an expansive beverage holder, it ensures that you stay refreshed on all your journeys. Its sleek design and advanced engineering make it not just a car, but a statement on wheels.",
    },
    {
      imageSrc: "../../../public/homer-with-engineers.png",
      title: "Engineered by Experts",
      description:
        "The Homer is a revolutionary vehicle designed with the modern driver in mind. Featuring an expansive beverage holder, it ensures that you stay refreshed on all your journeys. Its sleek design and advanced engineering make it not just a car, but a statement on wheels.",
    },
  ];

  return (
    <section>
      <hr className="my-8" />
      <div className="flex flex-col gap-12 md:flex-row">
        {items.map((item, index) => (
          <div key={index} className="flex-1">
            <img
              src={item.imageSrc}
              alt="The Homer Car"
              className="mb-6 h-auto w-full rounded-lg shadow-md"
            />
            <h2 className="mb-4 text-3xl font-semibold">{item.title}</h2>
            <p className="mb-4 text-lg">{item.description}</p>
          </div>
        ))}
      </div>
    </section>
  );
};

export default ProductAbout;
