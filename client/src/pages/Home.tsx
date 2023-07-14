import { useGetFiles } from "@/services/file.service";

const Home = () => {
  const { data, isLoading } = useGetFiles();
  console.log(data);
  if (isLoading) {
    return <h1>Loading</h1>;
  }
  return (
    <div className="flex flex-wrap gap-5 p-10">
      {data?.files?.map((file, index) => (
        <div key={file.name + index} className="w-72 h-72 rounded-md border-2">
          {file.name}
        </div>
      ))}
    </div>
  );
};

export default Home;
