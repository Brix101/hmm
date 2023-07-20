import FileCard from "@/components/FileCard";
import { cn } from "@/lib/utils";
import { useQueryFile } from "@/services/file.service";
import { useBoundStore } from "@/store";

const Home = () => {
  const {
    file: { pathHistory },
    setPathHistory,
  } = useBoundStore();

  const { data, isLoading, error } = useQueryFile(pathHistory);

  const breadCrumbs = pathHistory?.split("/");

  if (isLoading) {
    return <h1>Loading</h1>;
  }

  if (error) {
    return (
      <div
        className="relative py-3 px-4 text-red-700 bg-red-100 rounded border border-red-400"
        role="alert"
      >
        <strong className="font-bold">Message: </strong>
        <span className="block sm:inline">{error.message}</span>
      </div>
    );
  }

  function handleBreadCrumbsClick({ index }: { index: number }) {
    setPathHistory(`${breadCrumbs?.slice(0, index + 1).join("/")}`);
  }

  return (
    <>
      <div>
        <div className="flex items-center space-x-1 text-sm capitalize text-muted-foreground">
          {breadCrumbs?.map((url, index) => {
            const isActive = breadCrumbs?.length - 1 === index;

            const urlName = url.length > 1 ? url : "Files";
            return (
              <div className="flex items-center" key={index}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className={cn("w-4 h-4", index === 0 ? "hidden" : "")}
                  aria-hidden="true"
                >
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
                <button
                  className={`${isActive
                      ? "text-foreground"
                      : "cursor-pointer hover:underline "
                    }`}
                  onClick={() => handleBreadCrumbsClick({ index })}
                  disabled={isActive}
                >
                  {urlName}
                </button>
              </div>
            );
          })}
        </div>
      </div>

      <div className="grid gap-2 grid-cols-file">
        {data?.files?.map((file, index) => (
          <FileCard key={file.name + index} file={file} />
        ))}
      </div>
    </>
  );
};

export default Home;
