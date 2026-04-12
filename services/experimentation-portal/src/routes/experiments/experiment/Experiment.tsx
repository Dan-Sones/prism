import { useParams } from "react-router";

const Experiment = () => {
  const params = useParams();

  return <div>Experiment details page: {params.id}</div>;
};

export default Experiment;
