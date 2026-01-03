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
      <div className="flex flex-col items-center justify-center gap-4 p-8">
        <h1 className="text-5xl">Welcome to Prism!</h1>
        <PrimaryButton className="w-fit" onClick={onClick} rounded>
          New Experiment
        </PrimaryButton>
      </div>
    </Tile>
  );
};

export default CreateExperimentTile;
