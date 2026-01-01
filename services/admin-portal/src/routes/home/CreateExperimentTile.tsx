import { useNavigate } from "react-router";
import PrimaryButton from "../../components/button/PrimaryButton";
import Tile from "../../components/tiles/Tile";

const CreateExperimentTile = () => {
  const navigate = useNavigate();

  const onClick = () => {
    navigate("/create-experiment");
  };

  return (
    <Tile>
      <div className="flex flex-col gap-4 justify-center items-center p-8">
        <h1 className="text-5xl text-slate-100">Welcome to Prism!</h1>
        <PrimaryButton className="w-fit" onClick={onClick}>
          Create Experiment
        </PrimaryButton>
      </div>
    </Tile>
  );
};

export default CreateExperimentTile;
