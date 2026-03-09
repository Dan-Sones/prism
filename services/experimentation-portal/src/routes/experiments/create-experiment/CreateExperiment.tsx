import Card from "../../../components/card/Card";
import type { JourneyBarItemT } from "../../../components/journeyBar/JourneyBar";
import PageTitle from "../../../components/title/PageTitle";
import JourneyBar from "../../../components/journeyBar/JourneyBar";

const CreateExperiment = () => {
  const journeyBarItems: Array<JourneyBarItemT> = [
    {
      label: "Experiment details",
      onClick: () => {},
      complete: true,
    },
    {
      label: "Define variations",
      onClick: () => {},
      complete: false,
    },
    {
      label: "Targeting",
      onClick: () => {},
      complete: false,
    },
  ];

  return (
    <>
      <PageTitle>Create Experiment</PageTitle>
      <JourneyBar items={journeyBarItems} activeItemIndex={0} />
      <Card>yo yo yo</Card>
    </>
  );
};

export default CreateExperiment;
