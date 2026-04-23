import Card from "../../../components/card/Card";

interface PostAAConfigProps {
  experimentId: string;
}

const PostAAConfig = ({ experimentId }: PostAAConfigProps) => {
  return (
    <Card>
      <h2 className="font-mono">A/B Test Configuration</h2>

      <div className="flex flex-col flex-wrap gap-2"></div>
    </Card>
  );
};

export default PostAAConfig;
